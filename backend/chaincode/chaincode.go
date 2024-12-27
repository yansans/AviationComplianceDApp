package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode defines the chaincode structure
type SimpleChaincode struct{}

// Asset for Aviation Compliance
type Asset struct {
	ID          string `json:"id"`
	CompanyName string `json:"companyName"`
	AircraftID  string `json:"aircraftId"`
	Compliance  bool   `json:"compliance"`
	ReportDate  string `json:"reportDate"`
	Inspector   string `json:"inspector"`
	Description string `json:"description"`
}

// AssetHistory represents the history of an asset
type AssetHistory struct {
	Timestamp time.Time `json:"timestamp"`
	Asset     Asset     `json:"asset"`
}

// Init is called during chaincode instantiation to initialize the ledger
func (s *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	assets := []Asset{
		{ID: "asset1", CompanyName: "Airline A", AircraftID: "A123", Compliance: true, ReportDate: "2024-01-01", Inspector: "Inspector1", Description: "Routine Check"},
		{ID: "asset2", CompanyName: "Airline B", AircraftID: "B456", Compliance: false, ReportDate: "2024-02-15", Inspector: "Inspector2", Description: "Pending Maintenance"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to marshal asset: %s", err))
		}
		err = stub.PutState(asset.ID, assetJSON)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to add asset: %s", err))
		}
	}
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode
func (s *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	switch fn {
	case "CreateAsset":
		return s.CreateAsset(stub, args)
	case "ReadAsset":
		return s.ReadAsset(stub, args)
	case "UpdateCompliance":
		return s.UpdateCompliance(stub, args)
	case "GetHistory":
		return s.GetHistory(stub, args)
	default:
		return shim.Error("Invalid function name")
	}
}

// CreateAsset creates a new compliance report
func (s *SimpleChaincode) CreateAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	id := args[0]
	companyName := args[1]
	aircraftID := args[2]
	reportDate := args[3]
	inspector := args[4]
	description := args[5]
	compliance := args[6] == "true"

	exists, err := s.AssetExists(stub, id)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error checking asset existence: %s", err))
	}
	if exists {
		return shim.Error(fmt.Sprintf("Asset %s already exists", id))
	}

	asset := Asset{
		ID:          id,
		CompanyName: companyName,
		AircraftID:  aircraftID,
		Compliance:  compliance,
		ReportDate:  reportDate,
		Inspector:   inspector,
		Description: description,
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error marshalling asset: %s", err))
	}

	err = stub.PutState(id, assetJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error storing asset: %s", err))
	}

	return shim.Success(nil)
}

// ReadAsset retrieves an asset from the ledger
func (s *SimpleChaincode) ReadAsset(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := args[0]
	assetJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to read asset: %s", err))
	}
	if assetJSON == nil {
		return shim.Error(fmt.Sprintf("Asset %s does not exist", id))
	}

	return shim.Success(assetJSON)
}

// UpdateCompliance updates the compliance status of an asset
func (s *SimpleChaincode) UpdateCompliance(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	id := args[0]
	compliance := args[1] == "true"

	assetJSON, err := stub.GetState(id)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to read asset: %s", err))
	}
	if assetJSON == nil {
		return shim.Error(fmt.Sprintf("Asset %s does not exist", id))
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to unmarshal asset: %s", err))
	}

	asset.Compliance = compliance
	assetJSON, err = json.Marshal(asset)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to marshal updated asset: %s", err))
	}

	err = stub.PutState(id, assetJSON)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to store updated asset: %s", err))
	}

	return shim.Success(nil)
}

// GetHistory retrieves the history of an asset
func (s *SimpleChaincode) GetHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := args[0]
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to retrieve history: %s", err))
	}
	defer resultsIterator.Close()

	var history []AssetHistory
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(fmt.Sprintf("Error iterating history: %s", err))
		}

		var asset Asset
		if response.Value != nil {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return shim.Error(fmt.Sprintf("Failed to unmarshal asset: %s", err))
			}
		}

		history = append(history, AssetHistory{
			Timestamp: time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)),
			Asset:     asset,
		})
	}

	historyJSON, err := json.Marshal(history)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to marshal history: %s", err))
	}

	return shim.Success(historyJSON)
}

// AssetExists checks if an asset exists
func (s *SimpleChaincode) AssetExists(stub shim.ChaincodeStubInterface, id string) (bool, error) {
	assetJSON, err := stub.GetState(id)
	if err != nil {
		return false, err
	}

	return assetJSON != nil, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
