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

func UploadPackageImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		utils.APIError(c, http.StatusBadRequest, "File tidak ditemukan: "+err.Error())
		return
	}

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

	if file.Size > 5*1024*1024 {
		utils.APIError(c, http.StatusBadRequest, "Ukuran file terlalu besar. Maksimal 5MB")
		return
	}

	src, err := file.Open()
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal membuka file: "+err.Error())
		return
	}
	defer src.Close()

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	isProduction := os.Getenv("ENV") == "production" || os.Getenv("PORT") != ""

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		if isProduction {
			utils.APIError(c, http.StatusInternalServerError, "Cloudinary tidak dikonfigurasi. Silakan set environment variables: CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET")
			return
		}

		uploadDir := "./uploads/packages"
		
		if _, err := os.Stat("."); err != nil {
			tempDir := os.TempDir()
			uploadDir = filepath.Join(tempDir, "uploads", "packages")
		}

		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Gagal membuat direktori upload. Untuk production, silakan konfigurasi Cloudinary: "+err.Error())
			return
		}

		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("package_%d%s", timestamp, ext)
		filePath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(filePath)
		if err != nil {
			utils.APIError(c, http.StatusInternalServerError, "Gagal menyimpan file. Untuk production, silakan konfigurasi Cloudinary: "+err.Error())
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

		imageURL := fmt.Sprintf("/uploads/packages/%s", filename)

		utils.APIResponse(c, http.StatusOK, "Gambar berhasil diupload (local storage - development only)", gin.H{
			"image_url": imageURL,
			"filename":  filename,
		})
		return
	}

	ctx := context.Background()
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal menginisialisasi Cloudinary: "+err.Error())
		return
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("package_%d%s", timestamp, ext)

	if _, err := src.Seek(0, 0); err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal membaca file: "+err.Error())
		return
	}

	uploadResult, err := cld.Upload.Upload(ctx, src, uploader.UploadParams{
		PublicID: fmt.Sprintf("bunaken-boat/packages/%s", filename),
		Folder:   "bunaken-boat/packages",
		ResourceType: "image",
	})
	if err != nil {
		utils.APIError(c, http.StatusInternalServerError, "Gagal upload ke Cloudinary: "+err.Error())
		return
	}

	imageURL := uploadResult.SecureURL

	utils.APIResponse(c, http.StatusOK, "Gambar berhasil diupload", gin.H{
		"image_url": imageURL,
		"filename":  filename,
		"public_id": uploadResult.PublicID,
	})
}

