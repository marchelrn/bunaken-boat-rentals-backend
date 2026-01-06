package models

import "gorm.io/gorm"

type AddOn struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}