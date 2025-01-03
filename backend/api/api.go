package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	chaincodeID  = "basic"                                                    
	channelID    = "channel1"       
	port 		 = "localhost:7051"                                          
	aviationStackAPIURL = "https://api.aviationstack.com/v1/flights"
)
type Asset struct {
	ID          string `json:"id"`
	CompanyName string `json:"companyName"`
	AircraftID  string `json:"aircraftId"`
	Compliance  bool   `json:"compliance"`
	ReportDate  string `json:"reportDate"`
	Inspector   string `json:"inspector"`
	Description string `json:"description"`
}

type AssetHistory struct {
	Timestamp time.Time `json:"timestamp"`
	Asset     Asset     `json:"asset"`
}

type FlightData struct {
	FlightStatus   string `json:"flight_status"`
	Departure      string `json:"departure"`
	Arrival        string `json:"arrival"`
	FlightNumber   string `json:"flight_number"`
	AirlineName    string `json:"airline_name"`
	AircraftType   string `json:"aircraft_type"`
	DepartureTime  string `json:"departure_time"`
	ArrivalTime    string `json:"arrival_time"`
	DepartureCity  string `json:"departure_city"`
	ArrivalCity    string `json:"arrival_city"`
}

var gateway *client.Gateway
var clientConnection *grpc.ClientConn

func updateGrpcConnection(msp string) (*grpc.ClientConn, error) {
	tlsCertPath  := "../../fabric/test-network/organizations/peerOrganizations/org1.av.com/users/Admin@org1.av.com/msp/tlscacerts/tlsca.org1.av.com-cert.pem"
	port := "localhost:7051"
	if (msp == "Org2MSP") {
		tlsCertPath  = "../../fabric/test-network/organizations/peerOrganizations/org2.av.com/users/Admin@org2.av.com/msp/tlscacerts/tlsca.org2.av.com-cert.pem"
		port = "localhost:9051"
	}

	certificatePEM, err := os.ReadFile(tlsCertPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read TLS certificate: %w", err)
    }

    block, _ := pem.Decode(certificatePEM)
    if block == nil {
        return nil, fmt.Errorf("failed to decode PEM")
    }

    certPool := x509.NewCertPool()
    if !certPool.AppendCertsFromPEM(certificatePEM) {
        return nil, fmt.Errorf("failed to add certificate to pool")
    }

    transportCredentials := credentials.NewClientTLSFromCert(certPool, "")
    connection, err := grpc.Dial(port, grpc.WithTransportCredentials(transportCredentials))
    if err != nil {
        return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
    }

    return connection, nil
}

func FetchFlightData(flightID string) (*FlightData, error) {
    err := godotenv.Load("../../.env")
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    apiKey := os.Getenv("AVIATION_STACK_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("API key is missing")
    }

    url := fmt.Sprintf("%s?access_key=%s&flight_iata=%s", aviationStackAPIURL, apiKey, flightID)

    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch flight data: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    var flightDataResponse map[string]interface{}
    if err := json.Unmarshal(body, &flightDataResponse); err != nil {
        return nil, fmt.Errorf("failed to parse response JSON: %v", err)
    }

    flightDetails, ok := flightDataResponse["data"].([]interface{})
    if !ok || len(flightDetails) == 0 {
        return nil, fmt.Errorf("no flight data found for flight ID %s", flightID)
    }

    flight, ok := flightDetails[0].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid flight data format")
    }

    getString := func(data map[string]interface{}, key string) string {
        if value, ok := data[key].(string); ok {
            return value
        }
        return "Unknown"
    }

    getMap := func(data map[string]interface{}, key string) map[string]interface{} {
        if value, ok := data[key].(map[string]interface{}); ok {
            return value
        }
        return nil
    }

    departure := getMap(flight, "departure")
    arrival := getMap(flight, "arrival")
    airline := getMap(flight, "airline")
    flightInfo := getMap(flight, "flight")
    aircraft := getMap(flight, "aircraft")

    flightData := &FlightData{
        FlightStatus:   getString(flight, "flight_status"),
        Departure:      getString(departure, "estimated"),
        Arrival:        getString(arrival, "estimated"),
        FlightNumber:   getString(flightInfo, "iata"),
        AirlineName:    getString(airline, "name"),
        DepartureTime:  getString(departure, "estimated"),
        ArrivalTime:    getString(arrival, "estimated"),
        DepartureCity:  getString(departure, "airport"),
        ArrivalCity:    getString(arrival, "airport"),
        AircraftType:   getString(aircraft, "iata"),
    }

    return flightData, nil
}

