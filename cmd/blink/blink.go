package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"math/big"
)

const base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// shortenLongURL implements the logic for handling the shortening feature of the application.
func shortenLongURL(ctx context.Context, longURL string, redisClient *redis.Client) (string, error) {
	// Check if the long URL exists in the Redis cache.
	// If yes, return the corresponding shortened URL for it.
	exists, value, err := checkURLInCache(ctx, redisClient, longURL)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	if exists {
		return getURLFromHash(value), nil
	}

	// Since we have reached here, it means the entry for this URL doesn't
	// exist in the cache. In that case, let's create a shortended URL, set
	// it in the cache and then return it.
	encodedURL := encodeLongURL(longURL)
	err = setRedisKeyValue(ctx, redisClient, longURL, encodedURL)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}
	return getURLFromHash(encodedURL), nil
}

// encodeLongURL encodes the long URL sent as a parameter into a shorter URL
// using md5 and base62.
func encodeLongURL(longURL string) string {
	// Calculate MD5 hash of the long URL.
	hash := md5.Sum([]byte(longURL))

	// Convert the hash to a big integer.
	bigIntHash := new(big.Int)
	bigIntHash.SetBytes(hash[:])

	// Initialize variables for base62 conversion.
	var base62Chars []byte
	base := big.NewInt(62)

	// Convert the hash to base62.
	for bigIntHash.Cmp(big.NewInt(0)) > 0 {
		// Calculate the remainder from dividing by base (62).
		remainder := new(big.Int)
		bigIntHash.DivMod(bigIntHash, base, remainder)

		// Append the corresponding base62 character.
		base62Chars = append(base62Chars, base62Charset[remainder.Int64()])
	}

	// Reverse the base62 characters to get the correct order.
	reverse(base62Chars)
	return string(base62Chars)
}

func reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// getURLFromHash creates a formatted URL from the hash and returns it.
func getURLFromHash(hash string) string {
	return fmt.Sprintf("%s://%s/%s", serverScheme, hostName, hash)
}
