package cache

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type cacheEntry struct {
	expires time.Time
	data    interface{}
}

type memoryCache struct {
	sync.RWMutex
	entries map[string]*cacheEntry
}

func newMemoryCache() cacher {
	log.Println("Using in-memory cache")
	return &memoryCache{
		entries: make(map[string]*cacheEntry),
	}
}

func (c *memoryCache) getString(key string) string {
	c.RLock()
	defer c.RUnlock()

	entry := c.getEntry(key)
	if entry == nil {
		return ""
	}
	return entry.data.(string)
}

func (c *memoryCache) setString(key string, value string, expireSeconds int) {
	c.Lock()
	defer c.Unlock()

	entry := c.newEntry(expireSeconds)
	entry.data = value
	c.entries[key] = entry
}

func (c *memoryCache) newEntry(expireSeconds int) *cacheEntry {
	entry := &cacheEntry{}

	if expireSeconds > 0 {
		entry.expires = time.Now().Add(time.Duration(expireSeconds) * time.Second)
	}

	return entry
}

func (c *memoryCache) getEntry(key string) *cacheEntry {
	entry := c.entries[key]
	if entry == nil {
		return nil
	}
	if !entry.expires.IsZero() && time.Now().After(entry.expires) {
		// remove expired entry
		delete(c.entries, key)
		return nil
	}
	return entry
}
