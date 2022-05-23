package models

import "gorm.io/gorm"

type Response struct {
	Result []Planets `json:"results"`
}

type Planets struct {
	Name    string   `json:"name"`
	Climate string   `json:"climate"`
	Terrain string   `json:"terrain"`
	Films   []string `json:"films"`
}

type Planet struct {
	gorm.Model
	Name         string `json:"name"`
	Climate      string `json:"climate"`
	Terrain      string `json:"terrain"`
	MoviesNumber int    `json:"movies_number"`
}
