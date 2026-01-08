package controllers

import (
	"bunaken-boat-backend/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadPackageImage handles image upload for packages
func UploadPackageImage(c *gin.Context) {
	// Get the file from form data
	file, err := c.FormFile("image")
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "File tidak ditemukan: "+err.Error())
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	isAllowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		utils.APIError(c, http.StatusBadRequest, "Format file tidak didukung. Gunakan: jpg, jpeg, png, gif, atau webp")
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		utils.APIError(c, http.StatusBadRequest, "Ukuran file terlalu besar. Maksimal 5MB")
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "/uploads/packages"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal membuat direktori upload: "+err.Error())
		return
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("package_%d%s", timestamp, ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save the file
	src, err := file.Open()
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal membuka file: "+err.Error())
		return
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error())
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal menyalin file: "+err.Error())
		return
	}

	// Return the URL (relative path that can be served statically)
	// In production, you might want to return a full URL like: https://yourdomain.com/uploads/packages/filename.jpg
	imageURL := fmt.Sprintf("/uploads/packages/%s", filename)
	
	utils.APIResponse(c, http.StatusOK, "Gambar berhasil diupload", gin.H{
		"image_url": imageURL,
		"filename":  filename,
	})
}

