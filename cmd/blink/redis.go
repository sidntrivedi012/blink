package main

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

const redisPort int = 6379

// initRedisClient initializes a Redis client.
func initRedisClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%d", redisPort),
		Password: "",
		DB:       0,
	})
	return redisClient
}
