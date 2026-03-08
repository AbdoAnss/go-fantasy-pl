package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache is a Cache implementation backed by Redis.
// It is safe for concurrent use.
type RedisCache struct {
	client *redis.Client
	prefix string
}

// RedisOptions configures a RedisCache.
type RedisOptions struct {
	// Addr is the Redis server address in "host:port" form (default "localhost:6379").
	Addr string
	// Password is the Redis AUTH password. Leave empty for no authentication.
	Password string
	// DB is the Redis database index to use (default 0).
	DB int
	// KeyPrefix is an optional prefix applied to every key to avoid collisions
	// when multiple applications share the same Redis instance.
	KeyPrefix string
}

// NewRedisCache creates a RedisCache that connects to Redis using the supplied options.
// It performs a PING to verify connectivity and returns an error on failure.
func NewRedisCache(opts RedisOptions) (*RedisCache, error) {
	if opts.Addr == "" {
		opts.Addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: failed to connect to %s: %w", opts.Addr, err)
	}

	return &RedisCache{client: rdb, prefix: opts.KeyPrefix}, nil
}

// NewRedisCacheWithClient creates a RedisCache using an already-constructed *redis.Client.
// Useful for testing (e.g. with miniredis) or when finer-grained client configuration is required.
func NewRedisCacheWithClient(client *redis.Client, keyPrefix string) *RedisCache {
	return &RedisCache{client: client, prefix: keyPrefix}
}

func (r *RedisCache) prefixedKey(key string) string {
	if r.prefix == "" {
		return key
	}
	return r.prefix + ":" + key
}

// Set serializes value to JSON and stores it in Redis with the given TTL.
func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis cache: failed to marshal value for key %q: %w", key, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Set(ctx, r.prefixedKey(key), data, ttl).Err(); err != nil {
		return fmt.Errorf("redis cache: failed to set key %q: %w", key, err)
	}
	return nil
}

// Get retrieves the value for key from Redis and deserializes it into dest.
// Returns false if the key does not exist, has expired, or cannot be unmarshaled.
func (r *RedisCache) Get(key string, dest interface{}) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := r.client.Get(ctx, r.prefixedKey(key)).Bytes()
	if err != nil {
		return false
	}

	return json.Unmarshal(data, dest) == nil
}

// Delete removes the key from Redis.
func (r *RedisCache) Delete(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r.client.Del(ctx, r.prefixedKey(key))
}

// Clear removes all keys that match the cache prefix from Redis.
// If no prefix is set, it flushes the entire database (FLUSHDB).
// Use with caution in shared Redis instances — always configure a KeyPrefix.
func (r *RedisCache) Clear() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if r.prefix == "" {
		r.client.FlushDB(ctx)
		return
	}

	pattern := r.prefix + ":*"
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if len(keys) > 0 {
		r.client.Del(ctx, keys...)
	}
}

// Close closes the underlying Redis connection pool.
func (r *RedisCache) Close() error {
	return r.client.Close()
}
