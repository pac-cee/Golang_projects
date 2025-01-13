package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	GRPCPort     int
	MongoDBURI   string
	MongoDBName  string
	ProductSvcURL string
	UserSvcURL   string
}

func LoadConfig() (*Config, error) {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052" // Default port for order service
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

	productSvcURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productSvcURL == "" {
		productSvcURL = "localhost:50051"
	}

	userSvcURL := os.Getenv("USER_SERVICE_URL")
	if userSvcURL == "" {
		userSvcURL = "localhost:50053"
	}

	return &Config{
		GRPCPort:     grpcPort,
		MongoDBURI:   mongoURI,
		MongoDBName:  mongoDBName,
		ProductSvcURL: productSvcURL,
		UserSvcURL:   userSvcURL,
	}, nil
}
