package controllers

import (
	"bunaken-boat-backend/config"
	"bunaken-boat-backend/models"
	"bunaken-boat-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllPackages(c *gin.Context) {
	packages := []models.Package{} // Inisialisasi sebagai empty slice
	// Preload features or other relations if necessary, but currently they are JSON fields
	if err := config.DB.Find(&packages).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal Mengambil database")
		return
	}
	
	// Get language parameter (default: "id")
	lang := c.DefaultQuery("lang", "id")
	if lang != "id" && lang != "en" {
		lang = "id"
	}
	
		// Transform packages based on language
	transformedPackages := make([]map[string]interface{}, len(packages))
	for i, pkg := range packages {
		pkgMap := map[string]interface{}{
			"ID":          pkg.ID,
			"CreatedAt":   pkg.CreatedAt,
			"UpdatedAt":   pkg.UpdatedAt,
			"DeletedAt":   pkg.DeletedAt,
			"capacity":    pkg.Capacity,
			"duration":    pkg.Duration,
			"is_popular":  pkg.IsPopular,
			"image_url":   pkg.ImageURL,
			// Include all multi-language fields for admin dashboard
			"name_id":        pkg.NameID,
			"name_en":        pkg.NameEN,
			"description_id": pkg.DescriptionID,
			"description_en": pkg.DescriptionEN,
			"routes_id":      pkg.RoutesID,
			"routes_en":      pkg.RoutesEN,
			"features_id":    pkg.FeaturesID,
			"features_en":    pkg.FeaturesEN,
			"excludes_id":    pkg.ExcludesID,
			"excludes_en":    pkg.ExcludesEN,
			// Legacy fields
			"name":        pkg.Name,
			"description": pkg.Description,
			"routes":      pkg.Routes,
			"features":    pkg.Features,
			"excludes":    pkg.Excludes,
		}
		
		// Set language-specific fields with fallback to legacy fields
		if lang == "en" {
			// Use EN fields if available, otherwise fallback to ID fields
			if pkg.NameEN != "" {
				pkgMap["name"] = pkg.NameEN
			} else {
				pkgMap["name"] = pkg.NameID
			}
			if pkg.DescriptionEN != "" {
				pkgMap["description"] = pkg.DescriptionEN
			} else {
				pkgMap["description"] = pkg.DescriptionID
			}
			if len(pkg.RoutesEN) > 0 {
				pkgMap["routes"] = pkg.RoutesEN
			} else {
				pkgMap["routes"] = pkg.RoutesID
			}
			if len(pkg.FeaturesEN) > 0 {
				pkgMap["features"] = pkg.FeaturesEN
			} else {
				pkgMap["features"] = pkg.FeaturesID
			}
			if len(pkg.ExcludesEN) > 0 {
				pkgMap["excludes"] = pkg.ExcludesEN
			} else {
				pkgMap["excludes"] = pkg.ExcludesID
			}
		} else {
			// Use ID fields if available, otherwise fallback to legacy fields
			if pkg.NameID != "" {
				pkgMap["name"] = pkg.NameID
			} else if pkg.Name != "" {
				// Legacy fallback - use old Name field
				pkgMap["name"] = pkg.Name
			} else if pkg.Description != "" {
				// If name is empty but description exists, use description as name
				// (for backward compatibility with old data structure)
				pkgMap["name"] = pkg.Description
			} else {
				pkgMap["name"] = ""
			}
			if pkg.DescriptionID != "" {
				pkgMap["description"] = pkg.DescriptionID
			} else if pkg.Description != "" {
				// Legacy fallback - use old Description field
				pkgMap["description"] = pkg.Description
			} else {
				pkgMap["description"] = ""
			}
			if len(pkg.FeaturesID) > 0 {
				pkgMap["features"] = pkg.FeaturesID
			} else if len(pkg.Features) > 0 {
				// Legacy fallback - use old Features field
				pkgMap["features"] = pkg.Features
			} else {
				pkgMap["features"] = []string{}
			}
			if len(pkg.ExcludesID) > 0 {
				pkgMap["excludes"] = pkg.ExcludesID
			} else if len(pkg.Excludes) > 0 {
				// Legacy fallback - use old Excludes field
				pkgMap["excludes"] = pkg.Excludes
			} else {
				pkgMap["excludes"] = []string{}
			}
		}
		
		// Transform routes to match frontend format
		var routes []models.RouteDetail
		if lang == "en" {
			if len(pkg.RoutesEN) > 0 {
				routes = pkg.RoutesEN
			} else if len(pkg.RoutesID) > 0 {
				routes = pkg.RoutesID
			} else if len(pkg.Routes) > 0 {
				// Legacy fallback
				routes = pkg.Routes
			} else {
				routes = []models.RouteDetail{}
			}
		} else {
			if len(pkg.RoutesID) > 0 {
				routes = pkg.RoutesID
			} else if len(pkg.Routes) > 0 {
				// Legacy fallback - use old Routes field
				routes = pkg.Routes
			} else {
				routes = []models.RouteDetail{}
			}
		}
		
		// Transform routes to match frontend format
		// For main page display, use simplified format
		transformedRoutes := make([]map[string]string, len(routes))
		for j, route := range routes {
			routeMap := map[string]string{
				"price": route.Price,
			}
			if lang == "en" {
				if route.NameEN != "" {
					routeMap["name"] = route.NameEN
				} else if route.NameID != "" {
					routeMap["name"] = route.NameID
				} else {
					routeMap["name"] = ""
				}
			} else {
				if route.NameID != "" {
					routeMap["name"] = route.NameID
				} else {
					routeMap["name"] = ""
				}
			}
			transformedRoutes[j] = routeMap
		}
		pkgMap["routes"] = transformedRoutes
		
		// Also include full route details for admin dashboard
		// Transform RoutesID and RoutesEN to include name_id and name_en
		transformedRoutesID := make([]map[string]interface{}, len(pkg.RoutesID))
		for j, route := range pkg.RoutesID {
			transformedRoutesID[j] = map[string]interface{}{
				"name":    route.NameID,
				"name_id": route.NameID,
				"name_en": route.NameEN,
				"price":   route.Price,
			}
		}
		transformedRoutesEN := make([]map[string]interface{}, len(pkg.RoutesEN))
		for j, route := range pkg.RoutesEN {
			transformedRoutesEN[j] = map[string]interface{}{
				"name":    route.NameEN,
				"name_id": route.NameID,
				"name_en": route.NameEN,
				"price":   route.Price,
			}
		}
		pkgMap["routes_id"] = transformedRoutesID
		pkgMap["routes_en"] = transformedRoutesEN
		
		transformedPackages[i] = pkgMap
	}
	
	utils.APIResponse(c, http.StatusOK, "Berhasil mengambil data Packages", transformedPackages)
}

