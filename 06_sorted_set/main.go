package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Initialise client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	const leaderboardKey = "game:leaderboard"

	fmt.Println("-----1. Populating the Leaderboard (ZADD)----")

	// Add players with their baseline game scores
	rdb.ZAdd(ctx, leaderboardKey, redis.Z{Score: 4500, Member: "player_alice"})
	rdb.ZAdd(ctx, leaderboardKey, redis.Z{Score: 7200, Member: "player_bob"})
	rdb.ZAdd(ctx, leaderboardKey, redis.Z{Score: 3100, Member: "player_charlie"})

	fmt.Println("Initial scores registered.")

	fmt.Println("------Live score update  (ZINCRBY)-----")
	// Alice complete an objective and scores 3000 more points

	newScore, _ := rdb.ZIncrBy(ctx, leaderboardKey, 3000, "player_alice").Result()
	fmt.Println("player_alice score updated! New score: %.0f \n", newScore)

	fmt.Println("\n ------3. Fetching the Top Ranks (ZRANGE) --- ")

	// Fetch the top 3 players in descending order (highest score first)
	// We use ZRangeArgs with Rev: true to sort backwards

	topPlayers, _ := rdb.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{
		Key:   leaderboardKey,
		Start: 0,
		Stop:  2,
		Rev:   true,
	}).Result()

	fmt.Println("Current leaderboard top 3")

	for rank, p := range topPlayers {
		fmt.Printf("Rank %d | Name: %s | Score: %.0f\n", rank+1, p.Member, p.Score)
	}

	// Clean up
	rdb.Del(ctx, leaderboardKey)
}
