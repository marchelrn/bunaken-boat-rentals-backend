package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"` // Ini nanti berisi hash, bukan plain text
	Role     string `json:"role"`
}