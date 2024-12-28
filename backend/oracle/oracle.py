import os
import requests
from dotenv import load_dotenv

load_dotenv()

AVIATION_STACK_API_URL = "https://api.aviationstack.com/v1/flights"

class FlightData:
    def __init__(self, flight_status, departure, arrival, flight_number, airline_name,
                 aircraft_type, departure_time, arrival_time, departure_city, arrival_city):
        self.flight_status = flight_status
        self.departure = departure
        self.arrival = arrival
        self.flight_number = flight_number
        self.airline_name = airline_name
        self.aircraft_type = aircraft_type
        self.departure_time = departure_time
        self.arrival_time = arrival_time
        self.departure_city = departure_city
        self.arrival_city = arrival_city

def fetch_flight_data(flight_id):
    api_key = os.getenv("AVIATION_STACK_API_KEY")
    if not api_key:
        raise ValueError("API key is missing")

    url = f"{AVIATION_STACK_API_URL}?access_key={api_key}&flight_iata={flight_id}"

    response = requests.get(url)
    if response.status_code != 200:
        raise Exception(f"API request failed with status code {response.status_code}")

    flight_data_response = response.json()
    flight_details = flight_data_response.get("data", [])

    if not flight_details:
        raise ValueError(f"No flight data found for flight ID {flight_id}")

    flight = flight_details[0]

    flight_data = FlightData(
        flight_status=flight["flight_status"],
        departure=flight["departure"]["estimated"],
        arrival=flight["arrival"]["estimated"],
        flight_number=flight["flight"]["iata"],
        airline_name=flight["airline"]["name"],
        aircraft_type=flight["aircraft"]["iata"],
        departure_time=flight["departure"]["estimated"],
        arrival_time=flight["arrival"]["estimated"],
        departure_city=flight["departure"]["airport"],
        arrival_city=flight["arrival"]["airport"]
    )

    return flight_data

def handle_oracle_request(request_data):
    flight_id = request_data.get("data")[0]

    try:
        flight_data = fetch_flight_data(flight_id)

        result = {
            "flight_status": flight_data.flight_status,
            "departure_city": flight_data.departure_city,
            "arrival_city": flight_data.arrival_city,
            "flight_number": flight_data.flight_number,
            "airline_name": flight_data.airline_name,
            "aircraft_type": flight_data.aircraft_type,
            "departure_time": flight_data.departure_time,
            "arrival_time": flight_data.arrival_time
        }

        response = {
            "result": result
        }

        return response

    except Exception as e:
        return {"error": str(e)}

request_data = {
    "data": ["AA100"]  # Example flight ID (American Airlines Flight 100)
}

response = handle_oracle_request(request_data)
print(response)
