package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

// checkURLInCache checks if a key similar to the long URL that is being
// processed exists already in the redis cache.
func checkURLInCache(ctx context.Context, client *redis.Client, longURL string) (bool, string, error) {
	// Check if the long URL exists in the Redis cache
	val, err := client.Get(ctx, longURL).Result()
	if err == redis.Nil {
		// Key doesn't exist, the URL is not in the cache.
		return false, "", nil
	} else if err != nil {
		return false, "", err
	}

	// Key exists, the URL is in the cache.
	slog.Info("already existing url in cache", slog.String("long url", longURL), slog.String("hash", val))
	return true, val, nil
}

// setRedisKeyValue stores a key-value pair in the redis cache.
func setRedisKeyValue(ctx context.Context, redisClient *redis.Client, key string, value string) error {
	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting value for key %s: %v", key, err)
	}
	return nil
}

// getKeyValueFromCache fetches the value of a key from redis cache.
func getKeyValueFromCache(ctx context.Context, redisClient *redis.Client, key string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found: %s", key)
	} else if err != nil {
		return "", fmt.Errorf("error retrieving value for key %s: %v", key, err)
	}
	return val, nil
}
