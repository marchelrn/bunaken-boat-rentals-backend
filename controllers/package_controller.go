package controllers

import (
	"bunaken-boat-backend/config"
	"bunaken-boat-backend/models"
	"bunaken-boat-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllPackages(c *gin.Context) {
	packages := []models.Package{} 
	if err := config.DB.Find(&packages).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal Mengambil database")
		return
	}
	
	lang := c.DefaultQuery("lang", "id")
	if lang != "id" && lang != "en" {
		lang = "id"
	}
	
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
			"name_id":   pkg.NameID,
			"name_en":   pkg.NameEN,
			"routes_id": pkg.RoutesID,
			"routes_en":      pkg.RoutesEN,
			"features_id":    pkg.FeaturesID,
			"features_en":    pkg.FeaturesEN,
			"includes_id":    pkg.IncludesID,
			"includes_en":    pkg.IncludesEN,
			"excludes_id":    pkg.ExcludesID,
			"excludes_en":    pkg.ExcludesEN,
			"name":        pkg.Name,
			"routes":      pkg.Routes,
			"features":    pkg.Features,
			"includes":    pkg.Includes,
			"excludes":    pkg.Excludes,
		}
		
		if lang == "en" {
			if pkg.NameEN != "" {
				pkgMap["name"] = pkg.NameEN
			} else {
				pkgMap["name"] = pkg.NameID
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
			if len(pkg.IncludesEN) > 0 {
				pkgMap["includes"] = pkg.IncludesEN
			} else {
				pkgMap["includes"] = pkg.IncludesID
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
				pkgMap["name"] = pkg.Name	
			} else {
				pkgMap["name"] = ""
			}
			if len(pkg.FeaturesID) > 0 {
				pkgMap["features"] = pkg.FeaturesID
			} else if len(pkg.Features) > 0 {
				pkgMap["features"] = pkg.Features
			} else {
				pkgMap["features"] = []string{}
			}
			if len(pkg.ExcludesID) > 0 {
				pkgMap["excludes"] = pkg.ExcludesID
			} else if len(pkg.Excludes) > 0 {
				pkgMap["excludes"] = pkg.Excludes
			} else {
				pkgMap["excludes"] = []string{}
			}
		}
		
		var routes []models.RouteDetail
		if lang == "en" {
			if len(pkg.RoutesEN) > 0 {
				routes = pkg.RoutesEN
			} else if len(pkg.RoutesID) > 0 {
				routes = pkg.RoutesID
			} else if len(pkg.Routes) > 0 {
				routes = pkg.Routes
			} else {
				routes = []models.RouteDetail{}
			}
		} else {
			if len(pkg.RoutesID) > 0 {
				routes = pkg.RoutesID
			} else if len(pkg.Routes) > 0 {
				routes = pkg.Routes
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

	lang := c.DefaultQuery("lang", "id")
	if lang != "id" && lang != "en" {
		lang = "id"
	}
	
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
	
	var routes []models.RouteDetail
	if lang == "en" {
		if pkg.NameEN != "" {
			pkgMap["name"] = pkg.NameEN
		} else {
			pkgMap["name"] = pkg.NameID
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
		if len(pkg.IncludesEN) > 0 {
			pkgMap["includes"] = pkg.IncludesEN
		} else {
			pkgMap["includes"] = pkg.IncludesID
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
			pkgMap["name"] = pkg.Name
		} else {
			pkgMap["name"] = ""
		}
		if len(pkg.RoutesID) > 0 {
			routes = pkg.RoutesID
		} else if len(pkg.Routes) > 0 {
			routes = pkg.Routes
		} else {
			routes = []models.RouteDetail{}
		}
		if len(pkg.FeaturesID) > 0 {
			pkgMap["features"] = pkg.FeaturesID
		} else if len(pkg.Features) > 0 {
			pkgMap["features"] = pkg.Features
		} else {
			pkgMap["features"] = []string{}
		}
		if len(pkg.IncludesID) > 0 {
			pkgMap["includes"] = pkg.IncludesID
		} else if len(pkg.Includes) > 0 {
			pkgMap["includes"] = pkg.Includes
		} else {
			pkgMap["includes"] = []string{}
		}
		if len(pkg.ExcludesID) > 0 {
			pkgMap["excludes"] = pkg.ExcludesID
		} else if len(pkg.Excludes) > 0 {
			pkgMap["excludes"] = pkg.Excludes
		} else {
			pkgMap["excludes"] = []string{}
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

	utils.APIResponse(c, http.StatusOK, "Berhasil mengambil detail Package", pkgMap)
}

func CreatePackage(c *gin.Context) {
	var input models.Package
	
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Format JSON salah: "+err.Error())
		return
	}

	if input.NameID == "" && input.Name != "" {
		input.NameID = input.Name
	}
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
	if len(input.FeaturesID) == 0 && len(input.Features) > 0 {
		input.FeaturesID = input.Features
	}
	if len(input.IncludesID) == 0 && len(input.Includes) > 0 {
		input.IncludesID = input.Includes
	}
	if len(input.ExcludesID) == 0 && len(input.Excludes) > 0 {
		input.ExcludesID = input.Excludes
	}
	
	if input.Name == "" && input.NameID != "" {
		input.Name = input.NameID
	}
	if len(input.Routes) == 0 && len(input.RoutesID) > 0 {
		input.Routes = input.RoutesID
	}
	if len(input.Features) == 0 && len(input.FeaturesID) > 0 {
		input.Features = input.FeaturesID
	}
	if len(input.Includes) == 0 && len(input.IncludesID) > 0 {
		input.Includes = input.IncludesID
	}
	if len(input.Excludes) == 0 && len(input.ExcludesID) > 0 {
		input.Excludes = input.ExcludesID
	}
	
	if input.ImageURL == "" {
		input.ImageURL = ""
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

	if err := config.DB.First(&pkg, id).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "Package tidak ditemukan")
		return
	}

	var input models.Package
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Format JSON salah: "+err.Error())
		return
	}

	pkg.NameID = input.NameID
	pkg.NameEN = input.NameEN
	
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
	if input.IncludesID != nil {
		pkg.IncludesID = input.IncludesID
	}
	if input.IncludesEN != nil {
		pkg.IncludesEN = input.IncludesEN
	}
	if input.ExcludesID != nil {
		pkg.ExcludesID = input.ExcludesID
	}
	if input.ExcludesEN != nil {
		pkg.ExcludesEN = input.ExcludesEN
	}
	
	if input.Name != "" {
		pkg.Name = input.Name
		if pkg.NameID == "" {
			pkg.NameID = input.Name
		}
	}
	if len(input.Routes) > 0 {
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
		if len(pkg.FeaturesID) == 0 {
			pkg.FeaturesID = input.Features
		}
		pkg.Features = input.Features
	}
	if len(input.Excludes) > 0 {
		if len(pkg.ExcludesID) == 0 {
			pkg.ExcludesID = input.Excludes
		}
		pkg.Excludes = input.Excludes
	}
	
	if input.Capacity != "" {
		pkg.Capacity = input.Capacity
	}
	if input.Duration != "" {
		pkg.Duration = input.Duration
	}
	pkg.IsPopular = input.IsPopular
	
	existingImageURL := pkg.ImageURL
	if input.ImageURL != "" {
		pkg.ImageURL = input.ImageURL
	} else {
		pkg.ImageURL = existingImageURL
	}

	if err := config.DB.Model(&pkg).Select(
		"name_id", "name_en", "capacity", "duration", "is_popular", "image_url",
		"routes_id", "routes_en", "features_id", "features_en",
		"includes_id", "includes_en", "excludes_id", "excludes_en",
		"name", "routes", "features", "includes", "excludes",
	).Updates(&pkg).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal mengupdate database: "+err.Error())
		return
	}

	if err := config.DB.First(&pkg, id).Error; err != nil {
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