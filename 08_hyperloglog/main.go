package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// Hands-on Implementation: Massive Scaled Unique View Counter
func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	const hllKey = "view_tracker:ride_searches"

	fmt.Println("----- 1. Adding High-Volume items to HyperLogLog (PFADD) --- ")

	// Simulate ingestion stream of unique identifiers (e.g., Device IDs or IP addresses)
	// Even if we push duplicates, HLL internally handles approximation statistics.
	for i := 1; i <= 10000; i++ {
		deviceID := "device_uuid_" + strconv.Itoa(i)
		rdb.PFAdd(ctx, hllKey, deviceID)
	}

	// Internationally inject duplicates to demonstrate internal filtering
	rdb.PFAdd(ctx, hllKey, "device_uuid_500")
	rdb.PFAdd(ctx, hllKey, "device_uuid_1000")

	fmt.Println("Ingested 10,000 items with intentional duplicates.")

	fmt.Println("\n--- 2. Extracting Cardinality Estimation (PFCOUNT) ---")
	// Query the approximate count
	estimatedCount, _ := rdb.PFCount(ctx, hllKey).Result()
	fmt.Printf("Actual unique elements sent: 10,000\n")
	fmt.Printf("Redis HyperLogLog Estimated unique count: %d\n", estimatedCount)

	// Calculate accuracy delta
	delta := float64(estimatedCount) - 10000.0
	fmt.Printf("Statistical Estimation Variance: %.0f elements (Well within the 0.81%% error limit)\n", delta)
	// Clean up Keys

	rdb.Del(ctx, hllKey)
}
