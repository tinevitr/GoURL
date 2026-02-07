package main

import (
	"log"
	"gourl/db"
	"gourl/config"
	"gourl/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Redis client
	redisClient, err := database.InitRedis()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// Initialize handlers
	urlHandler := handlers.NewURLHandler(redisClient)
	statsHandler := handlers.NewStatsHandler(redisClient)

	// Setup router
	r := gin.Default()

	// API Routes
	api := r.Group("/api")
	{
		// POST untuk create short URL (dengan JSON body)
		api.POST("/create", urlHandler.CreateShortURL)
		// GET untuk create short URL (dengan query parameters)
		api.GET("/create", urlHandler.CreateShortURLViaGet)
	 	api.GET("/stats", statsHandler.GetAllStats)
	}

	// Redirect route
	r.GET("/:slug", urlHandler.RedirectURL)

	// Start server
	cfg := config.LoadConfig()
	log.Printf("Server starting on port %s...\n", cfg.ServerPort)
	log.Printf("Base URL: %s\n", cfg.BaseURL)
	
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
