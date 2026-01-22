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
	entries map[string]cacheEntry
	mu      sync.Mutex
}

func NewCache(cleanupInterval time.Duration) Cache {
	return Cache{make(map[string]cacheEntry), sync.Mutex{}}
}

func (c *Cache) Add(link string, data []byte) {
	c.entries[link] = cacheEntry{
		createdAt: time.Now(),
		data:      data,
	}
}

func (c *Cache) Get(link string) ([]byte, bool) {
	entry, ok := c.entries[link]
	return entry.data, ok
}