func GetPackageByID(c *gin.Context) {
	id := c.Param("id")
	var pkg models.Package

	if err := config.DB.First(&pkg, id).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "Package tidak ditemukan")
		return
	}

	// Get language parameter (default: "id")
	lang := c.DefaultQuery("lang", "id")
	if lang != "id" && lang != "en" {
		lang = "id"
	}
	
	// Transform package based on language
	pkgMap := map[string]interface{}{
		"ID":          pkg.ID,
		"CreatedAt":   pkg.CreatedAt,
		"UpdatedAt":   pkg.UpdatedAt,
		"DeletedAt":   pkg.DeletedAt,
		"capacity":    pkg.Capacity,
		"duration":    pkg.Duration,
		"is_popular":  pkg.IsPopular,
		"image_url":   pkg.ImageURL,
	}
	
	// Set language-specific fields with fallback
	var routes []models.RouteDetail
	if lang == "en" {
		if pkg.NameEN != "" {
			pkgMap["name"] = pkg.NameEN
		} else {
			pkgMap["name"] = pkg.NameID
		}
		if pkg.DescriptionEN != "" {
			pkgMap["description"] = pkg.DescriptionEN
		} else {
			pkgMap["description"] = pkg.DescriptionID
		}
		if len(pkg.RoutesEN) > 0 {
			routes = pkg.RoutesEN
		} else {
			routes = pkg.RoutesID
		}
		if len(pkg.FeaturesEN) > 0 {
			pkgMap["features"] = pkg.FeaturesEN
		} else {
			pkgMap["features"] = pkg.FeaturesID
		}
		if len(pkg.ExcludesEN) > 0 {
			pkgMap["excludes"] = pkg.ExcludesEN
		} else {
			pkgMap["excludes"] = pkg.ExcludesID
		}
	} else {
		if pkg.NameID != "" {
			pkgMap["name"] = pkg.NameID
		} else if pkg.Name != "" {
			// Legacy fallback
			pkgMap["name"] = pkg.Name
		} else if pkg.Description != "" {
			// If name is empty but description exists, use description as name
			pkgMap["name"] = pkg.Description
		} else {
			pkgMap["name"] = ""
		}
		if pkg.DescriptionID != "" {
			pkgMap["description"] = pkg.DescriptionID
		} else if pkg.Description != "" {
			// Legacy fallback
			pkgMap["description"] = pkg.Description
		} else {
			pkgMap["description"] = ""
		}
		if len(pkg.RoutesID) > 0 {
			routes = pkg.RoutesID
		} else if len(pkg.Routes) > 0 {
			// Legacy fallback
			routes = pkg.Routes
		} else {
			routes = []models.RouteDetail{}
		}
		if len(pkg.FeaturesID) > 0 {
			pkgMap["features"] = pkg.FeaturesID
		} else if len(pkg.Features) > 0 {
			// Legacy fallback
			pkgMap["features"] = pkg.Features
		} else {
			pkgMap["features"] = []string{}
		}
		if len(pkg.ExcludesID) > 0 {
			pkgMap["excludes"] = pkg.ExcludesID
		} else if len(pkg.Excludes) > 0 {
			// Legacy fallback
			pkgMap["excludes"] = pkg.Excludes
		} else {
			pkgMap["excludes"] = []string{}
		}
	}
	
	// Transform routes to match frontend format
	transformedRoutes := make([]map[string]string, len(routes))
	for j, route := range routes {
		routeMap := map[string]string{
			"price": route.Price,
		}
		if lang == "en" {
			if route.NameEN != "" {
				routeMap["name"] = route.NameEN
			} else {
				routeMap["name"] = route.NameID
			}
		} else {
			if route.NameID != "" {
				routeMap["name"] = route.NameID
			} else {
				routeMap["name"] = ""
			}
		}
		transformedRoutes[j] = routeMap
	}
	pkgMap["routes"] = transformedRoutes

	utils.APIResponse(c, http.StatusOK, "Berhasil mengambil detail Package", pkgMap)
}

