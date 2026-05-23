package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

// This Lua script check if stock exists, and if so, decrementing atomically
// It returns 1 for success, 0 for out-of-stock

const decrementScript = `
	local stock = tonumber(redis.call("GET", KEYS[1]))
	if stock and stock > 0 then
		redis.call("DECR", KEYS[1])
		return 1
	else 
		return 0
	end
	`

func main() {

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	promoKey := "promo:uber_50_off:stock"

	// Set initial stock to exactly 5
	rdb.Set(ctx, promoKey, 5, 0)
	fmt.Println("Promo stock initialized to 5.")

	// Simulate 100 concurrent users trying to claim the 5 promos at the exact same time

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// Execute the Lua script atomically
			result, err := rdb.Eval(ctx, decrementScript, []string{promoKey}).Result()

			if err != nil {
				log.Printf("User %d error: %v", userID, err)
				return
			}

			if result.(int64) == 1 {
				mu.Lock()
				successCount++
				mu.Unlock()
				fmt.Printf("User %d successfully claimed the promo! \n", userID)
			}
		}(i)
	}
	wg.Wait()

	finalStock, _ := rdb.Get(ctx, promoKey).Result()
	fmt.Printf("\n Total successful claims: %d\n", successCount)
	fmt.Printf("Final stock remaining in Redis: %s (Notice it never drops below 0!)\n", finalStock)
}
