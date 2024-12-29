package main

import (
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"crypto/x509"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	configPath   = "../../fabric/test-network/config/connection-profile.yaml"
	chaincodeID  = "basic"                                                    
	channelID    = "channel1"                                                 
	tlsCertPath  = "../../fabric/test-network/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem"                                  // TLS certificate path
	mspID        = "Org1MSP"
)

var gateway *client.Gateway

func newGrpcConnection() (*grpc.ClientConn, error) {
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
	// Create gRPC connection
	clientConnection, err := newGrpcConnection()
	if err != nil {
		return fmt.Errorf("failed to create gRPC connection: %v", err)
	}
	defer clientConnection.Close()

	// Init wallet and Load Identity
	store := &FileWalletStore{}
	identity, err := LoadIdentityFromFiles(mspID, "./key/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem", "./key/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk")
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

	// Create the gateway using the wallet identity
	gateway, err = client.Connect(
		&retrievedIdentity,
		client.WithClientConnection(clientConnection),
		// client.WithConnectionProfile(configPath),
	)
	if err != nil {
		return fmt.Errorf("failed to create gateway: %v", err)
	}

	return nil
}

// func queryChaincode(c *gin.Context) {
// 	key := c.Param("key")

// 	// Get the contract from the gateway
// 	contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

// 	response, err := contract.EvaluateTransaction("query", []string{key})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to query chaincode: %v", err)})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"result": string(response)})
// }

// func invokeChaincode(c *gin.Context) {
// 	var request struct {
// 		Key   string `json:"key"`
// 		Value string `json:"value"`
// 	}

// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	// Get the contract from the gateway
// 	contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

// 	_, err := contract.SubmitTransaction("invoke", []string{request.Key, request.Value})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to invoke chaincode: %v", err)})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Invoke successful"})
// }

func populateLedger() {
	// Get the contract from the gateway
	contract := gateway.GetNetwork(channelID).GetContract(chaincodeID)

	// Submit the transaction to initialize the ledger
	_, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("failed to submit InitLedger transaction: %v", err)
	}

	log.Println("Ledger initialized successfully")
}

func main() {
	// Init Fabric connection
	if err := initFabric(); err != nil {
		log.Fatalf("Failed to initialize Fabric Gateway: %v", err)
	}

	// Init ledger with sample assets
	// populateLedger()

	// Set up Gin router
	router := gin.Default()

	// router.GET("/query/:key", queryChaincode)
	// router.POST("/invoke", invokeChaincode)

	port := "8080"
	log.Printf("Server is running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
