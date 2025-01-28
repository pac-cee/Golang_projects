package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	fmt.Println("=== Lists (Queue/Stack) Example ===")
	// Lists - can be used as queues or stacks
	{
		// Push items to list (queue)
		err := rdb.LPush(ctx, "queue", "first", "second", "third").Err()
		if err != nil {
			fmt.Printf("Failed to push to list: %v\n", err)
			return
		}

		// Pop item from queue (FIFO - First In First Out)
		val, err := rdb.RPop(ctx, "queue").Result()
		if err != nil {
			fmt.Printf("Failed to pop from list: %v\n", err)
			return
		}
		fmt.Printf("Popped from queue: %v\n", val)

		// Get all items in list
		items, err := rdb.LRange(ctx, "queue", 0, -1).Result()
		if err != nil {
			fmt.Printf("Failed to get list range: %v\n", err)
			return
		}
		fmt.Printf("Remaining items in queue: %v\n", items)
	}

	fmt.Println("\n=== Sets Example ===")
	// Sets - collection of unique elements
	{
		// Add elements to set
		err := rdb.SAdd(ctx, "fruits", "apple", "banana", "orange").Err()
		if err != nil {
			fmt.Printf("Failed to add to set: %v\n", err)
			return
		}

		// Check if element exists in set
		exists, err := rdb.SIsMember(ctx, "fruits", "apple").Result()
		if err != nil {
			fmt.Printf("Failed to check set membership: %v\n", err)
			return
		}
		fmt.Printf("Is apple in fruits set? %v\n", exists)

		// Get all set members
		members, err := rdb.SMembers(ctx, "fruits").Result()
		if err != nil {
			fmt.Printf("Failed to get set members: %v\n", err)
			return
		}
		fmt.Printf("All fruits: %v\n", members)
	}

	fmt.Println("\n=== Hashes Example ===")
	// Hashes - maps between string fields and string values
	{
		// Set multiple hash fields
		err := rdb.HSet(ctx, "user:1", map[string]interface{}{
			"username": "johndoe",
			"email":    "john@example.com",
			"age":      "30",
		}).Err()
		if err != nil {
			fmt.Printf("Failed to set hash: %v\n", err)
			return
		}

		// Get specific hash field
		email, err := rdb.HGet(ctx, "user:1", "email").Result()
		if err != nil {
			fmt.Printf("Failed to get hash field: %v\n", err)
			return
		}
		fmt.Printf("User email: %v\n", email)

		// Get all hash fields
		userData, err := rdb.HGetAll(ctx, "user:1").Result()
		if err != nil {
			fmt.Printf("Failed to get all hash fields: %v\n", err)
			return
		}
		fmt.Printf("All user data: %v\n", userData)
	}

	fmt.Println("\n=== Sorted Sets Example ===")
	// Sorted Sets - ordered sets of unique elements
	{
		// Add members with scores
		users := []redis.Z{
			{Score: 100, Member: "user1"},
			{Score: 200, Member: "user2"},
			{Score: 150, Member: "user3"},
		}
		err := rdb.ZAdd(ctx, "scores", users...).Err()
		if err != nil {
			fmt.Printf("Failed to add to sorted set: %v\n", err)
			return
		}

		// Get top scores (highest to lowest)
		topScores, err := rdb.ZRevRangeWithScores(ctx, "scores", 0, -1).Result()
		if err != nil {
			fmt.Printf("Failed to get sorted set range: %v\n", err)
			return
		}
		fmt.Println("User scores (highest to lowest):")
		for _, z := range topScores {
			fmt.Printf("%v: %v\n", z.Member, z.Score)
		}

		// Get user rank
		rank, err := rdb.ZRank(ctx, "scores", "user2").Result()
		if err != nil {
			fmt.Printf("Failed to get rank: %v\n", err)
			return
		}
		fmt.Printf("user2 rank: %v\n", rank)
	}
}
