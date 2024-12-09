package cache

import (
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration time.Time
}

type Cache struct {
	items map[string]item
	mu    sync.RWMutex
}

// NewCache creates a new cache
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]item),
	}
}

// Set adds an item to the cache with a specific expiration time
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if item has expired
	if time.Now().After(item.expiration) {
		go c.Delete(key) // Clean up expired item
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]item)
}

func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, v := range c.items {
		if now.After(v.expiration) {
			delete(c.items, k)
		}
	}
}

// StartCleanupTask starts a goroutine that periodically cleans up expired items
func (c *Cache) StartCleanupTask(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.Cleanup()
		}
	}()
}
