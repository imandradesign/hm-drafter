package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type FormFields struct {
	FieldName        string    `json:"field_name"`
	FieldSlug      string `json:"field_slug"`
	FieldDescription      string `json:"field_description"`
}

// GetFormFields takes the API string that lists all tourney form fields, checks the `results` list entries and returns the form fields with their randomly assigned name
func GetFormFields(tournamentId string) (fields [][]string) {
	log.Println("Fetching form field data...")
	client := &http.Client{
    	Timeout: 10 * time.Second,
	}

	api := fmt.Sprintf(apiTemplate, "player-info-field", fmt.Sprintf("&tournament_id=%v", tournamentId), "")
		
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Decode the response into APIResponse struct
	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Fatal(err)
	}

	for _, field := range apiResponse.FResults {
		fields = append(fields, []string{
			field.FieldName,
			field.FieldSlug,
			field.FieldDescription,
		})
	}

	return fields
}