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

func (c *memoryCache) exists(key string) bool {
	c.RLock()
	defer c.RUnlock()

	entry := c.getEntry(key)
	if entry != nil {
		return true
	}
	return false
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

func (c *memoryCache) inc(key string, by int) int {
	c.Lock()
	defer c.Unlock()

	i := 0
	entry := c.getEntry(key)
	if entry == nil {
		entry = c.newEntry(0)
		c.entries[key] = entry
	} else {
		i = entry.data.(int)
	}

	i += by
	entry.data = i
	return i
}

func (c *memoryCache) hashSet(key string, fields ...string) {
	c.Lock()
	defer c.Unlock()

	var hash map[string]string
	entry := c.getEntry(key)
	if entry == nil {
		entry = c.newEntry(0)
		c.entries[key] = entry
		hash = make(map[string]string)
	} else {
		hash = entry.data.(map[string]string)
	}

	var k string
	for i, f := range fields {
		if i%2 == 0 {
			k = f
		} else if i%2 == 1 {
			hash[k] = f
		}
	}

	entry.data = hash
}

func (c *memoryCache) hashGet(key string, field string) string {
	c.RLock()
	defer c.RUnlock()

	entry := c.getEntry(key)
	if entry == nil {
		return ""
	}
	hash := entry.data.(map[string]string)
	return hash[field]
}

func (c *memoryCache) hashGetAll(key string) map[string]string {
	c.RLock()
	defer c.RUnlock()

	entry := c.getEntry(key)
	if entry == nil {
		return make(map[string]string)
	}
	return entry.data.(map[string]string)
}

func (c *memoryCache) hashRemove(key string, fields ...string) {
	c.Lock()
	defer c.Unlock()

	entry := c.getEntry(key)
	if entry == nil {
		return
	}
	hash := entry.data.(map[string]string)
	for _, f := range fields {
		delete(hash, f)
	}
	if len(hash) == 0 {
		c.deleteEntry(key)
	}
}

func (c *memoryCache) setGet(key string) []string {
	c.RLock()
	defer c.RUnlock()

	entry := c.getEntry(key)
	if entry == nil {
		return make([]string, 0)
	}
	hash := entry.data.(map[string]bool)
	result := make([]string, len(hash))
	i := 0
	for k := range hash {
		result[i] = k
		i++
	}
	return result
}

func (c *memoryCache) setAdd(key string, value string) {
	c.Lock()
	defer c.Unlock()

	var hash map[string]bool
	entry := c.getEntry(key)
	if entry == nil {
		entry = c.newEntry(0)
		c.entries[key] = entry
		hash = make(map[string]bool)
	} else {
		hash = entry.data.(map[string]bool)
	}

	hash[value] = true
	entry.data = hash
}

func (c *memoryCache) setRemove(key string, value string) {
	c.Lock()
	defer c.Unlock()

	entry := c.getEntry(key)
	if entry == nil {
		return
	}
	hash := entry.data.(map[string]bool)
	delete(hash, value)
	if len(hash) == 0 {
		c.deleteEntry(key)
	}
}

func (c *memoryCache) remove(key string) {
	c.Lock()
	defer c.Unlock()
	c.deleteEntry(key)
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

func (c *memoryCache) deleteEntry(key string) {
	delete(c.entries, key)
}
