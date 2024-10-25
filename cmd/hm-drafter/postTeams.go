package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Team struct {
	Name string `json:"name"`
	Tournament int `json:"tournament"`
}

func AddTeam(teamName string, tournamentID int) error {
	// Create the team struct to send to the API
	newTeam := Team{
		Name:       teamName,
		Tournament: tournamentID,
	}

	// Convert the team struct to JSON
	teamJSON, err := json.Marshal(newTeam)
	if err != nil {
		return fmt.Errorf("error marshalling team data: %v", err)
	}

	api := fmt.Sprintf("https://kqhivemind.com/api/tournament/?tournament_id=%v", tournamentID)

	// Make the POST request to the API
	resp, err := http.Post(api, "application/json", bytes.NewBuffer(teamJSON))
	if err != nil {
		return fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Failed to add team. Status: %v, Response: %s", resp.Status, string(body))
	}

	return nil
}