package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pac-cee/student-activity-tracker/config"
	"github.com/pac-cee/student-activity-tracker/database"
	"github.com/pac-cee/student-activity-tracker/handlers"
	"github.com/pac-cee/student-activity-tracker/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewMongoDB(cfg.MongoURI, cfg.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize services
	activityService := services.NewActivityService(db.Collection("activities"))

	// Initialize handlers
	activityHandler := handlers.NewActivityHandler(activityService)

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.AllowedOrigins
	router.Use(cors.New(corsConfig))

	// Register routes
	activityHandler.RegisterRoutes(router)

	// Start server
	log.Printf("Server starting on %s", cfg.ServerAddress)
	log.Fatal(router.Run(cfg.ServerAddress))
}
