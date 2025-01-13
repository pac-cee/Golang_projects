package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"product-service/config"
	"product-service/db"
	"product-service/handler"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize configuration
	cfg := config.LoadConfig()

	// Initialize MongoDB connection
	mongodb, err := db.NewMongoDB(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Close()

	// Initialize Redis for caching
	redis, err := db.NewRedisClient(cfg.RedisURI)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Create gRPC server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(handler.LoggingInterceptor),
	)

	// Register product service
	productHandler := handler.NewProductHandler(mongodb, redis)
	pb.RegisterProductServiceServer(server, productHandler)

	// Register reflection service on gRPC server
	reflection.Register(server)

	// Start listening
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down gRPC server...")
		server.GracefulStop()
	}()

	log.Printf("Product service starting on port %s", cfg.GRPCPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
