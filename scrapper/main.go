package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Define structs to match your JSON structure
type EvacuationSite struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Contacts    []string `json:"contacts"`
	Capacity    int      `json:"capacity"`
	Visibility  string   `json:"visibility"`
	HazardType  string   `json:"hazardType"`
	Lat         float64  `json:"lat"`
	Lng         float64  `json:"lng"`
	Images      []string `json:"images"`
	CreatedAt   string   `json:"createdAt"`
}

type Response struct {
	Success bool             `json:"success"`
	Data    []EvacuationSite `json:"data"`
}

func fetchEvacuationSites(apiURL string) {
	// Make HTTP GET reques
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatal("Error fetching URL: ", err)
	}
	defer resp.Body.Close()

	// Check if status code is OK
	if resp.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body: ", err)
	}

	// Parse JSON into struct
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal("Error parsing JSON: ", err)
	}


	// Check if the request was successful
	if !response.Success {
		log.Fatal("API request was not successful")
	}

	// Print the parsed data
	fmt.Println("Evacuation Sites:")
	for i, site := range response.Data {
		fmt.Printf("\nSite %d:\n", i+1)
		fmt.Printf("Name: %s\n", site.Name)
		fmt.Printf("Description: %s\n", site.Description)
		fmt.Printf("Contacts: %v\n", site.Contacts)
		fmt.Printf("Capacity: %d\n", site.Capacity)
		fmt.Printf("Hazard Type: %s\n", site.HazardType)
		fmt.Printf("Coordinates: (%f, %f)\n", site.Lat, site.Lng)
		fmt.Printf("Images: %v\n", site.Images)
		fmt.Printf("Created At: %s\n", site.CreatedAt)
	}
}

func main() {
	// Replace with your actual API endpoint URL
	apiURL := "https://admin-evacu-ease.vercel.app/api/locations"
	fmt.Printf("Fetching data from %s...\n", apiURL)
	fetchEvacuationSites(apiURL)
}
