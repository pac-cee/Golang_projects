package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	DB       DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Logger   LoggerConfig
	Swagger  SwaggerConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Mode         string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	Timeout  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret           string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type LoggerConfig struct {
	Level string
	File  string
}

type SwaggerConfig struct {
	Enabled bool
	Host    string
	BasePath string
}

type RateLimitConfig struct {
	Enabled bool
	Requests int
	Duration time.Duration
}

// LoadConfig reads configuration from file or environment variables
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set defaults if not specified
	setDefaults(&cfg)

	return &cfg, nil
}

func setDefaults(cfg *Config) {
	// Server defaults
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 10 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 10 * time.Second
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "development"
	}

	// Database defaults
	if cfg.DB.MaxConns == 0 {
		cfg.DB.MaxConns = 10
	}
	if cfg.DB.Timeout == 0 {
		cfg.DB.Timeout = 5 * time.Second
	}
	if cfg.DB.SSLMode == "" {
		cfg.DB.SSLMode = "disable"
	}

	// Redis defaults
	if cfg.Redis.Port == 0 {
		cfg.Redis.Port = 6379
	}
	if cfg.Redis.DB < 0 {
		cfg.Redis.DB = 0
	}

	// JWT defaults
	if cfg.JWT.AccessTokenTTL == 0 {
		cfg.JWT.AccessTokenTTL = 15 * time.Minute
	}
	if cfg.JWT.RefreshTokenTTL == 0 {
		cfg.JWT.RefreshTokenTTL = 24 * time.Hour
	}

	// Logger defaults
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}

	// Rate limit defaults
	if cfg.RateLimit.Enabled && cfg.RateLimit.Requests == 0 {
		cfg.RateLimit.Requests = 100
		cfg.RateLimit.Duration = time.Minute
	}
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetRedisAddr returns the Redis connection address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
