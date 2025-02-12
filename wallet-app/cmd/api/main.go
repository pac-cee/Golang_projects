package main

import (
    "log"
    "net/http"
    "os"
    "wallet-app/internal/config"
    "wallet-app/internal/handlers"
    "wallet-app/internal/middleware"
    "wallet-app/pkg/database"

    "github.com/gin-gonic/gin"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Initialize database
    db, err := database.NewPostgresDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Initialize router
    router := gin.Default()
    router.Use(middleware.Auth())

    // Initialize handlers
    walletHandler := handlers.NewWalletHandler(db)

    // Define routes
    api := router.Group("/api/v1")
    {
        api.POST("/wallets", walletHandler.Create)
        api.GET("/wallets/:id", walletHandler.Get)
        api.POST("/wallets/:id/deposit", walletHandler.Deposit)
        api.POST("/wallets/:id/withdraw", walletHandler.Withdraw)
        api.GET("/wallets/:id/transactions", walletHandler.GetTransactions)
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server starting on port %s", port)
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatal(err)
    }
}