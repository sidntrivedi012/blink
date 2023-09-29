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

// getURLFromHash fetches the long URL for a hash stored in Redis cache.
// Helps in the redirection workflow.
func getLongURLFromCache(ctx context.Context, client *redis.Client, value string) (bool, string, error) {
	// Get all keys in the Redis cache.
	keys, err := client.Keys(ctx, "*").Result()
	if err != nil {
		return false, "", err
	}

	// Iterate over keys and check if the value exists
	for _, key := range keys {
		val, err := client.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return false, "", err
		}

		if val == value {
			return true, key, nil
		}
	}
	return false, "", nil
}
