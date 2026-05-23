package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// Hands-on Implementation: Driver Proximity Search
func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	defer rdb.Close()

	const driversLocationKey = "drivers:locations"

	fmt.Println("---- 1. Indexing Geographical coordinates (GEOADD) ----")

	// Parameters: Context, Key, Longitude, Latitude, Member Name
	// Note - Redis uses Longitude first than Latitude

	err := rdb.GeoAdd(ctx, driversLocationKey,
		&redis.GeoLocation{Longitude: 106.827153, Latitude: -6.175110, Name: "driver_jakarta_center"},
		&redis.GeoLocation{Longitude: 106.811234, Latitude: -6.194567, Name: "driver_jakarta_south"},
		&redis.GeoLocation{Longitude: 107.619123, Latitude: -6.917464, Name: "driver_bandung_distant"},
	).Err()

	if err != nil {
		log.Fatalf("GeoAdd failed: %v", err)
	}

	fmt.Println("Successfully indexed live coordinate telemetry vectors.")

	fmt.Println("\n--- 2. Radius Search for Nearby Entities (GEOSEARCH) ---")
	// Search for drivers within a 5 Kilometer radius of a customer's location
	customerLong := 106.822910
	customerLat := -6.183200

	results, err := rdb.GeoSearchLocation(ctx, driversLocationKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  customerLong,
			Latitude:   customerLat,
			Radius:     5,
			RadiusUnit: "km",
			Sort:       "ASC", // Nearest driver first
			Count:      10,    // specify limit
		},
		WithDist:  true, // This returns the distance
		WithCoord: true, // Optional: returns coordinates
	}).Result()
	if err != nil {
		log.Fatalf("GeoSearch failed: %v", err)
	}

	fmt.Printf("Found %d drivers within a 5KM operational threshold:\n", len(results))
	for _, location := range results {
		fmt.Printf("  -> Name: %s | Distance: %.2f KM\n", location.Name, location.Dist)
	}

	rdb.Del(ctx, driversLocationKey)
}
