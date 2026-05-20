package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AbdoAnss/go-fantasy-pl/client"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	fmt.Printf("Connecting to Redis at %s...\n", redisAddr)

	c, err := client.NewClient(
		client.WithRedisCache(client.RedisOptions{
			Addr:      redisAddr,
			KeyPrefix: "fpl_example",
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Fetching teams (this will populate the cache)...")
	teams, err := c.Teams.GetAllTeams()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched %d teams.\n", len(teams))

	fmt.Println("Fetching teams again (this should come from Redis cache)...")
	teamsCached, err := c.Teams.GetAllTeams()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched %d teams from cache.\n", len(teamsCached))

	fmt.Println("Success! Redis caching is working.")
}
