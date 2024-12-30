package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// configPath   = "../../fabric/test-network/config/configtx.yaml"
	chaincodeID  = "basic"                                                    
	channelID    = "channel1"                                                 
	tlsCertPath  = "../../fabric/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem"                                  // TLS certificate path
	mspID        = "Org1MSP"
)
type Asset struct {
	ID          string `json:"id"`
	CompanyName string `json:"company_name"`
	AircraftID  string `json:"aircraft_id"`
	ReportDate  string `json:"report_date"`
	Inspector   string `json:"inspector"`
	Description string `json:"description"`
	Compliance  bool   `json:"compliance"`
}
type AssetHistory struct {
	Timestamp time.Time `json:"timestamp"`
	Asset     Asset     `json:"asset"`
}

var gateway *client.Gateway
var clientConnection *grpc.ClientConn

func newGrpcConnection() (*grpc.ClientConn, error) {
	// With TLS
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
    connection, err := grpc.Dial("localhost:7051", grpc.WithTransportCredentials(transportCredentials))
    if err != nil {
        return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
    }

    return connection, nil
}

func initFabric() error {
	var err error
	
	// Create gRPC connection
	clientConnection, err = newGrpcConnection()
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %v", err)
	}

	// Init wallet and Load Identity
	store := &FileWalletStore{}
	identity, err := LoadIdentityFromFiles(mspID, 
		"../../fabric/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem", 
		"../../fabric/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk")
	if err != nil {
		return fmt.Errorf("failed to load identity from files: %v", err)
	}

	// Create a wallet and add the identity
	walletInstance, err := NewWallet(identity, store)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}

	// Store the identity in the wallet with a label
	err = walletInstance.Put("user_identity", identity)
	if err != nil {
		return fmt.Errorf("failed to store identity in wallet: %v", err)
	}

	retrievedIdentity, err := walletInstance.Get("user_identity")
	if err != nil {
		return fmt.Errorf("failed to retrieve identity from wallet: %v", err)
	}

	signingImplementation, err := retrievedIdentity.Signer()
	if err != nil {
		panic(fmt.Sprintf("failed to get signing implementation: %v", err))
	}

	// Create the gateway using the wallet identity
	gateway, err = client.Connect(
		&retrievedIdentity,
		client.WithClientConnection(clientConnection),
		client.WithSign(signingImplementation))
	if err != nil {
		return fmt.Errorf("failed to create gateway: %v", err)
	}

	return nil
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
		CompanyName string `json:"company_name"`
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

	contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

	_, err := contract.SubmitTransaction(
		"CreateAsset", 
		request.ID, 
		request.CompanyName, 
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

	c.JSON(http.StatusOK, gin.H{"message": "Invoke successful"})
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
        c.JSON(http.StatusOK, gin.H{"message": "Asset compliance updated successfully"})
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

func main() {
	defer clientConnection.Close()

	// Init Fabric connection
	if err := initFabric(); err != nil {
		log.Fatalf("Failed to initialize Fabric Gateway: %v", err)
	}

	// Set up Gin router
	router := gin.Default()

	router.GET("/read_asset/:key", readAsset)
	router.GET("/asset_history/:id", getAssetHistory)
	router.GET("/asset_exists/:id", assetExists)
	router.POST("/create_asset", createAsset)
	router.POST("/update_compliance", updateCompliance)
	// router.POST("/populate", populateLedger)	// ONLY USE FOR TESTING PURPOSES

	port := "8080"
	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
