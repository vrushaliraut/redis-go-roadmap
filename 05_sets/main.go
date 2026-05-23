package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {

	ctx := context.Background()

	// Initialize redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	// keys
	const (
		user1Interests = "user:101:interests"
		user2Interests = "user:102:interests"
	)

	fmt.Println("---- 1. Unique Tracking (SADD & SISMEMBER) --- ")

	// SADD returns the number of *new* elements successfully added

	added, _ := rdb.SAdd(ctx, user1Interests, "golang", "redis", "docker", "golang", "redis").Result()

	// ignoring duplicate
	fmt.Printf("Added user 101. New unique items recorded: %d \n", added)

	// check if specific elements exists into the set
	isMember, _ := rdb.SIsMember(ctx, user1Interests, "redis").Result()
	fmt.Printf("Does User 101 like 'redis' ? : %t\n", isMember)

	fmt.Println("---- 2. Set Intersections (SINTER) ---")

	// Add interest for User 102
	rdb.SAdd(ctx, user2Interests, "rust", "docker", "golang", "kubernates")

	// Find cmmon interest between user1 and user2

	common, _ := rdb.SInter(ctx, user1Interests, user2Interests).Result()
	fmt.Println("Matching mutual interest found via SINTER:")
	for _, interest := range common {
		fmt.Printf(" .%s\n", interest)
	}

	// clean up task

	rdb.Del(ctx, user1Interests, user1Interests)
}
