package main

import (
	"log"
	"net/http"
	"os"

	"chat-app/database"
	"chat-app/handlers"
	"chat-app/websocket"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	defer database.Close()

	// Create a new hub
	hub := websocket.NewHub()
	go hub.Run()

	// Create router
	router := mux.NewRouter()

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(hub)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// WebSocket endpoint
	router.HandleFunc("/ws/{roomID}", chatHandler.ServeWs)

	// API endpoints
	router.HandleFunc("/api/rooms", chatHandler.GetRooms).Methods("GET")
	router.HandleFunc("/api/rooms", chatHandler.CreateRoom).Methods("POST")
	router.HandleFunc("/api/rooms/{roomID}", chatHandler.GetRoom).Methods("GET")
	router.HandleFunc("/api/rooms/{roomID}/messages", chatHandler.GetMessages).Methods("GET")

	// Serve index.html for the root path
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Add CORS middleware
	handler := handlers.CorsMiddleware(router)

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
