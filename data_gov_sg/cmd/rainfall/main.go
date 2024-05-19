package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type RainfallData struct {
	Metadata struct {
		Stations []struct {
			ID       string `json:"id"`
			DeviceID string `json:"device_id"`
			Name     string `json:"name"`
			Location struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"location"`
		} `json:"stations"`
		ReadingType string `json:"reading_type"`
		ReadingUnit string `json:"reading_unit"`
	} `json:"metadata"`
	Items []struct {
		Timestamp time.Time `json:"timestamp"`
		Readings  []struct {
			StationID string          `json:"station_id"`
			Value     decimal.Decimal `json:"value"`
		} `json:"readings"`
	} `json:"items"`
	APIInfo struct {
		Status string `json:"status"`
	} `json:"api_info"`
}

const (
	url string = "https://api.data.gov.sg/v1/environment/rainfall"
)

func main() {
	file, writer := createFile("rainfall.csv")
	defer file.Close()
	defer writer.Flush()

	var startTime, snapshotAt, now time.Time
	startTime = time.Date(2024, 05, 12, 0, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	now = time.Now()

	// Create a ticker to trigger 2 requests per second
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Channel to signal when to make a request
	requests := make(chan time.Time)

	// Function to make an HTTP request
	makeRequest := func(timestamp time.Time) {
		defer wg.Done()
		// Format the timestamp
		timestampStr := timestamp.Format("2006-01-02")
		fmt.Println("Making request for timestamp:", timestampStr)

		// Make the HTTP request
		url := fmt.Sprintf("%s?date=%s", url, timestampStr)
		snapshotAt = time.Now()
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error making HTTP request for %s: %v", timestampStr, err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body for %s: %v", timestampStr, err)
			return
		}

		if resp.StatusCode == 200 {
			var rainfall RainfallData
			err = json.Unmarshal(body, &rainfall)
			if err != nil {
				log.Printf("Error unmarshalling response body for %s: %v", timestampStr, err)
				return
			}

			writeToCSV(snapshotAt, &rainfall, writer)
		} else {
			log.Printf("Unexpected response status %d for %s. Exiting.", resp.StatusCode, timestampStr)
			close(requests)
			return
		}
	}

	// Launch a goroutine to send timestamps through the channel
	go func() {
		for t := startTime; t.Before(now); t = t.Add(24 * time.Hour) {
			requests <- t
		}
		close(requests)
	}()

	// Launch goroutines to process requests
	for t := range requests {
		<-ticker.C
		wg.Add(1)
		go makeRequest(t)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func createFile(file_name string) (*os.File, *csv.Writer) {
	// Create CSV file
	file, err := os.Create(file_name)
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}

	writer := csv.NewWriter(file)

	// Write CSV header
	header := []string{"reading_type", "reading_unit", "station_id", "rain_value", "measured_at", "snapshot_at"}
	writer.Write(header)

	return file, writer
}

func writeToCSV(snapshotAt time.Time, rainfallData *RainfallData, writer *csv.Writer) {
	for _, item := range rainfallData.Items {
		measuredAt := item.Timestamp
		for _, reading := range item.Readings {
			record := []string{
				rainfallData.Metadata.ReadingType,
				rainfallData.Metadata.ReadingUnit,
				reading.StationID,
				fmt.Sprintf("%v", reading.Value),
				measuredAt.Format(time.RFC3339),
				snapshotAt.Format(time.RFC3339),
			}
			writer.Write(record)
		}
	}
}
