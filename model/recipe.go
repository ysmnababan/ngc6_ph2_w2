package model

type Recipe struct {
	ID          int    `json:"id"`
	Name        string `json:"nama"`
	Description string `json:"description"`
	Rating      int    `json:"rating"`
}
