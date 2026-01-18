package main

import (
	"log"

	"github.com/sayeed1999/ride-sharing-golang-api/config"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/database"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/auth"
	"github.com/sayeed1999/ride-sharing-golang-api/internal/modules/trip"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db := database.InitDBWithErrorHandling(cfg)
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}
	defer database.CloseDBWithErrorHandling(db)

	// Run database migrations
	database.RunMigrationsWithErrorHandling(sqlDB, cfg)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Create API router group
	api := router.Group("/api")
	{
		// Auth routes under /api/auth/
		authGroup := api.Group("/auth")
		auth.ExposeRoutes(authGroup, db, cfg)

		// Trip routes under /api/trip/
		tripGroup := api.Group("/trip")
		trip.ExposeRoutes(tripGroup, db, cfg)
	}

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
