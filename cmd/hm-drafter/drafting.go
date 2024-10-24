package main

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