package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	configPath  = "./connection-profile.yaml" // Path to the connection profile
	chaincodeID = "mycc"                      // Your chaincode ID
	channelID   = "mychannel"                 // Your channel name
)

var (
	sdk          *fabsdk.FabricSDK
	channelClient *channel.Client
)

// Initialize Fabric SDK and Channel Client
func initFabric() error {
	// Load connection profile
	sdkConfig := config.FromFile(configPath)
	var err error
	sdk, err = fabsdk.New(sdkConfig)
	if err != nil {
		return fmt.Errorf("failed to create SDK: %v", err)
	}

	// Create a channel client for invoking/querying chaincode
	channelClient, err = channel.New(sdk.ChannelContext(channelID, fabsdk.WithUser("admin")))
	if err != nil {
		return fmt.Errorf("failed to create channel client: %v", err)
	}

	return nil
}

// Query chaincode handler
func queryChaincode(c *gin.Context) {
	key := c.Param("key")

	// Query the chaincode
	response, err := channelClient.Query(channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         "query", // Your chaincode function
		Args:        [][]byte{[]byte(key)},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query chaincode: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": string(response.Payload)})
}

// Invoke chaincode handler
func invokeChaincode(c *gin.Context) {
	var request struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Invoke the chaincode
	_, err := channelClient.Execute(channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         "invoke", // Your chaincode function
		Args:        [][]byte{[]byte(request.Key), []byte(request.Value)},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to invoke chaincode: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoke successful"})
}

func main() {
	// Initialize Fabric SDK
	if err := initFabric(); err != nil {
		log.Fatalf("Failed to initialize Fabric SDK: %v", err)
	}
	defer sdk.Close()

	// Create a Gin router
	router := gin.Default()

	// Define routes
	router.GET("/query/:key", queryChaincode)
	router.POST("/invoke", invokeChaincode)

	// Start the server
	port := "8080" // Adjust the port as needed
	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
