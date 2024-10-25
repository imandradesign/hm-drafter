package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Team struct {
	Name string `json:"name"`
	Tournament int `json:"tournament"`
}

type TeamApiResponse struct {
	Results []Team `json:"results"`
}

func GetTeams(tournamentID string) (teams [][]string) {
	log.Printf("Starting GetTeams Func. Tournament ID passed in: %v", tournamentID)
	api := fmt.Sprintf("https://kqhivemind.com/api/tournament/team/?tournament_id=%v&format=json", tournamentID)

	log.Printf("API call: %v", api)

	client := &http.Client{
    	Timeout: 10 * time.Second,
	}

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
	var teamApiResponse TeamApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamApiResponse); err != nil {
		log.Fatal(err)
	}

	log.Printf("teamApiResponse: %v", teamApiResponse)

	for _, team := range teamApiResponse.Results {
		log.Printf("Team retrieved: %v", team.Name)
		teams = append(teams, []string{team.Name})
	}

	log.Printf("TEAMS:\n%v", teams)
	return teams
}

func AddTeam(teamName string, tournamentID string) error {
	tournamentIDInt, err := strconv.Atoi(tournamentID)
	if err != nil {
		log.Print("Unable to convert tournament ID str to int in AddTeam() func.")
	}

	// Create the team struct to send to the API
	newTeam := Team{
		Name:       teamName,
		Tournament: tournamentIDInt,
	}

	// Convert the team struct to JSON
	teamJSON, err := json.Marshal(newTeam)
	if err != nil {
		return fmt.Errorf("error marshalling team data: %v", err)
	}

	api := fmt.Sprintf("https://kqhivemind.com/api/tournament/team/?tournament_id=%v&format=json", tournamentID)

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