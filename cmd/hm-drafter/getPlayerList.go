package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Players struct {
	Name       string    `json:"name"`
	Scene      string    `json:"scene"`
	Pronouns   string    `json:"pronouns"`
	FormFields map[string]string
}

type PlayersApiResponse struct {
	Results []Players `json:"results"`
}

func ParsePlayers(data map[string]interface{}) Players {
	// Extract static fields like name, scene, pronouns, checking if they are nil first
	player := Players{
		Name:     safeString(data["name"]),
		Scene:    safeString(data["scene"]),
		Pronouns: safeString(data["pronouns"]),
	}

	// Prepare dynamic form fields
	player.FormFields = make(map[string]string)
	for _, field := range formFields {
		fieldName := field[0] // e.g., "question216"
		if value, ok := data[fieldName]; ok {
			// Handle the different types the value could have
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	page := 1
	for {
		log.Printf("Fetching page: %d", page)

		// Build the API URL for the current page
		api := fmt.Sprintf(apiTemplate, "player", fmt.Sprintf("&tournament_id=%v", tournamentId), fmt.Sprintf("&page=%d", page))

		req, err := http.NewRequest("GET", api, nil)
		if err != nil {
			log.Fatal(err)
		}

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
			player := ParsePlayers(playerData.(map[string]interface{}))
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
