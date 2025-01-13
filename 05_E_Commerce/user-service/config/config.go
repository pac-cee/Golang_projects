package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	GRPCPort    int
	MongoDBURI  string
	MongoDBName string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50053" // Default port for user service
	}

	grpcPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid GRPC_PORT: %v", err)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoDBName := os.Getenv("MONGODB_DATABASE")
	if mongoDBName == "" {
		mongoDBName = "ecommerce"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-256-bit-secret"
	}

	return &Config{
		GRPCPort:    grpcPort,
		MongoDBURI:  mongoURI,
		MongoDBName: mongoDBName,
		JWTSecret:   jwtSecret,
	}, nil
}
