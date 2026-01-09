package models

import "gorm.io/gorm"

type AddOn struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name;type:varchar(255);default:''"`
	Price       string `json:"price" binding:"required" gorm:"column:price;type:varchar(50)"`
	Description string `json:"description" gorm:"column:description;type:text;default:''"`
	NameID        string `json:"name_id" gorm:"column:name_id;type:varchar(255);default:''"`
	NameEN        string `json:"name_en" gorm:"column:name_en;type:varchar(255);default:''"`
	DescriptionID string `json:"description_id" gorm:"column:description_id;type:text;default:''"`
	DescriptionEN string `json:"description_en" gorm:"column:description_en;type:text;default:''"`
}
