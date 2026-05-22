package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	// Initialize client

	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	defer rdb.Close()

	const queueKey = "tasks:email_queue"

	fmt.Println("---- 1. Producer: Enqueuing Items (LPUSH) ---")

	// Enqueue tasks into left side of the list

	tasks := []string{"send_welcome_email", "generate_invoice_pdf", "sync_crm_analytics"}

	for _, task := range tasks {
		// Command: LPUSH tasks: email_queue <task>
		length, err := rdb.LPush(ctx, queueKey, task).Result()
		if err != nil {
			log.Fatalf("LPUSH failed : %v", err)
		}
		fmt.Printf("Enqueued: '%s' | Current Queue Depth: %d \n", task, length)
	}

	fmt.Println("\n --- 2. Inspecting the Queue (LLEN & LRANGE) --- ")
	queueLen, _ := rdb.LLen(ctx, queueKey).Result()
	fmt.Printf("Total pending jobs: %d \n", queueLen)

	// Fetch items without consuming/removing them (0 is start index, -1 means last item)

	allTasks, _ := rdb.LRange(ctx, queueKey, 0, -1).Result()
	fmt.Println("Current raw array state in Redis (Notice:: Its reversed due to LPUSH):")

	for i, t := range allTasks {
		fmt.Printf(" Index [%d]: %s \n", i, t)
	}

	fmt.Println("\n ---- 3. Consumer: Processing Items (RPOP) ---")
	// Process items from the right side of the list (FIFO queue Mechanics)
	// We loop untill queue is empty

	for {
		// Command : RPOP tasks: email_queue

		task, err := rdb.LPop(ctx, queueKey).Result()
		if errors.Is(err, redis.Nil) {
			fmt.Println("Queue is now empty, Worker execution completed")
			break
		} else if err != nil {
			log.Fatalf("RPOP failed: %v", err)
		}

		fmt.Printf("Worker processed Job: %s\n", task)

	}

	// Clean up just in case
	rdb.Del(ctx, queueKey)
}
