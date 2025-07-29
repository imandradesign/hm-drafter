package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// RemoveCaptainsFromPlayers returns a new player list without captains
func RemoveCaptainsFromPlayers(players []Player, captains []Captain) (draftPlayers []Player) {
	captainSet := make(map[string]bool)
	
	// Create a set of captain names for quick lookups
	for _, captain := range captains {
		captainSet[captain.Name] = true
	}

	// Loop through players and only append those who are not captains
	for _, player := range players {
		if !captainSet[player.Name] {
			draftPlayers = append(draftPlayers, player)
		}
	}

	return draftPlayers
}


func ExtractCaptains(captainNamesList []string, players []Player) (captains []Captain) {
	// Create a map for quick lookup of captain names
	captainNames := make(map[string]bool)
	for _, name := range captainNamesList {
		captainNames[name] = true
	}

	// Iterate through players to find matches for captain names
	for _, player := range players {
		if captainNames[player.Name] {
			captains = append(captains, Captain{
				ID:      player.ID,
				Name:    player.Name,
				AltName: player.FormFields["altname"],
			})
		}
	}

	return captains
}


func GenerateDraftOrder(captains []Captain) (draftOrder []Captain) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a slice of numbers from 1 to the number of captains
	draftPositions := make([]int, len(captains))
	for i := range draftPositions {
		draftPositions[i] = i + 1
	}

	// Shuffle the draft positions using the local random generator
	r.Shuffle(len(draftPositions), func(i, j int) {
		draftPositions[i], draftPositions[j] = draftPositions[j], draftPositions[i]
	})

	// Pair each captain with a draft position and add to draftOrder
	for i, captain := range captains {
		draftOrder = append(draftOrder, Captain{
			ID:      captain.ID,
			Name:    captain.Name,
			AltName: captain.AltName,
			Order:   draftPositions[i],
		})
	}

	// Sort the draftOrder slice by the Order field in ascending order
	sort.Slice(draftOrder, func(i, j int) bool {
		return draftOrder[i].Order < draftOrder[j].Order
	})

	return draftOrder
}


// Adds or removes players from the unassignedCaptains var. If addCap bool is True, player list is searched and the one with the matching captain ID is appended. IF addCap is False, captain is removed from list.
func UpdateUnassignedCaptains(captainID string, captains []Captain, addCap bool) (unassignedCaptains []Captain) {
	// Convert captainID from string to float64
	id, err := strconv.ParseFloat(captainID, 64)
	if err != nil {
		log.Printf("Invalid captain ID: %v", captainID)
		return captains // return the original list if there's an error
	}

	// Search the full players list for the one with the same ID as the captain we're adding to the unassignedCaptains list
	if addCap == true {
		for _, player := range players {
			if id == player.ID {
				unassignedCaptains = append(unassignedCaptains, Captain{
					ID:      player.ID,
					Name:    player.Name,
					AltName: player.FormFields["altname"],
				})
				break // Only add one matching player
			}
		}
	}

	// Add all captains except the one recently assigned to a team back to the unassignedCaptains list.
	if addCap == false {
		// Filter out the captain with the matching ID
		for _, captain := range captains {
			if captain.ID != id {
				unassignedCaptains = append(unassignedCaptains, captain)
			}
		}
	}
	return unassignedCaptains
}