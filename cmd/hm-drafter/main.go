package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
	selectedTournament   []string
	formFields           [][]string
	tournamentID         string
	players              []Players
	playerCount          int
	captains             []CaptainDraft
	captainNamesFromForm []string
	captainCount         int
	draftPlayers         []Players
	draftOrder           []CaptainDraft
	remaininPlayerCount int
	currentCaptainIndex  int
	draftDirection       int
	teams                []TeamInfo
	unassignedCaptains   []CaptainDraft
)

// Helper function to create an HTTP request with the API key
func createRequest(method, url string, body io.Reader) (client *http.Client, req *http.Request) {
	client = &http.Client{
		Timeout: 10 * time.Second,
	}
	
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	// Retrieve the API key from the environment
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API key not set in environment")
	}

	// Attach the Authorization header
	req.Header.Add("Authorization", "Token "+apiKey)

	return client, req
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000" // Default
	}

	router := gin.New()
	router.Use(gin.Logger())

	// Load HTML templates
	router.LoadHTMLFiles("templates/index.html", "templates/drafting.html", "templates/teams.html", "templates/done.html")

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

	// Handle the form submission for captain selection
	router.POST("/confirm-captains", func(c *gin.Context) {
		// Get the selected captains from the form
		captainNamesFromForm = c.PostFormArray("selectedPlayers")
		captains := GetCaptainInfoFromNames(captainNamesFromForm, players)
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

		c.Redirect(http.StatusFound, "/teams")
	})

	// Teams page route
	router.GET("/teams", func(c *gin.Context) {
		teams = GetTeams(tournamentID, players)
		unassignedCaptains = draftOrder

		log.Printf("Unassigned Captain Data:\n%v", unassignedCaptains)

		c.HTML(http.StatusOK, "teams.html", gin.H{
			"selectedTournament":  selectedTournament,
			"remaininPlayerCount": remaininPlayerCount,
			"playerCount":         playerCount,
			"captainCount":        captainCount,
			"unassignedCaptains":  unassignedCaptains,
			"draftOrder":          draftOrder,
			"teams":               teams,
		})
	})

	// Handle the form submission for adding new teams
	router.POST("/add-team", func(c *gin.Context) {
		teamName := c.PostForm("teamAddition")

		AddNewTeam(teamName, tournamentID)

		c.Redirect(http.StatusFound, "/teams")
	})

	// Handle the form submission for deleting teams
	router.POST("/remove-team", func(c *gin.Context) {
		teamID := c.PostForm("teamDeletion")
		log.Printf("Team ID for removal: %v", teamID)

		teamName := GetTeamNameByID(teams, teamID)
		log.Printf("Team name for removal: %v", teamName)

		DeleteTeam(teamID, teamName, tournamentID)

		c.Redirect(http.StatusFound, "/teams")
	})

	router.POST("/assign-captain",func(c *gin.Context){
		cap := c.PostForm("captainID")
		team := c.PostForm("teamID")

		log.Printf("Captain ID: %v\nTeam ID: %v", cap, team)

		AssignPlayerToTeam(cap, team, tournamentID)
		unassignedCaptains = updateUnassignedCaptainList(cap, unassignedCaptains)

		c.Redirect(http.StatusFound, "/teams")
	})

	// Redirect to Drafting page after confirming teams
	router.POST("/confirm-teams",func(c *gin.Context) {
		// If no teams exist, return an error message
		if len(teams) == 0 {
			c.String(http.StatusBadRequest, "No teams created. Please create at least one team.")
			return
		}

		c.Redirect(http.StatusFound, "/drafting")
	})

	// Drafting page route
	router.GET("/drafting", func(c *gin.Context) {
		teams = GetTeams(tournamentID, players)
		currCaptain := draftOrder[currentCaptainIndex].Name

		c.HTML(http.StatusOK, "drafting.html", gin.H{
			"selectedTournament": selectedTournament,
			"captainCount": captainCount,
			"remaininPlayerCount": remaininPlayerCount,
			"draftOrder": draftOrder,
			"draftPlayers": draftPlayers,
			"currentCaptain": currCaptain,
			"teams": teams,
		})
	})

	// Handle the form submission for player selection & advance the draft turn
	router.POST("/pick-player", func(c *gin.Context) {
		currCaptain := draftOrder[currentCaptainIndex].Name
		selectedPlayer := c.PostForm("selectedPlayer")

		AddPlayerToDraftTeam(tournamentID, teams, currCaptain,selectedPlayer)

		// Get updated teams list
		teams = GetTeams(tournamentID, players)

		// Remove the selected player from the list
		draftPlayers = RemoveDraftedPlayers(draftPlayers, selectedPlayer)

		if len(draftPlayers) == 0 {
			c.Redirect(http.StatusFound, "/done")
			return
		}

		// Advance the draft turn
		advanceDraftTurn(draftOrder)

		// Re-render the drafting page with the updated current captain
		c.HTML(http.StatusOK, "drafting.html", gin.H{
			"selectedTournament": selectedTournament,
			"captainCount": captainCount,
			"remaininPlayerCount": remaininPlayerCount,
			"draftOrder": draftOrder,
			"draftPlayers": draftPlayers,
			"currentCaptain": currCaptain,
			"teams": teams,
		})
	})

	// Final page route
	router.GET("/done", func(c *gin.Context) {
		c.HTML(http.StatusOK, "done.html", gin.H{
			"selectedTournament": selectedTournament,
			"teams": teams,
		})
	})

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
