package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

const (
	apiTemplate = "https://kqhivemind.com/api/tournament/%s/?format=json%s%s"
	apiFormFieldsSlug = "player-info-field"
	apiTeamsSlug = "team"
	apiPlayersSlug = "player"
	apiAllTournamentsSlug = "tournament"
	scene = "kqpdx"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000" // Default
	}

	router := gin.New()
	router.Use(gin.Logger())

	// Load HTML templates
	router.LoadHTMLFiles("index.html")

	router.GET("/", func(c *gin.Context) {
		// Fetch tournament data using GetPDXTournies
		tournaments := GetPDXTournies(apiAllTournamentsSlug)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"tournaments": tournaments,
		})
	})

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}