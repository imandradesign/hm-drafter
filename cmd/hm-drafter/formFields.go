package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type FormFields struct {
	FieldName        string    `json:"field_name"`
	FieldSlug      string `json:"field_slug"`
	FieldDescription      string `json:"field_description"`
}

type FormApiResponse struct {
	Results []FormFields `json:"results"`
}

// GetFormFields takes the API string that lists all tourney form fields, checks the `results` list entries and returns the form fields with their randomly assigned name
func GetFormFields(tournamentId string) (fields [][]string) {
	log.Println("Fetching form field data...")

	api := fmt.Sprintf(apiTemplate, "player-info-field", fmt.Sprintf("&tournament_id=%v", tournamentId), "")
		
	client, req := createRequest("GET", api, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Decode the response into APIResponse struct
	var formApiResponse FormApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&formApiResponse); err != nil {
		log.Fatal(err)
	}

	for _, field := range formApiResponse.Results {
		fields = append(fields, []string{
			field.FieldName,
			field.FieldSlug,
			field.FieldDescription,
		})
	}

	log.Println("API data fetched.")
	return fields
}