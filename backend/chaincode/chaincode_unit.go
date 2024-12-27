package main

import (
	"encoding/json"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"github.com/stretchr/testify/assert"
)

// TestCreateAsset tests the CreateAsset function
func TestCreateAsset(t *testing.T) {
	chaincode := new(SimpleChaincode)
	mockStub := shimtest.NewMockStub("mockStub", chaincode)

	// Case 1: Asset creation success
	assetID := "asset1"
	args := [][]byte{
		[]byte("CreateAsset"),
		[]byte(assetID), []byte("Airline X"), []byte("A789"),
		[]byte("2024-12-01"), []byte("Inspector Y"), []byte("Passed Safety Check"),
		[]byte("true"),
	}
	response := mockStub.MockInvoke("1", args)

	assert.Equal(t, int32(shim.OK), response.Status, "Expected CreateAsset to succeed")
	assert.Empty(t, response.Message, "Expected no error message")

	// Verify asset in state
	state := mockStub.State[assetID]
	assert.NotNil(t, state, "Expected asset to be stored")
	var asset Asset
	err := json.Unmarshal(state, &asset)
	assert.NoError(t, err, "Expected unmarshalling asset to succeed")
	assert.Equal(t, assetID, asset.ID)
}

// TestReadAsset tests the ReadAsset function
func TestReadAsset(t *testing.T) {
	chaincode := new(SimpleChaincode)
	mockStub := shimtest.NewMockStub("mockStub", chaincode)

	// Initialize ledger with default assets
	mockStub.MockInit("1", [][]byte{[]byte("Init")})

	// Case 1: Read existing asset
	response := mockStub.MockInvoke("1", [][]byte{
		[]byte("ReadAsset"), []byte("asset1"),
	})
	assert.Equal(t, int32(shim.OK), response.Status, "Expected ReadAsset to succeed")

	var asset Asset
	err := json.Unmarshal(response.Payload, &asset)
	assert.NoError(t, err, "Expected unmarshalling asset to succeed")
	assert.Equal(t, "asset1", asset.ID)

	// Case 2: Read non-existent asset
	response = mockStub.MockInvoke("2", [][]byte{
		[]byte("ReadAsset"), []byte("nonexistent"),
	})
	assert.NotEqual(t, int32(shim.OK), response.Status, "Expected ReadAsset to fail for non-existent asset")
}

// TestUpdateCompliance tests the UpdateCompliance function
func TestUpdateCompliance(t *testing.T) {
	chaincode := new(SimpleChaincode)
	mockStub := shimtest.NewMockStub("mockStub", chaincode)

	// Initialize ledger with default assets
	mockStub.MockInit("1", [][]byte{[]byte("Init")})

	// Case 1: Update compliance status
	assetID := "asset1"
	args := [][]byte{[]byte("UpdateCompliance"), []byte(assetID), []byte("false")}
	response := mockStub.MockInvoke("1", args)

	assert.Equal(t, int32(shim.OK), response.Status, "Expected UpdateCompliance to succeed")
	assert.Empty(t, response.Message, "Expected no error message")

	// Verify compliance update
	state := mockStub.State[assetID]
	assert.NotNil(t, state, "Expected asset to be present")
	var asset Asset
	err := json.Unmarshal(state, &asset)
	assert.NoError(t, err, "Expected unmarshalling asset to succeed")
	assert.False(t, asset.Compliance)
}

// TestGetHistory tests the GetHistory function
func TestGetHistory(t *testing.T) {
	chaincode := new(SimpleChaincode)
	mockStub := shimtest.NewMockStub("mockStub", chaincode)

	// Initialize ledger with default assets
	mockStub.MockInit("1", [][]byte{[]byte("Init")})

	// Update compliance status to create history
	mockStub.MockInvoke("1", [][]byte{[]byte("UpdateCompliance"), []byte("asset1"), []byte("false")})
	mockStub.MockInvoke("2", [][]byte{[]byte("UpdateCompliance"), []byte("asset1"), []byte("true")})

	// Retrieve history for asset1
	response := mockStub.MockInvoke("3", [][]byte{[]byte("GetHistory"), []byte("asset1")})

	assert.Equal(t, int32(shim.OK), response.Status, "Expected GetHistory to succeed")

	var history []AssetHistory
	err := json.Unmarshal(response.Payload, &history)
	assert.NoError(t, err, "Expected unmarshalling history to succeed")
	assert.True(t, len(history) >= 2, "Expected at least 2 history entries")
}
