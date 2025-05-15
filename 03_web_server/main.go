package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./config"
	"./middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		jsonResponse(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: "Page not found",
		})
		return
	}

	jsonResponse(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Welcome to the Go Web Server API",
		Data: map[string]string{
			"version": "1.0.0",
			"status":  "running",
		},
	})
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Services page")
}

func setupServer(cfg *config.Config) *http.Server {
	// Register routes with middleware
	http.HandleFunc("/", middleware.LoggingMiddleware(homeHandler))
	http.HandleFunc("/about", middleware.LoggingMiddleware(aboutHandler))
	http.HandleFunc("/contact", middleware.LoggingMiddleware(contactHandler))
	http.HandleFunc("/services", middleware.LoggingMiddleware(servicesHandler))

	return &http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func main() {
	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load configuration
	cfg := config.NewConfig()

	// Create server
	server := setupServer(cfg)

	// Create channel for shutdown signals
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Printf("Server is starting on port %s...\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on port %s: %v\n", cfg.Port, err)
	}

	<-done
	log.Println("Server stopped")
}
