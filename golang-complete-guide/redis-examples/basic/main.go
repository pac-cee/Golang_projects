package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	ctx := context.Background()

	// Ping the Redis server to check connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("Failed to connect to Redis: %v\n", err)
		return
	}
	fmt.Println("Connected to Redis!")

	// Basic SET operation
	err := rdb.Set(ctx, "mykey", "Hello from Redis!", 0).Err()
	if err != nil {
		fmt.Printf("Failed to set key: %v\n", err)
		return
	}

	// Basic GET operation
	val, err := rdb.Get(ctx, "mykey").Result()
	if err != nil {
		fmt.Printf("Failed to get key: %v\n", err)
		return
	}
	fmt.Printf("mykey = %v\n", val)

	// SET with expiration (SETEX)
	err = rdb.SetEx(ctx, "tempkey", "I will expire in 1 minute", 1*time.Minute).Err()
	if err != nil {
		fmt.Printf("Failed to set expiring key: %v\n", err)
		return
	}

	// Get TTL (Time To Live)
	ttl, err := rdb.TTL(ctx, "tempkey").Result()
	if err != nil {
		fmt.Printf("Failed to get TTL: %v\n", err)
		return
	}
	fmt.Printf("tempkey will expire in %v\n", ttl)

	// Check if key exists
	exists, err := rdb.Exists(ctx, "mykey").Result()
	if err != nil {
		fmt.Printf("Failed to check key existence: %v\n", err)
		return
	}
	fmt.Printf("Does mykey exist? %v\n", exists == 1)

	// Delete a key
	err = rdb.Del(ctx, "mykey").Err()
	if err != nil {
		fmt.Printf("Failed to delete key: %v\n", err)
		return
	}
	fmt.Println("Successfully deleted mykey")

	// Try to get deleted key
	val, err = rdb.Get(ctx, "mykey").Result()
	if err == redis.Nil {
		fmt.Println("mykey no longer exists")
	} else if err != nil {
		fmt.Printf("Failed to get key: %v\n", err)
		return
	}
}
