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
	utils.APIResponse(c, http.StatusOK, "Berhasil mengambil data Packages", packages)
}

func GetPackageByID(c *gin.Context) {
	id := c.Param("id")
	var pkg models.Package

	if err := config.DB.First(&pkg, id).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "Package tidak ditemukan")
		return
	}

	utils.APIResponse(c, http.StatusOK, "Berhasil mengambil detail Package", pkg)
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