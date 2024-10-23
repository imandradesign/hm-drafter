package main

import (
	"html/template"
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch Portland tournaments using the GetPDXTournies function
	portlandTournies := GetPDXTournies(apiAllTournamentsSlug)

	// Parse and execute the HTML template
	tmpl := template.Must(template.ParseFiles("../../index.html"))
	tmpl.Execute(w, portlandTournies)
}

func main() {
	// Set up the home route and handler
	http.HandleFunc("/", homeHandler)

	// Start the web server
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
