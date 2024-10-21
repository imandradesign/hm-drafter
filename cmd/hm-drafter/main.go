package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	apiTemplate = "https://kqhivemind.com/api/tournament/%s/?format=json%s%s"
	apiFormFieldsSlug = "player-info-field"
	apiTeamsSlug = "team"
	apiPlayersSlug = "player"
	apiAllTournamentsSlug = "tournament"
	scene = "kqpdx"
)

type Tournament struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Date      string `json:"date"`
	SceneName string `json:"scene_name"`
}

type APIResponse struct {
	Results []Tournament `json:"results"`
}

// type player struct {
// 	Name     string   `json:"name"`
// 	Skill    string   `json:skill`
// 	Roles    []string `json:roles`
// 	Captain  string   `json:captain`
// }

// getPDXTournies takes the API string that lists all tournaments, checks the `results` list for entries where the `scene_name` is `kqpdx` and returns those events with their ID, name, and date in nested lists
func getPDXTournies(slug string) (portlandTournies [][]string) {
	client := &http.Client{}

	listLen := len(portlandTournies)
	page := 1

	for listLen < 10 {
		// Construct API URL with the slug and page number
		api := fmt.Sprintf(apiTemplate, slug, fmt.Sprintf("&page=%d", page), "")
		
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
		var apiResponse APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			log.Fatal(err)
		}

		// Loop through each item in `results` and filter by `scene_name == "kqpdx"`
		for _, tournament := range apiResponse.Results {
			if tournament.SceneName == scene {
				portlandTournies = append(portlandTournies, []string{
					fmt.Sprintf("%d", tournament.ID), tournament.Name, tournament.Date,
				})
			}

			listLen = len(portlandTournies)
			if listLen == 10 {
				return portlandTournies
			}
		}

		page++
		listLen = len(portlandTournies)
	}

	return portlandTournies
}

func main() {
	portlandTournies := getPDXTournies(apiAllTournamentsSlug)
	for _, tourney := range portlandTournies {
		fmt.Printf("ID: %s, Name: %s, Date: %s\n", tourney[0], tourney[1], tourney[2])
	}
}
