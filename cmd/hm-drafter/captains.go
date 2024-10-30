package main

import (
	"math/rand"
	"sort"
	"time"
)

type CaptainDraft struct {
	Name  string
	AltName string
	Order int
}

// RemoveCaptainsFromPlayers returns a new player list without captains
func RemoveCaptainsFromPlayers(players []Players, captains []string) (draftPlayers []Players) {
	captainSet := make(map[string]bool)
	
	// Create a set of captain names for quick lookups
	for _, captain := range captains {
		captainSet[captain] = true
	}

	// Loop through players and only append those who are not captains
	for _, player := range players {
		if !captainSet[player.Name] {
			draftPlayers = append(draftPlayers, player)
		}
	}

	return draftPlayers
}

func CaptainDraftOrder(captains []string) (draftOrder []CaptainDraft) {
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
		draftOrder = append(draftOrder, CaptainDraft{
			Name:  captain,
			Order: draftPositions[i],
		})
	}

	// Sort the draftOrder slice by the Order field in ascending order
	sort.Slice(draftOrder, func(i, j int) bool {
		return draftOrder[i].Order < draftOrder[j].Order
	})

	return draftOrder
}