func readAsset(c *gin.Context) {
	key := c.Param("key")

	contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

	response, err := contract.EvaluateTransaction("ReadAsset", key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query chaincode: %v", err)})
		return
	}

	var asset Asset

	err = json.Unmarshal(response, &asset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal response: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": asset})
}

func createAsset(c *gin.Context) {
	var request struct {
		ID          string `json:"id"`
		AircraftID  string `json:"aircraft_id"`
		ReportDate  string `json:"report_date"`
		Inspector   string `json:"inspector"`
		Description string `json:"description"`
		Compliance  bool   `json:"compliance"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	flightData, err := FetchFlightData(request.AircraftID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch flight data: %v", err)})
		return
	}

	companyName := flightData.AirlineName

	contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

	_, err = contract.SubmitTransaction(
		"CreateAsset", 
		request.ID, 
		companyName, 
		request.AircraftID, 
		request.ReportDate, 
		request.Inspector, 
		request.Description, 
		strconv.FormatBool(request.Compliance),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to invoke chaincode: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": request.ID})
}

func updateCompliance(c *gin.Context) {
    var request struct {
        ID        string `json:"id"`
        Compliance string `json:"compliance"` // "true" or "false"
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

    response, err := contract.SubmitTransaction("UpdateCompliance", request.ID, request.Compliance)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update compliance: %v", err)})
        return
    }

    if len(response) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": request.ID})
    } else {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update asset compliance"})
    }
}

func getAssetHistory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Asset ID is required"})
		return
	}

	contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

	result, err := contract.EvaluateTransaction("GetHistory", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to invoke chaincode: %v", err)})
		return
	}

	var history []AssetHistory
	err = json.Unmarshal(result, &history)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal history: %v", err)})
		return
	}

	c.JSON(http.StatusOK, history)
}

func assetExists(c *gin.Context) {
	id := c.Param("id")

	contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

	exists, err := contract.EvaluateTransaction("AssetExists", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query chaincode: %v", err)})
		return
	}

	var existsBool bool
	if err := json.Unmarshal(exists, &existsBool); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmarshal response: %v", err)})
		return
	}

	if existsBool {
		c.JSON(http.StatusOK, gin.H{"message": "Asset exists"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "Asset not found"})
	}
}

func populateLedger(c *gin.Context) {	// ONLY USE FOR TESTING PURPOSES
    contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

    _, err := contract.SubmitTransaction("CreateAsset", "asset123", "Company ABC", "A12345", "2024-12-30", "John Doe", "Engine check", "true")
    if err != nil {
        log.Fatalf("failed to submit CreateAsset transaction: %v", err)
    }

    c.JSON(http.StatusOK, gin.H{"message": "Ledger populated successfully"})
}

func decodeBase64(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func walletSignIn(c *gin.Context) {
    var requestBody map[string]string
    if err := c.BindJSON(&requestBody); err != nil {
        log.Printf("Failed to parse request body: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Extract and sanitize inputs
    encodedCert, certOk := requestBody["certificate"]
    encodedKey, keyOk := requestBody["privateKey"]
    mspContent, mspOk := requestBody["mspContent"]

    if !certOk || !keyOk || !mspOk {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields in the request body"})
        return
    }


    mspContent = strings.TrimSpace(mspContent)
    certificate, err := decodeBase64(encodedCert)
    privateKey, err := decodeBase64(encodedKey)

    // Create a new identity
    identity := NewX509Identity(mspContent, certificate, privateKey)

    // Close existing connections
    if clientConnection != nil {
        clientConnection.Close()
    }
    if gateway != nil {
        gateway.Close()
    }

    // Initialize wallet and store identity
    store := &FileWalletStore{}
    walletInstance, err := NewWallet(identity, store)
    if err != nil {
        log.Printf("Failed to create wallet: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
        return
    }

    err = walletInstance.Put("user_identity", identity)
    if err != nil {
        log.Printf("Failed to store identity in wallet: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store identity in wallet"})
        return
    }

    // Reconnect to Fabric with the new identity
    retrievedIdentity, err := walletInstance.Get("user_identity")
    if err != nil {
        log.Printf("Failed to retrieve identity from wallet: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve identity from wallet"})
        return
    }

    signingImplementation, err := retrievedIdentity.Signer()
    if err != nil {
        log.Printf("Failed to get signing implementation: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get signing implementation"})
        return
    }

    // Create a new gRPC connection
    clientConnection, err = updateGrpcConnection(mspContent)
    if err != nil {
        log.Printf("Failed to create gRPC connection: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create gRPC connection"})
        return
    }

    // Create a new Fabric gateway
    gateway, err = client.Connect(
        &retrievedIdentity,
        client.WithClientConnection(clientConnection),
        client.WithSign(signingImplementation),
    )
    if err != nil {
        log.Printf("Failed to create Fabric gateway: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Fabric gateway"})
        return
    }

    log.Printf("Reconnected to Fabric gateway successfully with new identity")
    c.JSON(http.StatusOK, gin.H{"message": "Reconnected to Fabric gateway successfully"})
}


func main() {
	defer clientConnection.Close()

	// Set up Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/read_asset/:key", readAsset)
	router.GET("/asset_history/:id", getAssetHistory)
	router.GET("/asset_exists/:id", assetExists)
	router.POST("/wallet_sign_in", walletSignIn)
	router.POST("/create_asset", createAsset)
	router.POST("/update_compliance", updateCompliance)
	// router.POST("/populate", populateLedger)	// ONLY USE FOR TESTING PURPOSES

	port := "8080"
	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
