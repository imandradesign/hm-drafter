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
	apiTemplate = "https://kqhivemind.com/api/tournament/%s/?format=json%s%s"
	apiFormFieldsSlug = "player-info-field"
	apiTeamsSlug = "team"
	apiPlayersSlug = "player"
	scene = "kqpdx"
)

var (
	selectedTournament []string
	formFields [][]string
	tournamentID string
	players []Players
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000" // Default
	}

	router := gin.New()
	router.Use(gin.Logger())

	// Load HTML templates
	router.LoadHTMLFiles("templates/index.html")

	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		// Fetch tournament data
		tournaments := GetPDXTournies()

		c.HTML(http.StatusOK, "index.html", gin.H{
			"tournaments": tournaments,
			"selectedTournament": selectedTournament,
			"players":            players,
		})
	})

	// Handle the form submission and store the selected tournament
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

		// Redirect back to the home page to show the selected tournament
		c.Redirect(http.StatusFound, "/")
	})

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}