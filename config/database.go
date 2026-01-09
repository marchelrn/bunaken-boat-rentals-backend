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
	log.Println("🔍 Checking database configuration...")
	
	dsn := os.Getenv("DATABASE_URL")
	
	if dsn == "" {
		log.Println("📝 DATABASE_URL not found, using individual DB_* variables...")
		
		// Ambil sslmode dari env, default require untuk production
		sslMode := os.Getenv("DB_SSLMODE")
		if sslMode == "" {
			// Jika DB_HOST bukan localhost, assume production
			dbHost := os.Getenv("DB_HOST")
			if dbHost != "localhost" && dbHost != "127.0.0.1" && dbHost != "" {
				sslMode = "require"
				log.Println("📝 Production mode detected, using sslmode=require")
			} else {
				sslMode = "disable"
				log.Println("📝 Development mode detected, using sslmode=disable")
			}
		} else {
			log.Printf("📝 Using DB_SSLMODE from environment: %s", sslMode)
		}
		
		dbHost := os.Getenv("DB_HOST")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbPort := os.Getenv("DB_PORT")
		
		log.Printf("🔍 DB_HOST: %s", dbHost)
		log.Printf("🔍 DB_USER: %s", dbUser)
		log.Printf("🔍 DB_NAME: %s", dbName)
		log.Printf("🔍 DB_PORT: %s", dbPort)
		
		if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
			log.Println("❌ Missing database configuration:")
			if dbHost == "" {
				log.Println("   - DB_HOST is missing")
			}
			if dbUser == "" {
				log.Println("   - DB_USER is missing")
			}
			if dbPassword == "" {
				log.Println("   - DB_PASSWORD is missing")
			}
			if dbName == "" {
				log.Println("   - DB_NAME is missing")
			}
			if dbPort == "" {
				log.Println("   - DB_PORT is missing")
			}
			log.Fatal("❌ Please set DATABASE_URL or all DB_* environment variables (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)")
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
		log.Println("✅ Database connection string built from individual variables")
	} else {
		log.Println("✅ DATABASE_URL found")
	}
	
	log.Println("🔄 Connecting to database...")
	
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	
	log.Println("✅ Database connected successfully")
	log.Println("🔄 Running database migrations...")
	
	// AutoMigrate
	if err := database.AutoMigrate(&models.Package{}, &models.User{}, &models.AddOn{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	
	log.Println("✅ Database migrations completed successfully")
	log.Println("📊 Tables created/verified: packages, users, addons")

	DB = database
}