package models

import "gorm.io/gorm"

type RouteDetail struct {
	NameID  string `json:"name_id"`
	NameEN  string `json:"name_en"`
	Price   string `json:"price"`
}

type Package struct {
	gorm.Model
	// Multi-language fields
	NameID        string        `json:"name_id" gorm:"column:name_id"`
	NameEN        string        `json:"name_en" gorm:"column:name_en"`
	DescriptionID string        `json:"description_id" gorm:"column:description_id"`
	DescriptionEN string        `json:"description_en" gorm:"column:description_en"`
	Capacity      string        `json:"capacity"`
	Duration      string        `json:"duration"`
	IsPopular     bool          `json:"is_popular"`
	ImageURL      string        `json:"image_url"` 
	
	// JSON DB - Multi-language
	RoutesID   []RouteDetail `json:"routes_id" gorm:"serializer:json;column:routes_id"`
	RoutesEN   []RouteDetail `json:"routes_en" gorm:"serializer:json;column:routes_en"`
	FeaturesID []string      `json:"features_id" gorm:"serializer:json;column:features_id"`
	FeaturesEN []string      `json:"features_en" gorm:"serializer:json;column:features_en"`
	ExcludesID []string      `json:"excludes_id" gorm:"serializer:json;column:excludes_id"`
	ExcludesEN []string      `json:"excludes_en" gorm:"serializer:json;column:excludes_en"`
	
	// Legacy fields for backward compatibility (will be deprecated)
	Name        string        `json:"name,omitempty" gorm:"-"`
	Description string        `json:"description,omitempty" gorm:"-"`
	Routes      []RouteDetail `json:"routes,omitempty" gorm:"-"`
	Features    []string      `json:"features,omitempty" gorm:"-"`
	Excludes    []string      `json:"excludes,omitempty" gorm:"-"`
}