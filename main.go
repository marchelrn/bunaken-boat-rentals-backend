package main

import (
	"bunaken-boat-backend/config"
	"bunaken-boat-backend/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("🚀 Starting application...")
	
	// Load .env hanya untuk development
	// Di production, gunakan environment variables dari platform
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️  .env file not found (this is OK in production)")
		} else {
			log.Println("✅ .env file loaded")
		}
	}
	
	// Validasi API_SECRET
	apiSecret := os.Getenv("API_SECRET")
	if apiSecret == "" {
		log.Fatal("❌ API_SECRET environment variable is required! Please set API_SECRET in environment variables.")
	}
	log.Println("✅ API_SECRET configured")
	
	// Connect to database
	log.Println("🔄 Initializing database connection...")
	config.ConnectDatabase()
	log.Println("✅ Database initialized")

	// Setup router
	log.Println("🔄 Setting up routes...")
	r := routes.SetupRouter()
	log.Println("✅ Routes configured")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("⚠️  PORT not set, using default: 8080")
	} else {
		log.Printf("📝 Using PORT from environment: %s", port)
	}
	
	log.Printf("🚀 Server starting on port %s", port)
	log.Println("✅ Application is ready to accept requests")
	
	// Start server - this will block
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}
}