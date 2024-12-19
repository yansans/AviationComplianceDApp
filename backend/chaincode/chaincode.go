package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SimpleChaincode defines the chaincode structure
type SimpleChaincode struct {
	contractapi.Contract
}

// Asset represents a simple asset
type Asset struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// InitLedger initializes the ledger with sample data
func (s *SimpleChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Name: "Asset One", Value: 100},
		{ID: "asset2", Name: "Asset Two", Value: 200},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to add asset: %v", err)
		}
	}
	return nil
}

// CreateAsset creates a new asset
func (s *SimpleChaincode) CreateAsset(ctx contractapi.TransactionContextInterface, id string, name string, value int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("asset %s already exists", id)
	}

	asset := Asset{
		ID:    id,
		Name:  name,
		Value: value,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset reads an asset from the ledger
func (s *SimpleChaincode) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// AssetExists checks if an asset exists
func (s *SimpleChaincode) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to check if asset exists: %v", err)
	}

	return assetJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
