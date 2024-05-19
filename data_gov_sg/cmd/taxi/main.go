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
)

const (
	url string = "https://api.data.gov.sg/v1/transport/taxi-availability"
)

type TaxiAvailability struct {
	Type string `json:"type"`
	Crs  struct {
		Type       string `json:"type"`
		Properties struct {
			Href string `json:"href"`
			Type string `json:"type"`
		} `json:"properties"`
	} `json:"crs"`
	Features []struct {
		Type     string `json:"type"`
		Geometry struct {
			Type        string      `json:"type"`
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Timestamp time.Time `json:"timestamp"`
			TaxiCount int       `json:"taxi_count"`
			APIInfo   struct {
				Status string `json:"status"`
			} `json:"api_info"`
		} `json:"properties"`
	} `json:"features"`
}

func main() {
	file, writer := createFile("taxi_availability.csv")
	defer file.Close()
	defer writer.Flush()

	var startTime, snapshotAt, now time.Time

	// Date of dataset availability
	// startTime = time.Date(2016, 4, 8, 0, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
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
		timestampStr := timestamp.Format("2006-01-02T15:04:05")
		fmt.Println("Making request for timestamp:", timestampStr)

		// Make the HTTP request
		url := fmt.Sprintf("%s?datetime=%s", url, timestampStr)
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
			var taxiAvailability TaxiAvailability
			err = json.Unmarshal(body, &taxiAvailability)
			if err != nil {
				log.Printf("Error unmarshalling response body for %s: %v", timestampStr, err)
				return
			}

			writeToCSV(snapshotAt, &taxiAvailability, writer)
		} else {
			log.Printf("Unexpected response status %d for %s. Exiting.", resp.StatusCode, timestampStr)
			close(requests)
			return
		}
	}

	// Launch a goroutine to send timestamps through the channel
	go func() {
		for t := startTime; t.Before(now); t = t.Add(5 * time.Minute) {
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
	header := []string{"type", "crs_type", "crs_properties_href", "crs_properties_type", "feature_type", "geometry_type", "coordinates", "timestamp", "taxi_count", "api_info_status", "snapshot_at"}
	writer.Write(header)

	return file, writer
}

func writeToCSV(snapshotAt time.Time, taxiAvailability *TaxiAvailability, writer *csv.Writer) {
	for _, feature := range taxiAvailability.Features {
		// Flatten coordinates into a string
		coordinates := ""
		for _, coord := range feature.Geometry.Coordinates {
			coordinates += fmt.Sprintf("(%f,%f) ", coord[0], coord[1])
		}

		record := []string{
			taxiAvailability.Type,
			taxiAvailability.Crs.Type,
			taxiAvailability.Crs.Properties.Href,
			taxiAvailability.Crs.Properties.Type,
			feature.Type,
			feature.Geometry.Type,
			coordinates,
			feature.Properties.Timestamp.Format(time.RFC3339),
			fmt.Sprintf("%d", feature.Properties.TaxiCount),
			feature.Properties.APIInfo.Status,
			snapshotAt.Format(time.RFC3339),
		}
		writer.Write(record)
	}
}
