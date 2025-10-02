package main

import (
	"h3-travel/config"
	"h3-travel/models"
	"h3-travel/routes"
	"log"

	_ "h3-travel/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title H3 Travel API
// @version 1.0
// @description API pour le projet H3 Travel
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Connexion DB
	config.ConnectDatabase()

	// Migration des modèles
	config.DB.AutoMigrate(&models.User{}, &models.Travel{}, &models.Order{})

	// Routes
	r := routes.SetupRouter()

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("🚀 Server running on port 8080")
	r.Run(":8080")
}
