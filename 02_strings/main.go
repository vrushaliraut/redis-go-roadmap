package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Initialize client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	// Define our keys
	const (
		sessionKey = "user:session:101"
		rateKey    = "user:hitcount:101"
	)

	fmt.Println("---- 1. Basic SET and GET with TTL ----")
	// Store a session token that automatically expires in 10 seconds
	// Command : SET user:session:101 "active_token_xyz" EX 10

	err := rdb.Set(ctx, sessionKey, "active_token_xyz", 10*time.Second).Err()

	if err != nil {
		log.Fatalf("SET failed: %v", err)
	}

	fmt.Println("Saved session with 10sec expiration")

	// Retrieve the value
	token, err := rdb.Get(ctx, sessionKey).Result()
	if err != nil {
		log.Fatalf("GET failed: %v", err)
	}
	fmt.Printf("Got session key %s\n", token)

	fmt.Println("---- 2. Inspect Metadata (STRLEN and TTL) ----")
	// Get string length
	strLen, _ := rdb.StrLen(ctx, sessionKey).Result()

	// Get remaining time to live
	ttl, _ := rdb.TTL(ctx, sessionKey).Result()
	fmt.Printf("Token Length: %d characters \n", strLen)
	fmt.Printf("Got TTL %v \n\n", ttl)

	fmt.Println("-----3.Atomic Counters (INCR & APPEND) ---- ")
	// INCR automatically initializes the key to 0 if it doesn't exist,
	// then increment it to 1

	for i := 1; i <= 3; i++ {
		currentHits, err := rdb.Incr(ctx, rateKey).Result()
		if err != nil {
			log.Fatalf("INCR failed: %v", err)
		}
		fmt.Printf("API Request #%d processed. "+
			"Current Total Hits stored in Redis: %d\n", i, currentHits)
	}

	rdb.Append(ctx, "user:notes:101", "Premium_").Err()
	rdb.Append(ctx, "user:notes:101", "Tier").Err()

	notes, _ := rdb.Get(ctx, "user:notes:101").Result()
	fmt.Printf("Appended Note Result: %s\n\n", notes)

	fmt.Println("--- 4. Data Cleanup & Expiration (DEL) ---")
	// Simulate checking if data expires
	fmt.Println("Waiting 3 seconds to observe TTL countdown...")
	time.Sleep(3 * time.Second)
	newTTL, _ := rdb.TTL(ctx, sessionKey).Result()
	fmt.Printf("New TTL %v \n\n", newTTL)

	// Explicitly deleting keys manually
	fmt.Println("Cleaning up hit counter and notes keys...")
	rdb.Del(ctx, rateKey, "user:notes:101")

	// Verify key deletion
	_, err = rdb.Get(ctx, rateKey).Result()
	if errors.Is(err, redis.Nil) {
		fmt.Println("Verification: 'user:hitcount:101' successfully deleted (Key does not exist).")
	}
}
