package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rishuishind/rivalTaskAssestment/backend/config"
	"github.com/rishuishind/rivalTaskAssestment/backend/middleware"
	"github.com/rishuishind/rivalTaskAssestment/backend/models"
	"github.com/rishuishind/rivalTaskAssestment/backend/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Connect to database
	config.ConnectDatabase()

	// Auto-migrate models
	if err := config.DB.AutoMigrate(&models.User{}, &models.Task{}); err != nil {
		log.Fatal("Failed to auto-migrate database: ", err)
	}
	log.Println("✅ Database migrated successfully")

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	// Create Gin router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(middleware.SetupCORS())

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
