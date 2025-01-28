package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// Caching example
func cacheExample(ctx context.Context, rdb *redis.Client) {
	cacheKey := "user:profile:123"
	
	// Try to get data from cache
	cachedData, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Cache miss - get data from "database" and cache it
		userData := `{"id": 123, "name": "John Doe", "email": "john@example.com"}`
		err = rdb.Set(ctx, cacheKey, userData, 1*time.Hour).Err()
		if err != nil {
			fmt.Printf("Failed to cache data: %v\n", err)
			return
		}
		fmt.Println("Data cached successfully")
		cachedData = userData
	} else if err != nil {
		fmt.Printf("Failed to get cached data: %v\n", err)
		return
	}
	
	fmt.Printf("User data: %s\n", cachedData)
}

// Rate limiting example
func rateLimitExample(ctx context.Context, rdb *redis.Client) {
	userID := "user123"
	action := "api_call"
	key := fmt.Sprintf("ratelimit:%s:%s", userID, action)
	
	// Check current count
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		fmt.Printf("Failed to increment counter: %v\n", err)
		return
	}
	
	// Set expiry for first request
	if count == 1 {
		rdb.Expire(ctx, key, 1*time.Hour)
	}
	
	// Check if rate limit exceeded
	if count > 10 {
		fmt.Println("Rate limit exceeded!")
		return
	}
	
	fmt.Printf("Request allowed. Count: %d/10\n", count)
}

// Session management example
func sessionExample(ctx context.Context, rdb *redis.Client) {
	sessionID := "sess:abc123"
	
	// Create session
	sessionData := map[string]interface{}{
		"user_id": 123,
		"username": "johndoe",
		"logged_in": true,
	}
	
	err := rdb.HSet(ctx, sessionID, sessionData).Err()
	if err != nil {
		fmt.Printf("Failed to create session: %v\n", err)
		return
	}
	
	// Set session expiry
	rdb.Expire(ctx, sessionID, 24*time.Hour)
	
	// Get session data
	data, err := rdb.HGetAll(ctx, sessionID).Result()
	if err != nil {
		fmt.Printf("Failed to get session: %v\n", err)
		return
	}
	
	fmt.Printf("Session data: %v\n", data)
}

// Pub/Sub example
func pubSubExample(ctx context.Context, rdb *redis.Client) {
	// Subscribe to channel
	pubsub := rdb.Subscribe(ctx, "notifications")
	defer pubsub.Close()
	
	// Start subscriber in goroutine
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				fmt.Printf("Failed to receive message: %v\n", err)
				return
			}
			fmt.Printf("Received message: %s\n", msg.Payload)
		}
	}()
	
	// Publish messages
	err := rdb.Publish(ctx, "notifications", "Hello subscribers!").Err()
	if err != nil {
		fmt.Printf("Failed to publish message: %v\n", err)
		return
	}
	
	// Wait a bit for message to be received
	time.Sleep(time.Second)
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	
	ctx := context.Background()
	
	fmt.Println("=== Caching Example ===")
	cacheExample(ctx, rdb)
	
	fmt.Println("\n=== Rate Limiting Example ===")
	rateLimitExample(ctx, rdb)
	
	fmt.Println("\n=== Session Management Example ===")
	sessionExample(ctx, rdb)
	
	fmt.Println("\n=== Pub/Sub Example ===")
	pubSubExample(ctx, rdb)
}
