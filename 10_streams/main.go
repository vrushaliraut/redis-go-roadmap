package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Hands-on Implementation: Transaction/Event Logging Stream

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	const orderStreamKey = "stream:orders"

	fmt.Println("--- 1. Appending Telemetry Logs to Stream (XADD) ---")
	// XAdd generates a unique, auto-incremented ID containing <timestamp>-<sequence>
	// We pass a map of key-value attributes for the event payload

	id1, _ := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: orderStreamKey,
		ID:     "*", // Auto-generated ID
		Values: map[string]interface{}{"order_id": "ORD-552", "amount": 25.50, "status": "created"},
	}).Result()

	fmt.Printf("Appended Event 1 to stream. Assigned ID: %s\n", id1)

	id2, _ := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: orderStreamKey,
		ID:     "*", // Auto-generated ID
		Values: map[string]interface{}{"order_id": "ORD-553", "amount": 25.50, "status": "created"},
	}).Result()

	fmt.Printf("Appended Event 2 to stream. Assigned ID: %s\n", id2)

	fmt.Println("\n--- 2. Querying Ranges of Stream Logs (XRANGE) ---")
	// Fetch stream events from start ("-") to end ("+")
	events, err := rdb.XRange(ctx, orderStreamKey, "-", "+").Result()
	if err != nil {
		log.Fatalf("XRNGE Failed: %v", err)
	}

	fmt.Printf("Retrieved %d events historical stream entries:\n", len(events))
	for _, entry := range events {
		fmt.Printf(" Event Id: %s | OrderID: %s | Status: %s \n",
			entry.ID, entry.Values["order_id"], entry.Values["status"])
	}

	rdb.Del(ctx, orderStreamKey)
}
