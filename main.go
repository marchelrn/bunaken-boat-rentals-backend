package main

import (
	"bunaken-boat-backend/config"
	"bunaken-boat-backend/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env hanya untuk development
	// Di production, gunakan environment variables dari platform
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️  .env file not found (this is OK in production)")
		}
	}
	
	// Validasi API_SECRET
	apiSecret := os.Getenv("API_SECRET")
	if apiSecret == "" {
		log.Fatal("❌ API_SECRET environment variable is required!")
	}
	log.Println("✅ API_SECRET configured")
	
	config.ConnectDatabase()

	r := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("🚀 Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}