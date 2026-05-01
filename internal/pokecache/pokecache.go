package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

func NewCache(cleanupInterval time.Duration) *Cache {
	c := Cache{sync.Mutex{}, make(map[string]cacheEntry)}
	go c.reapLoop(cleanupInterval)
	return &c
}

func (c *Cache) Add(link string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[link] = cacheEntry{
		createdAt: time.Now(),
		data:      data,
	}
}

func (c *Cache) Get(link string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[link]
	return entry.data, ok
}

func (c *Cache) reapLoop(ticks time.Duration) {
	ticker := time.NewTicker(ticks)
	defer ticker.Stop()

	for t := range ticker.C {
		deadline := t.Add(-ticks)
		c.mu.Lock()
		for link, entry := range c.entries {
			if entry.createdAt.Before(deadline) {
				delete(c.entries, link)
			}
		}
		c.mu.Unlock()
	}
}
