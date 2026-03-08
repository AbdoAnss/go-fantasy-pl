package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Cache defines the interface for all cache implementations.
// Implementations must be safe for concurrent use.
type Cache interface {
	// Get retrieves a value from the cache and unmarshals it into dest.
	// Returns true if the key exists and has not expired, false otherwise.
	Get(key string, dest interface{}) bool
	// Set stores a value in the cache with the given TTL.
	Set(key string, value interface{}, ttl time.Duration) error
	// Delete removes a key from the cache.
	Delete(key string)
	// Clear removes all keys from the cache.
	Clear()
}

type item struct {
	value      []byte
	expiration time.Time
}

// MemoryCache is an in-memory Cache implementation backed by a map.
// It is safe for concurrent use.
type MemoryCache struct {
	items map[string]item
	mu    sync.RWMutex
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]item),
	}
}

// Set serializes value to JSON and stores it with the given TTL.
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: failed to marshal value for key %q: %w", key, err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:      data,
		expiration: time.Now().Add(ttl),
	}
	return nil
}

// Get deserializes the cached value into dest.
// Returns false if the key does not exist or has expired.
func (c *MemoryCache) Get(key string, dest interface{}) bool {
	c.mu.RLock()
	it, exists := c.items[key]
	c.mu.RUnlock()

	if !exists {
		return false
	}

	if time.Now().After(it.expiration) {
		go c.Delete(key)
		return false
	}

	return json.Unmarshal(it.value, dest) == nil
}

// Delete removes a key from the cache.
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear removes all items from the cache.
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]item)
}

// Cleanup removes all expired items from the cache.
func (c *MemoryCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, v := range c.items {
		if now.After(v.expiration) {
			delete(c.items, k)
		}
	}
}

// StartCleanupTask starts a goroutine that periodically removes expired items.
func (c *MemoryCache) StartCleanupTask(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.Cleanup()
		}
	}()
}
