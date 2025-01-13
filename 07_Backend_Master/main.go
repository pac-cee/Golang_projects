package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-master/config"
	"backend-master/pkg/auth"
	"backend-master/pkg/database"
	"backend-master/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(logger.Config{
		Level:   cfg.Logger.Level,
		File:    cfg.Logger.File,
		Console: true,
	})
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize PostgreSQL
	db, err := database.NewPostgreSQL(context.Background(), database.PostgresConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
		MaxConns: cfg.DB.MaxConns,
		Timeout:  cfg.DB.Timeout,
	})
	if err != nil {
		log.Fatal("Failed to connect to database", err, nil)
	}
	defer db.Close()

	// Initialize Redis
	redis, err := database.NewRedis(database.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Fatal("Failed to connect to Redis", err, nil)
	}
	defer redis.Close()

	// Initialize JWT manager
	tokenManager := auth.NewTokenManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// Initialize Gin
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Middleware
	router.Use(
		gin.Recovery(),
		middleware.Logger(log),
		middleware.CORS(),
	)

	if cfg.RateLimit.Enabled {
		router.Use(middleware.RateLimit(cfg.RateLimit.Requests, cfg.RateLimit.Duration))
	}

	// Initialize handlers
	handlers := handler.NewHandler(handler.Dependencies{
		Log:          log,
		DB:           db,
		Redis:        redis,
		TokenManager: tokenManager,
		Config:       cfg,
	})

	// Register routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Auth.Register)
			auth.POST("/login", handlers.Auth.Login)
			auth.POST("/refresh", handlers.Auth.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.Auth(tokenManager))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("", handlers.User.List)
				users.GET("/:id", handlers.User.Get)
				users.PUT("/:id", handlers.User.Update)
				users.DELETE("/:id", handlers.User.Delete)
			}

			// Product routes
			products := protected.Group("/products")
			{
				products.POST("", handlers.Product.Create)
				products.GET("", handlers.Product.List)
				products.GET("/:id", handlers.Product.Get)
				products.PUT("/:id", handlers.Product.Update)
				products.DELETE("/:id", handlers.Product.Delete)
			}
		}
	}

	// Swagger documentation
	if cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Health check
	router.GET("/health", handlers.HealthCheck)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server
	go func() {
		log.Info("Starting server", map[string]interface{}{
			"port": cfg.Server.Port,
			"mode": cfg.Server.Mode,
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", err, nil)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", err, nil)
	}

	log.Info("Server exited properly", nil)
}
