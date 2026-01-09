package routes

import (
	"bunaken-boat-backend/controllers"
	"bunaken-boat-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/uploads", "/uploads")

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	r.Use(cors.New(config))

	api := r.Group("/api")
	{
		// when admin is registered just disable this endpoint
		// api.POST("/auth/register", controllers.Register) 
		api.POST("/auth/login", controllers.Login)

		api.GET("/packages", controllers.GetAllPackages)
		api.GET("/packages/:id", controllers.GetPackageByID)
		
		api.GET("/addons", controllers.GetAllAddOns)
		api.GET("/addons/:id", controllers.GetAddOnByID)
		
		protected := api.Group("/admin")
		protected.Use(middleware.JwtAuthMiddleware())
		{
			protected.POST("/packages", controllers.CreatePackage)
			protected.PUT("/packages/:id", controllers.UpdatePackage)
			protected.DELETE("/packages/:id", controllers.DeletePackage)
			protected.POST("/packages/upload-image", controllers.UploadPackageImage)
			
			protected.POST("/addons", controllers.CreateAddOn)
			protected.PUT("/addons/:id", controllers.UpdateAddOn)
			protected.DELETE("/addons/:id", controllers.DeleteAddOn)
			
			protected.PUT("/change-password", controllers.ChangePassword)
		}
	}

	return r
}
