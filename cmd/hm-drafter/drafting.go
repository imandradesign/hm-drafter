package main

import (
	"fmt"
	"log"
)

type UpdatePlayerTeamRequest struct {
    Team int `json:"team"`
}

func advanceDraftTurn(draftOrder []CaptainDraft) {
	// Move to the next captain
	currentCaptainIndex += draftDirection

	// If we reach the end of the list (ascending), reverse direction and repeat last captain
	if currentCaptainIndex >= len(draftOrder) {
		currentCaptainIndex = len(draftOrder) - 1
		draftDirection = -1
	} else if currentCaptainIndex < 0 { // If we reach the start (descending), reverse direction again
		currentCaptainIndex = 0
		draftDirection = 1
	}
}

func RemoveDraftedPlayers(draftPlayers []Players, selectedPlayer string) (updatedDraftPlayers []Players) {
	for _, player := range draftPlayers {
		if player.Name != selectedPlayer {
			updatedDraftPlayers = append(updatedDraftPlayers, player)
		}
	}

	return updatedDraftPlayers
}

func AddPlayerToDraftTeam(tournamentID string, teams []TeamInfo, captain string, draftedPlayer string) {
	api := fmt.Sprintf("https://kqhivemind.com/api/tournament/player/?tournament_id=%v&format=json", tournamentID)

	client, req := createRequest("PATCH", api, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	
}