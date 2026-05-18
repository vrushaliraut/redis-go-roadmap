package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {

	// Initialize context with strict 3 seconds timeout for connection phase

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// configure redis client option

	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",              // No password set in compose command
		DB:           0,               // Default logical database
		DialTimeout:  5 * time.Second, // Connection establishing timeout
		ReadTimeout:  5 * time.Second, // Socket read timeout
		WriteTimeout: 5 * time.Second, // Socket write timeout
	})

	// Ensure client connections pool is closed when the application terminates

	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("failed to close redis client: %v", err)
		}
	}()

	// Execute PING - verify connection
	pong, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("failed to ping redis: %v", err)
	}

	// If successful, Redis returns a string payload "PONG"
	fmt.Printf("Connection Successful! Redis responded with : %s\n", pong)
}
