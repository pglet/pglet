package cache

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wangjia184/sortedset"
)

type cacheEntry struct {
	expires time.Time
	data    interface{}
}

type lockEntry struct {
	m     *memoryCache // point back to M, so we can synchronize removing this mentry when cnt==0
	el    sync.Mutex   // entry-specific lock
	count int          // reference count
	key   interface{}  // key in ma
}

type memoryCache struct {
	sync.RWMutex
	entries map[string]*cacheEntry
	// pubsub
	// subscriber is a channel
	pubsubLock         sync.RWMutex
	channelSubscribers map[string]map[chan []byte]bool
	subscribers        map[chan []byte]string
	// locks
	ml          sync.Mutex                 // lock for entry map
	lockEntries map[interface{}]*lockEntry // entry map
}

func newMemoryCache() cacher {
	log.Println("Using in-memory cache")
	mc := &memoryCache{
		entries:            make(map[string]*cacheEntry),
		channelSubscribers: make(map[string]map[chan []byte]bool),
		subscribers:        make(map[chan []byte]string),
		lockEntries:        make(map[interface{}]*lockEntry),
	}
	//go mc.dumpData()
	return mc
}

func (c *memoryCache) dumpData() {
	dataPath := filepath.Join(os.TempDir(), "pglet-memory-cache")
	log.Println("Memory cache dump:", dataPath)

	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C

		log.Println("MEMORY STATE DUMP")
		log.Println("channel subscribers:")
		for k, v := range c.channelSubscribers {
			log.Println("channel:", k, "subscribers:", len(v))
		}
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

//
// Hashes
// =============================

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

//
// Sets
// =============================

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

//
// Sorted Sets
// =============================

func (c *memoryCache) sortedSetAdd(key string, value string, score int64) {
	c.Lock()
	defer c.Unlock()

	var set *sortedset.SortedSet
	entry := c.getEntry(key)
	if entry == nil {
		entry = c.newEntry(0)
		c.entries[key] = entry
		set = sortedset.New()
	} else {
		set = entry.data.(*sortedset.SortedSet)
	}

	set.AddOrUpdate(value, sortedset.SCORE(score), nil)
	entry.data = set
}

func (c *memoryCache) sortedSetPopRange(key string, min int64, max int64) []string {
	c.Lock()
	defer c.Unlock()

	entry := c.getEntry(key)
	if entry == nil {
		return make([]string, 0)
	}
	set := entry.data.(*sortedset.SortedSet)
	nodes := set.GetByScoreRange(sortedset.SCORE(min), sortedset.SCORE(max), &sortedset.GetByScoreRangeOptions{})
	result := make([]string, len(nodes))
	for i, node := range nodes {
		result[i] = node.Key()
		set.Remove(result[i])
	}
	return result
}

func (c *memoryCache) sortedSetRemove(key string, value string) {
	c.Lock()
	defer c.Unlock()

	entry := c.getEntry(key)
	if entry == nil {
		return
	}
	set := entry.data.(*sortedset.SortedSet)
	set.Remove(value)
	if set.GetCount() == 0 {
		c.deleteEntry(key)
	}
}

func (c *memoryCache) remove(keys ...string) {
	c.Lock()
	defer c.Unlock()

	for _, key := range keys {
		c.deleteEntry(key)
	}
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

//
// PubSub
// =============================

func (c *memoryCache) subscribe(channel string) chan []byte {
	c.pubsubLock.Lock()
	defer c.pubsubLock.Unlock()

	subscribers := c.channelSubscribers[channel]
	if subscribers == nil {
		subscribers = make(map[chan []byte]bool)
		c.channelSubscribers[channel] = subscribers
	}

	ch := make(chan []byte)
	subscribers[ch] = true
	c.subscribers[ch] = channel
	return ch
}

func (c *memoryCache) unsubscribe(ch chan []byte) {
	c.pubsubLock.Lock()
	defer c.pubsubLock.Unlock()

	channel := c.subscribers[ch]
	if channel == "" {
		return
	}

	subscribers := c.channelSubscribers[channel]
	if subscribers == nil {
		return
	}

	close(ch)
	delete(subscribers, ch)

	if len(subscribers) == 0 {
		delete(c.channelSubscribers, channel)
	}
}

func (c *memoryCache) send(channel string, message []byte) {
	c.pubsubLock.RLock()
	defer c.pubsubLock.RUnlock()

	subscribers := c.channelSubscribers[channel]
	if subscribers == nil {
		return
	}

	for ch := range subscribers {
		select {
		case ch <- message:
			// Message sent to subscriber
		default:
			// No listeners
		}
	}
}

//
// Locks
// Source: https://stackoverflow.com/questions/40931373/how-to-gc-a-map-of-mutexes-in-go
// =============================
func (c *memoryCache) lock(key string) Unlocker {

	// read or create entry for this key atomically
	c.ml.Lock()
	e, ok := c.lockEntries[key]
	if !ok {
		e = &lockEntry{m: c, key: key}
		c.lockEntries[key] = e
	}
	e.count++ // ref count
	c.ml.Unlock()

	// acquire lock, will block here until e.cnt==1
	e.el.Lock()

	return e
}

// Unlock releases the lock for this entry.
func (me *lockEntry) Unlock() {

	m := me.m

	//log.Println("LOCK ENTRIES:", len(me.m.lockEntries))

	// decrement and if needed remove entry atomically
	m.ml.Lock()
	e, ok := m.lockEntries[me.key]
	if !ok { // entry must exist
		m.ml.Unlock()
		log.Errorf("Unlock requested for key=%v but no entry found", me.key)
	}
	e.count--        // ref count
	if e.count < 1 { // if it hits zero then we own it and remove from map
		delete(m.lockEntries, me.key)
	}
	m.ml.Unlock()

	// now that map stuff is handled, we unlock and let
	// anything else waiting on this key through
	e.el.Unlock()
}
