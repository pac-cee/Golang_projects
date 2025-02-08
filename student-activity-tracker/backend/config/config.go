package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	MongoURI      string
	DatabaseName  string
	ServerAddress string
	AllowedOrigins []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:  getEnv("DB_NAME", "student_tracker"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		AllowedOrigins: []string{getEnv("ALLOWED_ORIGIN", "http://localhost:4321")},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
