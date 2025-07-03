package models

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type PostCardResponse struct {
	ALT      string
	IMG      string
	CATEGORY string
	LINK     string
	TITLE    string
	EXCERPT  string
	VIEWS    string
	AUTHOR   string
	DATE     string
	Slug     string
}
