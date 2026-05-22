package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Initialize client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	// Unique key for our user resource
	const userKey = "user:profile:99"

	fmt.Println("--- Creating a Hash (HSET) --- ")

	// storing structural field natively using map[string]interface{}
	// Command: HSET user:profile:99 username "johndoe" email "john.example.com"
	// role engineer

	err := rdb.HSet(ctx, userKey, map[string]interface{}{
		"username": "johndoe",
		"email":    "john@example.com",
		"role":     "engineer",
		"status":   "active",
	}).Err()

	if err != nil {
		log.Fatal("HSET failed: %v", err)
	}
	fmt.Println("User profile hash successfully created in Redis.")

	fmt.Println("\n--- 2. Fetching Specific Fields (HGET & HEXISTS) ---")
	// Verify if a field exists before querying it

	exists, _ := rdb.HExists(context.Background(), userKey, "email").Result()
	if exists {
		// Fetch only email.field

		email, err := rdb.HGet(ctx, userKey, "email").Result()
		if err != nil {
			log.Fatalf("HGET failed: %v", err)
		}
		fmt.Println("User email successfully fetched in Redis. %s \n ", email)
	}

	fmt.Println("\n--- 3. Modifying Specific Fields (HSET / HINCRBY) ---")
	// Update just the status field without touching the rest of the object
	rdb.HSet(ctx, userKey, "status", "suspended")

	// Hashes cam also contain numeric fields that support atomic incrementing

	rdb.HSet(ctx, userKey, "login_count", 1)
	newCount, _ := rdb.HIncrBy(ctx, userKey, "login_count", 1).Result()
	fmt.Printf("Updated Field 'status' to suspended. 'login_count' atomically incremented to: %d\n", newCount)

	fmt.Println("\n -- 4. Fetching complete Object (HGETALL) ----")
	// HGetAll fetches all fields and values, which go-redis automatically map
	// into a go map[string]string

	profileMap, err := rdb.HGetAll(ctx, userKey).Result()
	if err != nil {
		log.Fatalf("HGETALL failed: %v", err)
	}

	fmt.Println("Iterating through retrieved HGetAll Map: ")
	for field, value := range profileMap {
		fmt.Printf("  -> %s \n", field, value)
	}

	fmt.Println("5. Removing Specific fields (HDEL) ---")
	// Delete just the temporary login count tracking field
	deleteRows, _ := rdb.HDel(ctx, userKey, "login_count").Result()
	fmt.Printf("Remove fields count via HDEL: %d\n", deleteRows)

	// clean up hash
	rdb.Del(ctx, userKey)
}
