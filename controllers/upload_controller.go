package controllers

import (
	"bunaken-boat-backend/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// UploadPackageImage handles image upload for packages using Cloudinary
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

	// Open file
	src, err := file.Open()
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal membuka file: "+err.Error())
		return
	}
	defer src.Close()

	// Check if Cloudinary is configured
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	// If Cloudinary is not configured, fallback to local storage
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		// Fallback to local storage
		uploadDir := "./uploads/packages"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Gagal membuat direktori upload: "+err.Error())
			return
		}

		// Generate unique filename
		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("package_%d%s", timestamp, ext)
		filePath := filepath.Join(uploadDir, filename)

		// Save the file locally
		dst, err := os.Create(filePath)
		if err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error())
			return
		}
		defer dst.Close()

		if _, err := src.Seek(0, 0); err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Gagal membaca file: "+err.Error())
			return
		}

		if _, err := dst.ReadFrom(src); err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Gagal menyalin file: "+err.Error())
			return
		}

		// Return the URL (relative path)
		imageURL := fmt.Sprintf("/uploads/packages/%s", filename)

		utils.APIResponse(c, http.StatusOK, "Gambar berhasil diupload (local storage)", gin.H{
			"image_url": imageURL,
			"filename":  filename,
		})
		return
	}

	// Initialize Cloudinary
	ctx := context.Background()
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal menginisialisasi Cloudinary: "+err.Error())
		return
	}

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("package_%d%s", timestamp, ext)

	// Reset file reader to beginning
	if _, err := src.Seek(0, 0); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal membaca file: "+err.Error())
		return
	}

	// Upload to Cloudinary
	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: fmt.Sprintf("bunaken-boat/packages/%s", filename),
		Folder:   "bunaken-boat/packages",
		ResourceType: "image",
	})
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal upload ke Cloudinary: "+err.Error())
		return
	}

	// Return the secure URL from Cloudinary
	imageURL := uploadResult.SecureURL

	utils.APIResponse(c, http.StatusOK, "Gambar berhasil diupload", gin.H{
		"image_url": imageURL,
		"filename":  filename,
		"public_id": uploadResult.PublicID,
	})
}

