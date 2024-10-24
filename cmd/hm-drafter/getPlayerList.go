package hmdrafter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Players struct {
	Name        string    `json:"name"`
	Scene        string    `json:"scene"`
	Pronouns        string    `json:"pronouns"`
	FormFields  map[string]string
}

type PlayersApiResponse struct {
	Results []Players `json:"results"`
}

func ParsePlayers(data map[string]interface{}) Players {
	// Extract static fields like name, scene, pronouns
	player := Players{
		Name:     data["name"].(string),
		Scene:    data["scene"].(string),
		Pronouns: data["pronouns"].(string),
	}

	// Prepare dynamic form fields
	player.FormFields = make(map[string]string)
	for _, field := range formFields {
		fieldName := field[0] // e.g., "question216"
		if value, ok := data[fieldName]; ok {
			player.FormFields[field[1]] = value.(string) // field[1] is the slug like "discord"
		}
	}
	return player
}

func GetPlayersData(tournamentId string) (players []Players) {
	log.Println("Fetching player data...")
	client := &http.Client{
    	Timeout: 10 * time.Second,
	}

	page := 1
	for {
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

		// Close the response body after reading
		func() {
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
				return
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
				page = 0 // Break the loop
			}
		}()

		// Break the loop if there are no more pages
		if page == 0 {
			break
		}
	}

	return players
}