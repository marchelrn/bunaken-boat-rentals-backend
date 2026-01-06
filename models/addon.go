package models

import "gorm.io/gorm"

type AddOn struct {
	gorm.Model
	Name        string `json:"name" binding:"required"`
	Price       string `json:"price" binding:"required"`
	Description string `json:"description"`
}
