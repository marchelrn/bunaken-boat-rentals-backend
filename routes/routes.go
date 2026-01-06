package routes

import (
	"bunaken-boat-backend/controllers"
	"bunaken-boat-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true 
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		// Auth Routes (Public)
		api.POST("/auth/register", controllers.Register) // Gunakan sekali buat bikin admin pertama, lalu bisa dikomentari/diamankan
		api.POST("/auth/login", controllers.Login)

		// Public Package Routes
		api.GET("/packages", controllers.GetAllPackages)
		api.GET("/packages/:id", controllers.GetPackageByID)
		
		// Protected Routes (Butuh Token)
		protected := api.Group("/admin")
		protected.Use(middleware.JwtAuthMiddleware())
		{
			protected.POST("/packages", controllers.CreatePackage)
			protected.PUT("/packages/:id", controllers.UpdatePackage)
			protected.DELETE("/packages/:id", controllers.DeletePackage)
		}
	}

	return r
}
