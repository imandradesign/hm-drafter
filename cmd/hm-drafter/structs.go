package main

type FormFields struct {
	FieldName        string    `json:"field_name"`
	FieldSlug      string `json:"field_slug"`
	FieldDescription      string `json:"field_description"`
}

type FormApiResponse struct {
	Results []FormFields `json:"results"`
}

type Player struct {
	Name       string            `json:"name"`
	ID         float64           `json:"id"`
	Scene      string            `json:"scene"`
	Pronouns   string            `json:"pronouns"`
	Tournament int               `json:"tournament"`
	Team       int               `json:"team"`
	Image      string            `json:"image"`
	FormFields map[string]string `json:"form_fields,omitempty"`
}

type PlayersApiResponse struct {
	Results []map[string]interface{} `json:"results"`
}

type Captain struct {
	ID float64
	Name  string
	AltName string
	Order int
}

type Team struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Tournament int    `json:"tournament"`
}

type TeamInfo struct {
	ID      int
	Name    string
	Players []Player `json:"players,omitempty"`
}

type TeamApiResponse struct {
	Results []Team `json:"results"`
}

type UpdatePlayerTeamRequest struct {
    Team int `json:"team"`
}

type Tournament struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Date      string `json:"date"`
	SceneName string `json:"scene_name"`
}

type TourneyAPIResponse struct {
	Results []Tournament `json:"results"`
}