func CreatePackage(c *gin.Context) {
	var input models.Package
	
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Format JSON salah: "+err.Error())
		return
	}

	// Ensure multi-language fields are set from legacy fields if needed
	// If NameID is empty but Name is provided, use Name as NameID
	if input.NameID == "" && input.Name != "" {
		input.NameID = input.Name
	}
	// If DescriptionID is empty but Description is provided, use Description as DescriptionID
	if input.DescriptionID == "" && input.Description != "" {
		input.DescriptionID = input.Description
	}
	// If RoutesID is empty but Routes is provided, convert Routes to RoutesID
	if len(input.RoutesID) == 0 && len(input.Routes) > 0 {
		convertedRoutes := make([]models.RouteDetail, len(input.Routes))
		for i, r := range input.Routes {
			routeName := ""
			if r.NameID != "" {
				routeName = r.NameID
			} else if r.NameEN != "" {
				routeName = r.NameEN
			}
			convertedRoutes[i] = models.RouteDetail{
				NameID: routeName,
				NameEN: routeName,
				Price:  r.Price,
			}
		}
		input.RoutesID = convertedRoutes
	}
	// If FeaturesID is empty but Features is provided, use Features as FeaturesID
	if len(input.FeaturesID) == 0 && len(input.Features) > 0 {
		input.FeaturesID = input.Features
	}
	// If ExcludesID is empty but Excludes is provided, use Excludes as ExcludesID
	if len(input.ExcludesID) == 0 && len(input.Excludes) > 0 {
		input.ExcludesID = input.Excludes
	}
	
	// Ensure legacy fields are set for backward compatibility
	if input.Name == "" && input.NameID != "" {
		input.Name = input.NameID
	}
	if input.Description == "" && input.DescriptionID != "" {
		input.Description = input.DescriptionID
	}
	if len(input.Routes) == 0 && len(input.RoutesID) > 0 {
		// Convert RoutesID to Routes for legacy support
		input.Routes = input.RoutesID
	}
	if len(input.Features) == 0 && len(input.FeaturesID) > 0 {
		input.Features = input.FeaturesID
	}
	if len(input.Excludes) == 0 && len(input.ExcludesID) > 0 {
		input.Excludes = input.ExcludesID
	}

	result := config.DB.Create(&input)
	if result.Error != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal Menyimpan ke database: "+result.Error.Error())
		return
	}
	utils.APIResponse(c, http.StatusOK, "Berhasil membuat package baru", input)
}

