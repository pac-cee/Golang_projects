package config

import (
	"os"
)

type Config struct {
	GRPCPort string
	MongoURI string
	MongoDB  string
	RedisURI string
}

func LoadConfig() *Config {
	return &Config{
		GRPCPort: getEnv("GRPC_PORT", ":50051"),
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "ecommerce"),
		RedisURI: getEnv("REDIS_URI", "redis://localhost:6379/0"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
