package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Players struct {
	Name       string            `json:"name"`
	ID         float64           `json:"user"`
	Scene      string            `json:"scene"`
	Pronouns   string            `json:"pronouns"`
	Team       int               `json:"team"`
	FormFields map[string]string `json:"form_fields,omitempty"`
}

type PlayersApiResponse struct {
	Results []Players `json:"results"`
}

func ParsePlayers(data map[string]interface{}) Players {
	// Extract static fields like name, scene, pronouns, and ID
	player := Players{
		Name:       safeString(data["name"]),
		Scene:      safeString(data["scene"]),
		Pronouns:   safeString(data["pronouns"]),
		ID:         data["user"].(float64),
		Team:       int(data["team"].(float64)),
		FormFields: make(map[string]string),
	}

	// Process dynamic form fields if they exist
	for _, field := range formFields {
		fieldName := field[0]
		if value, ok := data[fieldName]; ok {
			switch v := value.(type) {
			case string:
				player.FormFields[field[1]] = v
			case []interface{}:
				var strValues []string
				for _, item := range v {
					strValues = append(strValues, fmt.Sprintf("%v", item))
				}
				player.FormFields[field[1]] = strings.Join(strValues, ", ")
			default:
				player.FormFields[field[1]] = fmt.Sprintf("%v", value)
			}
		}
	}

	return player
}

// Helper function to safely convert a field to a string, handling nil cases
func safeString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// GetPlayersData retrieves all player data for the specified tournament ID, returning a slice of Players
func GetPlayersData(tournamentID string) (players []Players) {
	log.Println("Fetching player data...")

	page := 1
	for {
		// Build the API URL for the current page
		api := fmt.Sprintf(apiTemplate, "player", fmt.Sprintf("&tournament_id=%v", tournamentID), fmt.Sprintf("&page=%d", page))
		client, req := createRequest("GET", api, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Parse the response into a PlayersApiResponse
		var playerApiResponse PlayersApiResponse
		err = json.NewDecoder(resp.Body).Decode(&playerApiResponse)
		if err != nil {
			log.Fatal(err)
		}

		// If no results are returned, exit the loop
		if len(playerApiResponse.Results) == 0 {
			log.Println("No more pages.")
			break
		}

		// Add the current page's players to the slice
		players = append(players, playerApiResponse.Results...)

		// Increment the page for the next iteration
		page++
	}

	log.Printf("API data fetched.\nPLAYERS:\n%v", players)
	return players
}
