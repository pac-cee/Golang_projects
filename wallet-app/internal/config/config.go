package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL     string `mapstructure:"DATABASE_URL"`
	Port            string `mapstructure:"PORT"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
	Environment     string `mapstructure:"ENVIRONMENT"`
	RedisURL        string `mapstructure:"REDIS_URL"`
	MaxTransactions int    `mapstructure:"MAX_TRANSACTIONS"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	config := &Config{
		Port:            getEnvOrDefault("PORT", "8080"),
		Environment:     getEnvOrDefault("ENVIRONMENT", "development"),
		MaxTransactions: 100,
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
