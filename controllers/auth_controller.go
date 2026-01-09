package controllers

import (
	"bunaken-boat-backend/config"
	"bunaken-boat-backend/models"
	"bunaken-boat-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal mengenkripsi password")
		return
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.APIError(c, http.StatusBadRequest, "Username mungkin sudah digunakan")
		return
	}

	utils.APIResponse(c, http.StatusOK, "Registrasi berhasil", nil)
}

func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.APIError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		utils.APIError(c, http.StatusBadRequest, "Username atau password salah")
		return
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.APIError(c, http.StatusBadRequest, "Username atau password salah")
		return
	}

	// Generate Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal generate token")
		return
	}

	utils.APIResponse(c, http.StatusOK, "Login berhasil", gin.H{"token": token, "role": user.Role})
}