func UpdatePackage(c *gin.Context) {
	id := c.Param("id")
	var pkg models.Package

	// 1. Cari dulu apakah ada
	if err := config.DB.First(&pkg, id).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "Package tidak ditemukan")
		return
	}

	// 2. Bind JSON baru ke variabel input
	var input models.Package
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Format JSON salah: "+err.Error())
		return
	}

	// 3. Update multi-language fields (always update, even if empty)
	// Check if field is present in JSON by checking if it's explicitly set
	// For now, we'll update if provided (non-zero value) or if explicitly set to empty
	
	// Update multi-language fields - always update if provided
	// Always update these fields if they are in the input
	pkg.NameID = input.NameID
	pkg.NameEN = input.NameEN
	pkg.DescriptionID = input.DescriptionID
	pkg.DescriptionEN = input.DescriptionEN
	
	// Update arrays - always replace if provided (even if empty array)
	if input.RoutesID != nil {
		pkg.RoutesID = input.RoutesID
	}
	if input.RoutesEN != nil {
		pkg.RoutesEN = input.RoutesEN
	}
	if input.FeaturesID != nil {
		pkg.FeaturesID = input.FeaturesID
	}
	if input.FeaturesEN != nil {
		pkg.FeaturesEN = input.FeaturesEN
	}
	if input.ExcludesID != nil {
		pkg.ExcludesID = input.ExcludesID
	}
	if input.ExcludesEN != nil {
		pkg.ExcludesEN = input.ExcludesEN
	}
	
	// Legacy fields for backward compatibility
	if input.Name != "" {
		pkg.Name = input.Name
		// If NameID is empty, use Name as NameID
		if pkg.NameID == "" {
			pkg.NameID = input.Name
		}
	}
	if input.Description != "" {
		pkg.Description = input.Description
		// If DescriptionID is empty, use Description as DescriptionID
		if pkg.DescriptionID == "" {
			pkg.DescriptionID = input.Description
		}
	}
	if len(input.Routes) > 0 {
		// Convert Routes to RoutesID if RoutesID is empty
		if len(pkg.RoutesID) == 0 && len(input.Routes) > 0 {
			convertedRoutes := make([]models.RouteDetail, len(input.Routes))
			for i, r := range input.Routes {
				routeName := ""
				if r.NameID != "" {
					routeName = r.NameID
				} else if r.NameEN != "" {
					routeName = r.NameEN
				}
				convertedRoutes[i] = models.RouteDetail{
					NameID: routeName,
					NameEN: routeName,
					Price:  r.Price,
				}
			}
			pkg.RoutesID = convertedRoutes
		}
		if len(input.Routes) > 0 {
			pkg.Routes = input.Routes
		}
	}
	if len(input.Features) > 0 {
		// If FeaturesID is empty, use Features as FeaturesID
		if len(pkg.FeaturesID) == 0 {
			pkg.FeaturesID = input.Features
		}
		pkg.Features = input.Features
	}
	if len(input.Excludes) > 0 {
		// If ExcludesID is empty, use Excludes as ExcludesID
		if len(pkg.ExcludesID) == 0 {
			pkg.ExcludesID = input.Excludes
		}
		pkg.Excludes = input.Excludes
	}
	
	// Common fields (always update)
	if input.Capacity != "" {
		pkg.Capacity = input.Capacity
	}
	if input.Duration != "" {
		pkg.Duration = input.Duration
	}
	pkg.IsPopular = input.IsPopular
	// Update image_url if provided (even if empty string)
	// This allows updating image_url while preserving existing value if not provided
	// Note: If you want to clear image_url, send empty string explicitly
	if input.ImageURL != "" {
		pkg.ImageURL = input.ImageURL
	}
	// If ImageURL is empty string in input, we don't update it to preserve existing image
	// This prevents accidentally clearing the image_url when updating other fields
	// IMPORTANT: image_url is preserved if not provided or if empty string is sent

	if err := config.DB.Save(&pkg).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal mengupdate database: "+err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Berhasil update package", pkg)
}

func DeletePackage(c *gin.Context) {
	id := c.Param("id")
	var pkg models.Package

	if err := config.DB.First(&pkg, id).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "Package tidak ditemukan")
		return
	}

	if err := config.DB.Unscoped().Delete(&pkg).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal menghapus data")
		return
	}

	utils.APIResponse(c, http.StatusOK, "Berhasil menghapus package", nil)
}


func ForceResetPackage(c *gin.Context) {
    // 1. Matikan Sensor Keamanan (Foreign Key Check)
    // Supaya database tidak protes kalau kita hapus paksa induk datanya
    config.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")

    // 2. Lakukan Factory Reset (Truncate)
    // Menghapus semua data & mereset ID kembali ke 1
    if err := config.DB.Exec("TRUNCATE TABLE packages").Error; err != nil {
        // Jangan lupa nyalakan lagi sensornya kalau error
        config.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
        utils.APIError(c, http.StatusInternalServerError, "Gagal reset tabel: " + err.Error())
        return
    }

    // 3. Nyalakan Kembali Sensor Keamanan
    config.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

    utils.APIResponse(c, http.StatusOK, "Tabel berhasil di-reset paksa ke ID 1", nil)
}