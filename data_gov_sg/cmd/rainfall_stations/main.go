package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/shopspring/decimal"
)

type RainfallStation struct {
	Stations []struct {
		ID       string `json:"id"`
		DeviceID string `json:"device_id"`
		Name     string `json:"name"`
		Location struct {
			Latitude  decimal.Decimal `json:"latitude"`
			Longitude decimal.Decimal `json:"longitude"`
		} `json:"location"`
	} `json:"stations"`
}

const (
	outputFilePath string = "rainfall_stations.csv"
	inputFilePath  string = "cmd/rainfall_stations/stations.json"
)

func main() {
	// Read the JSON file
	file, err := os.ReadFile(inputFilePath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Unmarshal the JSON data
	var stationsData RainfallStation
	if err := json.Unmarshal(file, &stationsData); err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// Create CSV file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer outputFile.Close()

	// Create a CSV writer
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Write the CSV header
	header := []string{"id", "device_id", "name", "location"}
	writer.Write(header)

	// Write the station data
	for _, station := range stationsData.Stations {
		location := fmt.Sprintf("(%s,%s)", station.Location.Latitude.String(), station.Location.Longitude.String())
		record := []string{station.ID, station.DeviceID, station.Name, location}
		writer.Write(record)
	}

	fmt.Println("CSV file created successfully.")
}
