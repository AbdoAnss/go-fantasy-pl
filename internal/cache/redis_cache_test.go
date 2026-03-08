package cache_test

import (
	"testing"
	"time"

	"github.com/AbdoAnss/go-fantasy-pl/internal/cache"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestRedisCache spins up an in-process miniredis server and returns a
// RedisCache connected to it, along with a cleanup function.
func newTestRedisCache(t *testing.T, keyPrefix string) (*cache.RedisCache, *miniredis.Miniredis) {
	t.Helper()
	mr, err := miniredis.Run()
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return cache.NewRedisCacheWithClient(rdb, keyPrefix), mr
}

// TestRedisCache_Contract validates that RedisCache satisfies the same
// behavioral contract as MemoryCache by reusing the shared test suite.
func TestRedisCache_Contract(t *testing.T) {
	// Use a key prefix so that Clear() is not a no-op
	c, mr := newTestRedisCache(t, "test")
	defer mr.Close()

	testCacheImplementation(t, c)
}

func TestRedisCache_KeyPrefix(t *testing.T) {
	c, mr := newTestRedisCache(t, "fpl")
	defer mr.Close()

	err := c.Set("teams", "value", time.Minute)
	require.NoError(t, err)

	// Verify the actual Redis key includes the prefix
	assert.True(t, mr.Exists("fpl:teams"), "key should be stored with prefix")
	assert.False(t, mr.Exists("teams"), "unprefixed key should not exist")

	var got string
	ok := c.Get("teams", &got)
	assert.True(t, ok)
	assert.Equal(t, "value", got)
}

func TestRedisCache_TTLEnforcement(t *testing.T) {
	c, mr := newTestRedisCache(t, "")
	defer mr.Close()

	err := c.Set("ttl_key", "value", 100*time.Millisecond)
	require.NoError(t, err)

	// Fast-forward miniredis time so the key expires
	mr.FastForward(200 * time.Millisecond)

	var got string
	ok := c.Get("ttl_key", &got)
	assert.False(t, ok, "key should have expired")
}

func TestRedisCache_Clear_WithPrefix(t *testing.T) {
	c, mr := newTestRedisCache(t, "ns")
	defer mr.Close()

	require.NoError(t, c.Set("a", "1", time.Minute))
	require.NoError(t, c.Set("b", "2", time.Minute))

	c.Clear()

	var v string
	assert.False(t, c.Get("a", &v))
	assert.False(t, c.Get("b", &v))
}

func TestRedisCache_Clear_NoPrefix_IsNoOp(t *testing.T) {
	c, mr := newTestRedisCache(t, "")
	defer mr.Close()

	// Write a key directly into miniredis (simulating another application's data)
	mr.Set("other_app_key", "sensitive")

	// Clear without a prefix should NOT delete any keys
	c.Clear()

	assert.True(t, mr.Exists("other_app_key"),
		"Clear with no prefix must not delete keys from other applications")
}

func TestRedisCache_NewRedisCache_ConnectionError(t *testing.T) {
	_, err := cache.NewRedisCache(cache.RedisOptions{
		Addr: "localhost:19999", // nothing listening here
	})
	assert.Error(t, err, "should fail when Redis is not reachable")
}

func TestRedisCache_NewRedisCache_Defaults(t *testing.T) {
	// We can't connect to a real Redis, but we can confirm the constructor
	// returns an error (rather than panicking) when no server is reachable.
	_, err := cache.NewRedisCache(cache.RedisOptions{})
	assert.Error(t, err, "should fail when default Redis address is not reachable")
}

func TestRedisCache_Close(t *testing.T) {
	c, mr := newTestRedisCache(t, "")
	defer mr.Close()

	err := c.Close()
	assert.NoError(t, err)
}
