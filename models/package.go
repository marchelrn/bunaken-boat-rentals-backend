package models

import "gorm.io/gorm"

type RouteDetail struct {
	NameID  string `json:"name_id"`
	NameEN  string `json:"name_en"`
	Price   string `json:"price"`
}

type Package struct {
	gorm.Model
	Name   string `json:"name" gorm:"column:name"`
	
	NameID   string `json:"name_id" gorm:"column:name_id"`
	NameEN   string `json:"name_en" gorm:"column:name_en"`
	Capacity string `json:"capacity"`
	Duration      string        `json:"duration"`
	IsPopular     bool          `json:"is_popular"`
	ImageURL      string        `json:"image_url" gorm:"column:image_url;type:varchar(500);default:''"` 
	
	RoutesID   []RouteDetail `json:"routes_id" gorm:"serializer:json;column:routes_id"`
	RoutesEN   []RouteDetail `json:"routes_en" gorm:"serializer:json;column:routes_en"`
	FeaturesID []string      `json:"features_id" gorm:"serializer:json;column:features_id"`
	FeaturesEN []string      `json:"features_en" gorm:"serializer:json;column:features_en"`
	IncludesID []string      `json:"includes_id" gorm:"serializer:json;column:includes_id"`
	IncludesEN []string      `json:"includes_en" gorm:"serializer:json;column:includes_en"`
	ExcludesID []string      `json:"excludes_id" gorm:"serializer:json;column:excludes_id"`
	ExcludesEN []string      `json:"excludes_en" gorm:"serializer:json;column:excludes_en"`
	
	Routes      []RouteDetail `json:"routes,omitempty" gorm:"serializer:json;column:routes"`
	Features    []string      `json:"features,omitempty" gorm:"serializer:json;column:features"`
	Includes    []string      `json:"includes,omitempty" gorm:"serializer:json;column:includes"`
	Excludes    []string      `json:"excludes,omitempty" gorm:"serializer:json;column:excludes"`
}