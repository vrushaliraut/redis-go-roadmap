package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Go Implementation: Simulating a Pub/Sub Channel
func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	channel := "dispatch:driver_zone_A"

	//1. Setup the subscriber(ctx, channel)
	pubsub := rdb.Subscribe(ctx, channel)
	defer pubsub.Close()

	// Wait for subscription confirmation
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	// Go routine to listen for messages continuously
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Printf("Driver App Received Dispatch: %s \n", msg.Payload)
		}
	}()

	// Setup the publisher
	fmt.Println("Dispatcher sending ride requests")
	rdb.Publish(ctx, channel, "Ride Request: Airport pickup (John)")
	time.Sleep(1 * time.Second)

	rdb.Publish(ctx, channel, "ride Request: Downtown drop-off (Sarah)")
	time.Sleep(1 * time.Second)
	fmt.Println("Finished publishing")
}
