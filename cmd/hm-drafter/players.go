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
	Image      string            `json:"image"`
	FormFields map[string]string `json:"form_fields,omitempty"`
}

type PlayersApiResponse struct {
	Results []map[string]interface{} `json:"results"`
}

func ParsePlayers(data map[string]interface{}) Players {
	player := Players{
		ID:         data["user"].(float64),
		Name:       safeString(data["name"]),
		Scene:      safeString(data["scene"]),
		Pronouns:   safeString(data["pronouns"]),
		Image:      safeString(data["image"]),
		FormFields: make(map[string]string),
	}

	// Check if "team" exists and is not nil, then convert it safely
	if teamVal, ok := data["team"]; ok && teamVal != nil {
		player.Team = int(teamVal.(float64)) // safely cast to int if "team" exists and is float64
	} else {
		player.Team = 0 // Default or placeholder value if "team" is missing or nil
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

		// Parse each player result and add it to the players slice
		for _, playerData := range playerApiResponse.Results {
			parsedPlayer := ParsePlayers(playerData)
			players = append(players, parsedPlayer)
		}

		// Increment the page for the next iteration
		page++
	}

	log.Printf("API data fetched.\nPLAYERS:\n%v", players)
	return players
}
