package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

const (
	apiTemplate      = "https://kqhivemind.com/api/tournament/%s/?format=json%s%s"
	apiFormFieldsSlug = "player-info-field"
	apiTeamsSlug      = "team"
	apiPlayersSlug    = "player"
	scene            = "kqpdx"
)

var (
	selectedTournament  []string
	formFields          [][]string
	tournamentID        string
	players             []Players
	playerCount         int
	captains            []string
	captainCount        int
	draftPlayers        []Players
	draftOrder          []CaptainDraft
	remaininPlayerCount int
	currentCaptainIndex int
	draftDirection      int
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000" // Default
	}

	router := gin.New()
	router.Use(gin.Logger())

	// Load HTML templates
	router.LoadHTMLFiles("templates/index.html", "templates/drafting.html")

	router.Static("/static", "./static")

	// Home page route
	router.GET("/", func(c *gin.Context) {
		// Fetch tournament data
		tournaments := GetPDXTournies()

		c.HTML(http.StatusOK, "index.html", gin.H{
			"tournaments":        tournaments,
			"selectedTournament": selectedTournament,
			"playerCount":        playerCount,
			"players":            players,
		})
	})

	// Handle the form submission for tournament selection
	router.POST("/confirm", func(c *gin.Context) {
		// Get the selected tournament ID from the form
		tournamentIndexStr := c.PostForm("tournament")
		tournamentIndex, err := strconv.Atoi(tournamentIndexStr)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid tournament selection")
			return
		}

		// Fetch tournaments and retrieve the selected one
		tournaments := GetPDXTournies()
		if tournamentIndex >= 0 && tournamentIndex < len(tournaments) {
			selectedTournament = tournaments[tournamentIndex]
		} else {
			c.String(http.StatusBadRequest, "Tournament not found")
			return
		}

		tournamentID = selectedTournament[0]

		// Fetch form fields for selected tournament
		formFields = GetFormFields(tournamentID)

		// Fetch player data
		players = GetPlayersData(tournamentID)
		playerCount = len(players)
		log.Printf("# of players: %v", playerCount)

		// Stay on homepage when confirming tournament selection
		c.Redirect(http.StatusFound, "/")
	})

	// Handle the captain confirmation POST request
	router.POST("/confirm-captains", func(c *gin.Context) {
		// Get the selected captains from the form
		captains = c.PostFormArray("selectedPlayers")
		captainCount = len(captains)

		// If no captains were selected, return an error message
		if len(captains) == 0 {
			c.String(http.StatusBadRequest, "No captains selected. Please select at least one captain.")
			return
		}

		// Save initial draft info
		draftPlayers = RemoveCaptainsFromPlayers(players, captains)
		remaininPlayerCount = len(draftPlayers)

		// Set initial values for the draft state
		draftOrder = CaptainDraftOrder(captains)
		currentCaptainIndex = 0 // Start with the first captain
		draftDirection = 1      // Start with ascending order

		c.Redirect(http.StatusFound, "/drafting")
	})

	// Drafting page route (accessible at /drafting)
	router.GET("/drafting", func(c *gin.Context) {
		c.HTML(http.StatusOK, "drafting.html", gin.H{
			"selectedTournament": selectedTournament,
			"captainCount": captainCount,
			"remaininPlayerCount": remaininPlayerCount,
			"draftOrder": draftOrder,
			"draftPlayers": draftPlayers,
			"currentCaptain": draftOrder[currentCaptainIndex].Name,
		})
	})

	// Handle the player selection and advance the draft turn
	router.POST("/pick-player", func(c *gin.Context) {
		selectedPlayer := c.PostForm("selectedPlayer")

		// Remove the selected player from the list
		draftPlayers = RemoveDraftedPlayers(draftPlayers, selectedPlayer)

		// Advance the draft turn
		advanceDraftTurn(draftOrder)

		// Re-render the drafting page with the updated current captain
		c.HTML(http.StatusOK, "drafting.html", gin.H{
			"selectedTournament": selectedTournament,
			"captainCount": captainCount,
			"remaininPlayerCount": remaininPlayerCount,
			"draftOrder": draftOrder,
			"draftPlayers": draftPlayers,
			"currentCaptain": draftOrder[currentCaptainIndex].Name,
		})
	})


	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
