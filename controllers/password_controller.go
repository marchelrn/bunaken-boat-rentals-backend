package controllers

import (
	"bunaken-boat-backend/config"
	"bunaken-boat-backend/models"
	"bunaken-boat-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.APIError(c, http.StatusUnauthorized, "User ID tidak ditemukan")
		return
	}

	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Format JSON salah: "+err.Error())
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.APIError(c, http.StatusNotFound, "User tidak ditemukan")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Password lama salah")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal mengenkripsi password baru")
		return
	}

	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal mengupdate password: "+err.Error())
		return
	}

	utils.APIResponse(c, http.StatusOK, "Password berhasil diubah", nil)
}

