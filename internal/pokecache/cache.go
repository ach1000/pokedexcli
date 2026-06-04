package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache is a thread-safe in-memory store that automatically evicts entries
// older than the configured interval.
type Cache struct {
	mu       sync.RWMutex
	entries  map[string]cacheEntry
	interval time.Duration
}

// Stats contains basic runtime cache metrics.
type Stats struct {
	ItemCount       int
	AverageLifetime time.Duration
}

// NewCache creates a cache that evicts entries older than interval.
// The reap loop runs in a background goroutine until the process exits.
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop()
	return c
}

// Add stores val under key, overwriting any existing entry.
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

// Get retrieves the value for key. The bool is false when the key is absent.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	return entry.val, ok
}

// Stats returns a snapshot of the cache size and average item age.
func (c *Cache) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := len(c.entries)
	if count == 0 {
		return Stats{}
	}

	now := time.Now()
	var totalLifetime time.Duration
	for _, entry := range c.entries {
		totalLifetime += now.Sub(entry.createdAt)
	}

	return Stats{
		ItemCount:       count,
		AverageLifetime: totalLifetime / time.Duration(count),
	}
}

// reapLoop evicts entries that are older than the cache interval.
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()
	cutoff := time.Now().Add(-c.interval)
	for key, entry := range c.entries {
		if entry.createdAt.Before(cutoff) {
			delete(c.entries, key)
		}
	}
}
