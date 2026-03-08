package client

import (
	"testing"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/internal/cache"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient(
		WithTimeout(20*time.Second),
		WithRateLimit(30, time.Minute),
	)

	// Use Testify's assert package for better readability
	assert.NoError(t, err, "Expected no error creating client")
	assert.NotNil(t, c, "Expected non-nil client")
	assert.Equal(t, baseURL, c.baseURL, "Expected baseURL to be %s, got %s", baseURL, c.baseURL)
}

func TestNewClient_WithRedisCache(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	c, err := NewClient(
		WithRedisCache(cache.RedisOptions{Addr: mr.Addr(), KeyPrefix: "test"}),
	)
	require.NoError(t, err)
	assert.NotNil(t, c)
}

func TestNewClient_WithRedisCache_UnreachableServer(t *testing.T) {
	_, err := NewClient(
		WithRedisCache(cache.RedisOptions{Addr: "localhost:19999"}),
	)
	assert.Error(t, err, "should fail when Redis server is unreachable")
}
