package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Middleware function to log requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// JSON response helper
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// HomeHandler handles the home route
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		respondWithJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: "Route not found",
		})
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Welcome to the Go Web Server",
	})
}

// HealthCheckHandler handles the health check endpoint
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Service is healthy",
		Data: map[string]string{
			"version": "1.0.0",
			"uptime":  time.Now().String(),
		},
	})
}

func main() {
	// Create a new mux router
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/health", HealthCheckHandler)

	// Wrap the mux with the logging middleware
	handler := loggingMiddleware(mux)

	// Configure the server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
