package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	const (
		day1Key  = "analytics:active_users:day1"
		day2Key  = "analytics:active_users:day2"
		unionKey = "analytics:active_users:total_active"
	)

	// User IDs correspond directly to the bit index offset
	userA := int64(1024)
	userB := int64(500000) // 500k offset
	userC := int64(750000)

	fmt.Println("------1. Setting and getting individual bits (O(1))----")

	// Mark users active on Day1
	rdb.SetBit(ctx, day1Key, userA, 1)
	rdb.SetBit(ctx, day1Key, userB, 1)
	// Mark users active on Day2

	rdb.SetBit(ctx, day2Key, userB, 1)
	rdb.SetBit(ctx, day2Key, userC, 1)

	// check if userA was active on Day2

	isActiveDay2, err := rdb.GetBit(ctx, day2Key, userA).Result()
	if err != nil {
		log.Fatalf("GetBit failed: %v", err)
	}
	fmt.Printf("Was User ID %d active on Day 2?: %t\n", userA, isActiveDay2 == 1)

	fmt.Println("\n ------2. Counting Set Bits (BITCOUNT) ---")

	// Count how many total users were active  on Day1
	day1Count, _ := rdb.BitCount(ctx, day1Key, nil).Result()
	fmt.Printf("Total unique users active on Day 1: %d \n", day1Count)

	fmt.Println("\n--- 3. Server-side Bit Operations across Keys (BITOP) ---")
	// Calculate the union of Day1 and Day2 (Total unique active users across both days)

	// This happens entirely in Redis memory without dragging bits over the network wire.
	_, err = rdb.BitOpOr(ctx, unionKey, day1Key, day2Key).Result()
	if err != nil {
		log.Fatalf("BitOp failed: %v", err)
	}

	totalUniqueActive, _ := rdb.BitCount(ctx, unionKey, nil).Result()
	fmt.Printf("Total unique users across Day1 AND Day2 %d\n", totalUniqueActive)

	// Clean up keys
	rdb.Del(ctx, day1Key, day2Key, unionKey)
}
