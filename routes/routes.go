package routes

import (
	"bunaken-boat-backend/controllers"
	"bunaken-boat-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Serve static files for uploads
	r.Static("/uploads", "./uploads")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true 
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		// Auth Routes (Public)
		// api.POST("/auth/register", controllers.Register) // Dinonaktifkan untuk keamanan - gunakan hanya untuk membuat admin pertama
		api.POST("/auth/login", controllers.Login)

		// Public Package Routes
		api.GET("/packages", controllers.GetAllPackages)
		api.GET("/packages/:id", controllers.GetPackageByID)
		
		// Public Add-On Routes
		api.GET("/addons", controllers.GetAllAddOns)
		api.GET("/addons/:id", controllers.GetAddOnByID)
		
		// Protected Routes (Butuh Token)
		protected := api.Group("/admin")
		protected.Use(middleware.JwtAuthMiddleware())
		{
			// Package Routes
			protected.POST("/packages", controllers.CreatePackage)
			protected.PUT("/packages/:id", controllers.UpdatePackage)
			protected.DELETE("/packages/:id", controllers.DeletePackage)
			protected.POST("/packages/upload-image", controllers.UploadPackageImage)
			
			// Add-On Routes
			protected.POST("/addons", controllers.CreateAddOn)
			protected.PUT("/addons/:id", controllers.UpdateAddOn)
			protected.DELETE("/addons/:id", controllers.DeleteAddOn)
			
			// Auth Routes
			protected.PUT("/change-password", controllers.ChangePassword)
		}
	}

	return r
}
