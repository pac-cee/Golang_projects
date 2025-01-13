package main

import (
	"log"
	"net/http"
	"os"

	"book-api/config"
	"book-api/database"
	"book-api/handlers"
	"book-api/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Create router
	router := mux.NewRouter()

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	// Auth routes
	router.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// Book routes (protected by auth middleware)
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("/books", bookHandler.GetBooks).Methods("GET")
	api.HandleFunc("/books", bookHandler.CreateBook).Methods("POST")
	api.HandleFunc("/books/{id}", bookHandler.GetBook).Methods("GET")
	api.HandleFunc("/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	api.HandleFunc("/books/{id}", bookHandler.DeleteBook).Methods("DELETE")

	// Add logging middleware
	router.Use(middleware.LoggingMiddleware)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
