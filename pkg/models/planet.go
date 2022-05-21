package models

import "gorm.io/gorm"

type Planet struct {
	gorm.Model
	Name         string `json:"name"`
	Weather      string `json:"weather"`
	Ground       string `json:"ground"`
	MoviesNumber int    `json:"movies_number"`
}
