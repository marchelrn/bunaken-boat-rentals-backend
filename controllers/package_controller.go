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
			} else {
				// Legacy fallback - try to get from old Name field if exists
				pkgMap["name"] = ""
			}
			if pkg.DescriptionID != "" {
				pkgMap["description"] = pkg.DescriptionID
			} else {
				pkgMap["description"] = ""
			}
			if len(pkg.RoutesID) > 0 {
				pkgMap["routes"] = pkg.RoutesID
			} else {
				pkgMap["routes"] = []models.RouteDetail{}
			}
			if len(pkg.FeaturesID) > 0 {
				pkgMap["features"] = pkg.FeaturesID
			} else {
				pkgMap["features"] = []string{}
			}
			if len(pkg.ExcludesID) > 0 {
				pkgMap["excludes"] = pkg.ExcludesID
			} else {
				pkgMap["excludes"] = []string{}
			}
		}
		
		// Transform routes to match frontend format
		var routes []models.RouteDetail
		if lang == "en" {
			if len(pkg.RoutesEN) > 0 {
				routes = pkg.RoutesEN
			} else {
				routes = pkg.RoutesID
			}
		} else {
			if len(pkg.RoutesID) > 0 {
				routes = pkg.RoutesID
			} else {
				routes = []models.RouteDetail{}
			}
		}
		
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
		} else {
			pkgMap["name"] = ""
		}
		if pkg.DescriptionID != "" {
			pkgMap["description"] = pkg.DescriptionID
		} else {
			pkgMap["description"] = ""
		}
		if len(pkg.RoutesID) > 0 {
			routes = pkg.RoutesID
		} else {
			routes = []models.RouteDetail{}
		}
		if len(pkg.FeaturesID) > 0 {
			pkgMap["features"] = pkg.FeaturesID
		} else {
			pkgMap["features"] = []string{}
		}
		if len(pkg.ExcludesID) > 0 {
			pkgMap["excludes"] = pkg.ExcludesID
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

	// 3. Update field. 
	// Kita gunakan Model(&pkg).Updates(input) agar hanya field yang dikirim yang berubah,
	// TAPI hati-hati dengan field boolean false atau string kosong (zero values).
	// Jika ingin replace total, lebih aman assign manual atau gunakan map.
	
	// Sederhananya kita update field-field penting:
	pkg.Name = input.Name
	pkg.Description = input.Description
	pkg.Capacity = input.Capacity
	pkg.Duration = input.Duration
	pkg.IsPopular = input.IsPopular
	pkg.ImageURL = input.ImageURL
	pkg.Routes = input.Routes
	pkg.Features = input.Features
	pkg.Excludes = input.Excludes

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