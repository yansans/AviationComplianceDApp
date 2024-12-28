package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const aviationStackAPIURL = "https://api.aviationstack.com/v1/flights"

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

	// fmt.Println("Raw API Response: ", string(body))

	var flightDataResponse map[string]interface{}
	if err := json.Unmarshal(body, &flightDataResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	flightDetails, ok := flightDataResponse["data"].([]interface{})
	if !ok || len(flightDetails) == 0 {
		return nil, fmt.Errorf("no flight data found for flight ID %s", flightID)
	}

	flight := flightDetails[0].(map[string]interface{})
	flightData := &FlightData{
		FlightStatus:   flight["flight_status"].(string),
		Departure:      flight["departure"].(map[string]interface{})["estimated"].(string),
		Arrival:        flight["arrival"].(map[string]interface{})["estimated"].(string),
		FlightNumber:   flight["flight"].(map[string]interface{})["iata"].(string),
		AirlineName:    flight["airline"].(map[string]interface{})["name"].(string),
		AircraftType:   flight["aircraft"].(map[string]interface{})["iata"].(string),
		DepartureTime:  flight["departure"].(map[string]interface{})["estimated"].(string),
		ArrivalTime:    flight["arrival"].(map[string]interface{})["estimated"].(string),
		DepartureCity:  flight["departure"].(map[string]interface{})["airport"].(string),
		ArrivalCity:    flight["arrival"].(map[string]interface{})["airport"].(string),
	}

	return flightData, nil
}

func main() {
	flightID := "AA100" // Example Flight ID (American Airlines Flight 100)

	flightData, err := FetchFlightData(flightID)
	if err != nil {
		log.Fatalf("Error fetching flight data: %v", err)
	}

	fmt.Printf("Flight Status: %s\n", flightData.FlightStatus)
	fmt.Printf("Departure: %s at %s\n", flightData.DepartureCity, flightData.DepartureTime)
	fmt.Printf("Arrival: %s at %s\n", flightData.ArrivalCity, flightData.ArrivalTime)
	fmt.Printf("Flight Number: %s\n", flightData.FlightNumber)
	fmt.Printf("Airline: %s\n", flightData.AirlineName)
	fmt.Printf("Aircraft: %s\n", flightData.AircraftType)
}
