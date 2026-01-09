package config

import (
	"bunaken-boat-backend/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	
	if dsn == "" {
		// Ambil sslmode dari env, default require untuk production
		sslMode := os.Getenv("DB_SSLMODE")
		if sslMode == "" {
			// Jika DB_HOST bukan localhost, assume production
			dbHost := os.Getenv("DB_HOST")
			if dbHost != "localhost" && dbHost != "127.0.0.1" && dbHost != "" {
				sslMode = "require"
			} else {
				sslMode = "disable"
			}
		}
		
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")
		
		if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
			log.Fatal("❌ Database configuration missing! Please set DATABASE_URL or DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT")
		}
		
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			dbHost,
			dbUser,
			dbPassword,
			dbName,
			dbPort,
			sslMode,
		)
	}
	
	log.Println("🔄 Connecting to database...")
	log.Printf("📝 Database host: %s", os.Getenv("DB_HOST"))
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
		panic("Koneksi Database Gagal! :" + err.Error())
	}
	
	log.Println("✅ Database connected successfully")
	log.Println("🔄 Running database migrations...")
	
	// AutoMigrate
	if err := database.AutoMigrate(&models.Package{}, &models.User{}, &models.AddOn{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
		panic("Migration Gagal! :" + err.Error())
	}
	
	log.Println("✅ Database migrations completed successfully")
	log.Println("📊 Tables: packages, users, addons")

	DB = database
}