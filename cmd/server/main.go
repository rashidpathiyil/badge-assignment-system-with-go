package main

import (
	"log"
	"os"

	"github.com/badge-assignment-system/internal/api"
	"github.com/badge-assignment-system/internal/models"
	"github.com/badge-assignment-system/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to the database
	db, err := models.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create service layer
	svc := service.NewService(db)

	// Set up the HTTP server
	router := setupServer(svc)

	// Get the port to listen on
	port := getEnv("PORT", "8080")

	// Start the server
	log.Printf("Server starting on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupServer configures the HTTP server
func setupServer(svc *service.Service) *gin.Engine {
	// Set Gin mode
	mode := getEnv("GIN_MODE", "debug")
	gin.SetMode(mode)

	// Create a new Gin router
	router := gin.New()

	// Use logger and recovery middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Create API handler
	handler := api.NewHandler(svc)

	// Set up routes
	api.SetupRoutes(router, handler)

	return router
}

// getEnv gets an environment variable or a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
