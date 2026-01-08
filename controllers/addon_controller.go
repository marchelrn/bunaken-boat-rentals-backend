package controllers

import (
	"bunaken-boat-backend/config"
	"bunaken-boat-backend/models"
	"bunaken-boat-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllAddOns(c *gin.Context) {
	lang := c.Query("lang")
	if lang == "" {
		lang = "id"
	}

	addOns := []models.AddOn{}
	if err := config.DB.Find(&addOns).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal mengambil data add-ons")
		return
	}

	// Transform add-ons based on language
	transformedAddOns := make([]map[string]interface{}, len(addOns))
	for i, addOn := range addOns {
		addOnMap := map[string]interface{}{
			"ID":          addOn.ID,
			"CreatedAt":   addOn.CreatedAt,
			"UpdatedAt":   addOn.UpdatedAt,
			"DeletedAt":   addOn.DeletedAt,
			"price":      addOn.Price,
			// Include all multi-language fields for admin dashboard
			"name_id":        addOn.NameID,
			"name_en":        addOn.NameEN,
			"description_id": addOn.DescriptionID,
			"description_en": addOn.DescriptionEN,
		}

		// For public API, return language-specific fields
		if lang == "en" {
			if addOn.NameEN != "" {
				addOnMap["name"] = addOn.NameEN
			} else if addOn.NameID != "" {
				addOnMap["name"] = addOn.NameID
			} else {
				addOnMap["name"] = addOn.Name
			}
			if addOn.DescriptionEN != "" {
				addOnMap["description"] = addOn.DescriptionEN
			} else if addOn.DescriptionID != "" {
				addOnMap["description"] = addOn.DescriptionID
			} else {
				addOnMap["description"] = addOn.Description
			}
		} else {
			// Default to ID
			if addOn.NameID != "" {
				addOnMap["name"] = addOn.NameID
			} else if addOn.NameEN != "" {
				addOnMap["name"] = addOn.NameEN
			} else {
				addOnMap["name"] = addOn.Name
			}
			if addOn.DescriptionID != "" {
				addOnMap["description"] = addOn.DescriptionID
			} else if addOn.DescriptionEN != "" {
				addOnMap["description"] = addOn.DescriptionEN
			} else {
				addOnMap["description"] = addOn.Description
			}
		}

		transformedAddOns[i] = addOnMap
	}

	utils.APIResponse(c, http.StatusOK, "Berhasil mengambil data Add-Ons", transformedAddOns)
}

func GetAddOnByID(c *gin.Context) {
	id := c.Param("id")
	var addOn models.AddOn

	if err := config.DB.First(&addOn, id).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "Add-On tidak ditemukan")
		return
	}

	utils.APIResponse(c, http.StatusOK, "Berhasil mengambil detail Add-On", addOn)
}

func CreateAddOn(c *gin.Context) {
	var input models.AddOn

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Format JSON salah: "+err.Error())
		return
	}

	// Set legacy fields from multi-language fields if needed
	if input.Name == "" {
		if input.NameID != "" {
			input.Name = input.NameID
		} else if input.NameEN != "" {
			input.Name = input.NameEN
		}
	}
	if input.Description == "" {
		if input.DescriptionID != "" {
			input.Description = input.DescriptionID
		} else if input.DescriptionEN != "" {
			input.Description = input.DescriptionEN
		}
	}

	// Ensure multi-language fields are set from legacy fields if needed
	if input.NameID == "" && input.Name != "" {
		input.NameID = input.Name
	}
	if input.NameEN == "" && input.Name != "" {
		input.NameEN = input.Name
	}
	if input.DescriptionID == "" && input.Description != "" {
		input.DescriptionID = input.Description
	}
	if input.DescriptionEN == "" && input.Description != "" {
		input.DescriptionEN = input.Description
	}

	result := config.DB.Create(&input)
	if result.Error != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal menyimpan ke database: "+result.Error.Error())
		return
	}
	utils.APIResponse(c, http.StatusOK, "Berhasil membuat add-on baru", input)
}

func UpdateAddOn(c *gin.Context) {
	id := c.Param("id")
	var addOn models.AddOn

	if err := config.DB.First(&addOn, id).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "Add-On tidak ditemukan")
		return
	}

	var input models.AddOn
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Format JSON salah: "+err.Error())
		return
	}

	// Update price (always required)
	addOn.Price = input.Price

	// Update multi-language fields if provided
	if input.NameID != "" {
		addOn.NameID = input.NameID
	}
	if input.NameEN != "" {
		addOn.NameEN = input.NameEN
	}
	if input.DescriptionID != "" {
		addOn.DescriptionID = input.DescriptionID
	}
	if input.DescriptionEN != "" {
		addOn.DescriptionEN = input.DescriptionEN
	}

	// Update legacy fields if provided, or set from multi-language fields
	if input.Name != "" {
		addOn.Name = input.Name
	} else {
		// Set legacy field from multi-language fields if not provided
		if addOn.NameID != "" {
			addOn.Name = addOn.NameID
		} else if addOn.NameEN != "" {
			addOn.Name = addOn.NameEN
		}
	}

	if input.Description != "" {
		addOn.Description = input.Description
	} else {
		// Set legacy field from multi-language fields if not provided
		if addOn.DescriptionID != "" {
			addOn.Description = addOn.DescriptionID
		} else if addOn.DescriptionEN != "" {
			addOn.Description = addOn.DescriptionEN
		}
	}

	if err := config.DB.Save(&addOn).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal mengupdate database: "+err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Berhasil update add-on", addOn)
}

func DeleteAddOn(c *gin.Context) {
	id := c.Param("id")
	var addOn models.AddOn

	if err := config.DB.First(&addOn, id).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "Add-On tidak ditemukan")
		return
	}

	if err := config.DB.Unscoped().Delete(&addOn).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal menghapus data")
		return
	}

	utils.APIResponse(c, http.StatusOK, "Berhasil menghapus add-on", nil)
}

