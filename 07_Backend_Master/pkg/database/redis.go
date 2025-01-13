package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis represents a Redis client
type Redis struct {
	Client *redis.Client
}

// RedisConfig holds the configuration for Redis connection
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedis creates a new Redis client
func NewRedis(cfg RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis{Client: client}, nil
}

// Close closes the Redis client
func (r *Redis) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

// Set stores a key-value pair with expiration
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.Client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a value by key and unmarshals it into the provided interface
func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil // Key doesn't exist
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a key
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// SetNX sets a key-value pair if the key doesn't exist (useful for locks)
func (r *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.Client.SetNX(ctx, key, data, expiration).Result()
}

// Lock implements a distributed lock
func (r *Redis) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return r.SetNX(ctx, fmt.Sprintf("lock:%s", key), true, expiration)
}

// Unlock releases a distributed lock
func (r *Redis) Unlock(ctx context.Context, key string) error {
	return r.Delete(ctx, fmt.Sprintf("lock:%s", key))
}

// CacheGet attempts to get a value from cache, if not found calls fetch function and caches result
func (r *Redis) CacheGet(ctx context.Context, key string, dest interface{}, fetch func() (interface{}, error), expiration time.Duration) error {
	// Try to get from cache
	err := r.Get(ctx, key, dest)
	if err != nil {
		return err
	}

	// If value exists in cache, return it
	if dest != nil {
		return nil
	}

	// If not in cache, fetch it
	value, err := fetch()
	if err != nil {
		return fmt.Errorf("failed to fetch value: %w", err)
	}

	// Store in cache
	if err := r.Set(ctx, key, value, expiration); err != nil {
		return fmt.Errorf("failed to cache value: %w", err)
	}

	// Copy fetched value to destination
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal fetched value: %w", err)
	}

	return json.Unmarshal(data, dest)
}
