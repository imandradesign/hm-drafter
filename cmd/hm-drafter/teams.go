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

func GetTeams(tournamentID string, players []Player) (teams []TeamInfo) {
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
		teamInfo := TeamInfo{ID: team.ID, Name: team.Name, Players: []Player{}}
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


func DeleteTeam(teamID string, teamName string, tournamentID string) (playerIDs []string) {
	// Construct the API URL to retrieve team information
	apiGet := fmt.Sprintf("https://kqhivemind.com/api/tournament/team/%s/?tournament_id=%v", teamID, tournamentID)

	// Create a GET request to retrieve the team information
	client, reqGet := createRequest("GET", apiGet, nil)
	reqGet.Header.Set("Content-Type", "application/json")

	// Make the GET request to fetch team data
	respGet, err := client.Do(reqGet)
	if err != nil {
		log.Fatalf("error making team GET request: %v", err)
	}
	defer respGet.Body.Close()

	// Check if GET request was successful
	if respGet.StatusCode == http.StatusOK {
		var teamData struct {
			Players []struct {
				ID string `json:"id"`
			} `json:"players"`
		}

		// Parse the JSON response to get the player IDs
		err = json.NewDecoder(respGet.Body).Decode(&teamData)
		if err != nil {
			log.Fatalf("error decoding team data: %v", err)
		}

		// Collect all player IDs from the team
		for _, player := range teamData.Players {
			playerIDs = append(playerIDs, player.ID)
		}

		log.Printf("Players in team %s (ID: %s): %v", teamName, teamID, playerIDs)
	} else {
		log.Printf("Failed to fetch team details. Status: %v", respGet.Status)
	}

	// Now, construct the API URL for the DELETE request
	apiDelete := fmt.Sprintf("https://kqhivemind.com/api/tournament/team/%s/?tournament_id=%v", teamID, tournamentID)
	client, reqDelete := createRequest("DELETE", apiDelete, nil)
	reqDelete.Header.Set("Content-Type", "application/json")

	// Make the DELETE request to delete the team
	respDelete, err := client.Do(reqDelete)
	if err != nil {
		log.Fatalf("error making team DELETE request: %v", err)
	}
	defer respDelete.Body.Close()

	// Check the response status code for DELETE
	if respDelete.StatusCode != http.StatusOK && respDelete.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(respDelete.Body)
		log.Fatalf("Failed to delete team. Status: %v, Response: %s", respDelete.Status, string(body))
	} else {
		log.Printf("Team ID %s deleted successfully with status: %v", teamID, respDelete.Status)
	}

	// Return the list of player IDs that were in the team before deletion
	return playerIDs
}


func GetTeamNameByID(teams []TeamInfo, teamID string) string {
	teamIDInt, err := strconv.Atoi(teamID)
	if err != nil {
		log.Printf("Can't convert team ID (%v) str to int", teamID)
	}

    for _, team := range teams {
        if team.ID == teamIDInt {
            return team.Name
        }
    }
    return ""
}
