package main

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckURLInCache(t *testing.T) {
	db, mock := redismock.NewClientMock()
	url := "https://google.com/helloworld"

	mock.ExpectGet(url).RedisNil()
	ctx := context.Background()
	exists, value, err := checkURLInCache(ctx, db, url)

	assert.NoError(t, err)
	assert.Equal(t, false, exists)
	assert.Equal(t, value, "")
}

func TestSetRedisKeyValue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	dummyKey := "https://google.com/helloworld"
	dummyValue := "anyrandomhash"

	mock.ExpectSet(dummyKey, dummyValue, 0).SetVal("OK")
	ctx := context.Background()
	err := setRedisKeyValue(ctx, db, dummyKey, dummyValue)
	assert.NoError(t, err)
}

func TestGetLongURLFromCache(t *testing.T) {
	ctx := context.Background()
	db, mock := redismock.NewClientMock()
	expectedKeys := []string{
		"https://github.com/sidntrivedi012",
		"https://github.com/hello",
	}
	mock.ExpectKeys("*").SetVal(expectedKeys)
	mock.ExpectGet(expectedKeys[0]).SetVal("wronghash")
	mock.ExpectGet(expectedKeys[1]).SetVal("anotherrandomhash")

	exists, longURL, err := getLongURLFromCache(ctx, db, "anotherrandomhash")

	assert.NoError(t, err)
	assert.Equal(t, exists, true)
	assert.Equal(t, longURL, expectedKeys[1])
}

func TestFetchMostShortenedURLs(t *testing.T) {
	db, mock := redismock.NewClientMock()

	expectedKeys := []string{
		"https://github.com/sidntrivedi012",
		"https://github.com/hello",
		"https://github.com/infracloud",
		"https://google.com/helloworld",
		"https://google.com/binpack",
		"https://infracloud.io/hello",
		"https://twitter.com/sidntrivedi012",
		"https://linkedin.com/in/siddhantntrivedi",
	}
	mock.ExpectKeys("*").SetVal(expectedKeys)

	ctx := context.Background()
	result, err := fetchMostShortenedURLs(ctx, db)

	assert.NoError(t, err)
	assert.Contains(t, result, "github.com\t3")
	assert.Contains(t, result, "google.com\t2")
}
