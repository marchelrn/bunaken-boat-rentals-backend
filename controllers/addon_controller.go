package controllers

import (
	"bunaken-boat-backend/config"
	"bunaken-boat-backend/models"
	"bunaken-boat-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllAddOns(c *gin.Context) {
	addOns := []models.AddOn{}
	if err := config.DB.Find(&addOns).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal mengambil data add-ons")
		return
	}
	utils.APIResponse(c, http.StatusOK, "Berhasil mengambil data Add-Ons", addOns)
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

	addOn.Name = input.Name
	addOn.Price = input.Price
	addOn.Description = input.Description

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

