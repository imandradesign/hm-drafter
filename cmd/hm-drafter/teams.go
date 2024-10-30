package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Team struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Tournament int    `json:"tournament"`
}

type TeamInfo struct {
	ID      int
	Name    string
	Players []Players `json:"players,omitempty"`
}

type TeamApiResponse struct {
	Results []Team `json:"results"`
}

func GetTeams(tournamentID string, players []Players) (teams []TeamInfo) {
	log.Printf("Starting GetTeams Func. Tournament ID passed in: %v", tournamentID)

	// Fetch Teams from the API
	api := fmt.Sprintf("https://kqhivemind.com/api/tournament/team/?tournament_id=%v&format=json", tournamentID)
	client, req := createRequest("GET", api, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var teamApiResponse TeamApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamApiResponse); err != nil {
		log.Fatal(err)
	}

	// Create a map of team IDs to TeamInfo pointers for quick access
	teamMap := make(map[int]*TeamInfo)
	for _, team := range teamApiResponse.Results {
		teamInfo := TeamInfo{ID: team.ID, Name: team.Name, Players: []Players{}}
		teams = append(teams, teamInfo)
		teamMap[team.ID] = &teams[len(teams)-1] // Map each TeamInfo by its ID
	}

	// Iterate over players and add them to the matching team in teamMap
	for _, player := range players {
		if team, found := teamMap[player.Team]; found {
			team.Players = append(team.Players, player)
		}
	}

	log.Printf("TEAMS with Players:\n%v", teams)
	return teams
}


func AddNewTeam(teamName string, tournamentID string) {
	// Convert tournament ID to an integer
	tournamentIDInt, err := strconv.Atoi(tournamentID)
	if err != nil {
		log.Print("Unable to convert tournament ID string to int in AddTeam() func.")
	}

	// Create the team struct to send to the API
	newTeam := Team{
		Name:       teamName,
		Tournament: tournamentIDInt,
	}

	// Convert the team struct to JSON
	teamJSON, err := json.Marshal(newTeam)
	if err != nil {
		log.Fatalf("error marshalling team data: %v", err)
	}

	api := fmt.Sprintf("https://kqhivemind.com/api/tournament/team/?tournament_id=%v&format=json", tournamentID)

	client, req := createRequest("POST", api, bytes.NewBuffer(teamJSON))
	req.Header.Set("Content-Type", "application/json")

	// Make the POST request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to add team. Status: %v, Response: %s", resp.Status, string(body))
	}

	log.Printf("Request to add team returned status: %v", resp.Status)
}


func DeleteTeam(teamID string, tournamentID string) {
	// Convert team ID to an integer
	teamIDInt, err := strconv.Atoi(teamID)
	if err != nil {
		log.Print("Unable to convert team ID string to int in DeleteTeam() func.")
		return
	}

	// Create a payload that matches the API's requirements for deletion
	deletePayload := map[string]interface{}{
		"id":       teamIDInt,
		"action":   "delete", // Adjust this based on the API's expected fields for deletion
		"tournament": tournamentID,
	}

	// Convert the payload to JSON
	deleteJSON, err := json.Marshal(deletePayload)
	if err != nil {
		log.Fatalf("error marshalling team delete data: %v", err)
	}

	// Use POST method for the delete request
	api := fmt.Sprintf("https://kqhivemind.com/api/tournament/team/?tournament_id=%v&format=json", tournamentID)
	client, req := createRequest("POST", api, bytes.NewBuffer(deleteJSON))
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making team DELETE request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Failed to delete team. Status: %v, Response: %s", resp.Status, string(body))
	}

	log.Printf("Request to delete team ID %d returned status: %v", teamIDInt, resp.Status)
}
