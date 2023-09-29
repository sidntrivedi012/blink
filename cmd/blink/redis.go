package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/url"
	"sort"
	"strings"
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

// fetchMostShortenedURLs fetches the top three most shortened domains
// from the Redis database.
func fetchMostShortenedURLs(ctx context.Context, client *redis.Client) (string, error) {
	longURLFrequencyMap := make(map[string]int)
	// Get all keys in the Redis cache.
	longURLs, err := client.Keys(ctx, "*").Result()
	if err != nil {
		return "", err
	}

	// Iterate over keys and check if the value exists
	for _, longURL := range longURLs {
		parsedURL, err := url.Parse(longURL)
		if err != nil {
			return "", err
		}

		urlHostName := strings.TrimPrefix(parsedURL.Hostname(), "www.")
		slog.Info("log", slog.String(longURL, urlHostName))
		longURLFrequencyMap[urlHostName] = longURLFrequencyMap[urlHostName] + 1
	}

	// Now, as we have the map of the hostnames alongwith the frequency
	// of their shortening, we can sort the map in descending order of values
	// and then pick the top three.
	type hostCount struct {
		name  string
		count int
	}

	var hostCountSlice []hostCount
	for k, v := range longURLFrequencyMap {
		hostCountSlice = append(hostCountSlice, hostCount{k, v})
	}
	sort.Slice(hostCountSlice, func(i, j int) bool {
		return hostCountSlice[i].count > hostCountSlice[j].count
	})

	// Picking top three from the slice sorted in descending order.
	topThreeHosts := ""
	for i := 0; i <= 2; i++ {
		if i == len(hostCountSlice) {
			break
		}
		topThreeHosts = topThreeHosts + fmt.Sprintf("%s\t%d\n", hostCountSlice[i].name, hostCountSlice[i].count)
	}
	return topThreeHosts, nil
}
