package main

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeLongURL(t *testing.T) {
	longURL := "http://localhost:8080/long-url-to-be-shortened"

	expectedEncodedURL := "EjpUVwkp4aqEyls6ZbbOc"

	result := encodeLongURL(longURL)
	assert.Equal(t, expectedEncodedURL, result)
}

func TestGetURLFromHash(t *testing.T) {
	hash := "g5gFCB6"
	expectedURL := "http://localhost:8080/g5gFCB6"

	result := getURLFromHash(hash)
	assert.Equal(t, expectedURL, result)
}

func TestShortenLongURL(t *testing.T) {
	db, mock := redismock.NewClientMock()
	longURL := "https://www.example.com/long-url-to-be-shortened"
	expectedEncodedURL := "cHlpOX8M0Oxece4FG6nVbZ"

	mock.ExpectGet(longURL).RedisNil()
	mock.Regexp().ExpectSet(longURL, expectedEncodedURL, 0).SetVal("OK")
	ctx := context.Background()
	result, err := shortenLongURL(ctx, longURL, db)

	assert.NoError(t, err)
	assert.Equal(t, getURLFromHash(expectedEncodedURL), result)
}
