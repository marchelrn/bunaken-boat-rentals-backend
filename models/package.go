package models

import "gorm.io/gorm"

type RouteDetail struct {
	Name  string `json:"name"`
	Price string `json:"price"`
}

type Package struct {
	gorm.Model
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Capacity    string        `json:"capacity"`
	Duration    string        `json:"duration"`
	IsPopular   bool          `json:"is_popular"`
	ImageURL    string        `json:"image_url"` 
	
	// JSON DB
	Routes      []RouteDetail `json:"routes" gorm:"serializer:json"`
	Features    []string      `json:"features" gorm:"serializer:json"`
	Excludes    []string      `json:"excludes" gorm:"serializer:json"`
}