package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Players struct {
	Name       string    `json:"name"`
	ID         float64       `json:"user"`
	Scene      string    `json:"scene"`
	Pronouns   string    `json:"pronouns"`
	FormFields map[string]string
	TeamInfo *TeamInfo `json:"team_info,omitempty"`
}

type PlayersApiResponse struct {
	Results []Players `json:"results"`
}

func ParsePlayers(data map[string]interface{}, teamMap map[int]TeamInfo) Players {
	// Extract static fields like name, scene, pronouns, and ID
	player := Players{
		Name:     safeString(data["name"]),
		Scene:    safeString(data["scene"]),
		Pronouns: safeString(data["pronouns"]),
		ID:       data["user"].(float64),
	}

	// Prepare dynamic form fields
	player.FormFields = make(map[string]string)
	for _, field := range formFields {
		fieldName := field[0] // e.g., "question216"
		if value, ok := data[fieldName]; ok {
			// Handle different types the value could have
			switch v := value.(type) {
			case string:
				player.FormFields[field[1]] = v // field[1] is the slug like "discord"
			case []interface{}:
				// Convert array to a string representation (comma-separated)
				var strValues []string
				for _, item := range v {
					strValues = append(strValues, fmt.Sprintf("%v", item))
				}
				player.FormFields[field[1]] = strings.Join(strValues, ", ") // Join without brackets
			default:
				player.FormFields[field[1]] = fmt.Sprintf("%v", value)
			}
		}
	}

	// Assign team info if team ID exists in teamMap
	if teamID, ok := data["team"].(float64); ok {
		intTeamID := int(teamID)
		if team, found := teamMap[intTeamID]; found {
			player.TeamInfo = &team
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

// GetPlayersData uses the tournament ID to retrieve player and form field data from that specific event. Returns player count as well.
func GetPlayersData(tournamentId string) (players []Players) {
	log.Println("Fetching player data...")

	// Fetch teams for the tournament
    teams := GetTeams(tournamentId)
    teamMap := make(map[int]TeamInfo)
    for _, team := range teams {
        teamMap[team.ID] = team
    }

	page := 1
	for {
		// Build the API URL for the current page
		api := fmt.Sprintf(apiTemplate, "player", fmt.Sprintf("&tournament_id=%v", tournamentId), fmt.Sprintf("&page=%d", page))
		client, req := createRequest("GET", api, nil)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Parse the response into a generic map
		var rawResponse map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&rawResponse)
		if err != nil {
			log.Fatal(err)
		}

		// Check if the results array is empty
		results, ok := rawResponse["results"].([]interface{})
		if !ok || len(results) == 0 {
			log.Println("No players found.")
			break
		}

		// Parse each player
		for _, playerData := range results {
			player := ParsePlayers(playerData.(map[string]interface{}), teamMap)
			players = append(players, player)
		}

		// Check if there is a next page
		if nextURL, ok := rawResponse["next"].(string); ok && nextURL != "" {
			page++
		} else {
			log.Println("No more pages.")
			break
		}
	}

	log.Printf("PLAYERS: \n%v", players)
	log.Println("API data fetched.")
	return players
}
