package cache

import (
	"os"

	"github.com/pglet/pglet/internal/config"
)

type cacher interface {
	exists(key string) bool
	getString(key string) string
	setString(key string, value string, expireSeconds int)
	inc(key string, by int) int
	remove(keys ...string)
	// hashes
	hashSet(key string, fields ...string)
	hashGet(key string, field string) string
	hashGetAll(key string) map[string]string
	hashRemove(key string, fields ...string)
	// sets
	setGet(key string) []string
	setAdd(key string, value string)
	setRemove(key string, value string)
	// sorted sets
	sortedSetAdd(key string, value string, score int64)
	sortedSetPopRange(key string, min int64, max int64) []string
	sortedSetRemove(key string, value string)
	// pubsub
	subscribe(channel string) chan []byte
	unsubscribe(ch chan []byte)
	send(channel string, message []byte)
	// locks
	lock(key string) Unlocker
}

type Unlocker interface {
	Unlock()
}

var cache cacher

func Init() {
	redisAddr := os.Getenv(config.RedisAddr)
	redisPassword := os.Getenv(config.RedisPassword)

	if redisAddr != "" {
		cache = newRedisCache(redisAddr, redisPassword)
	} else {
		cache = newMemoryCache()
	}
}

func Exists(key string) bool {
	return cache.exists(key)
}

func GetString(key string) string {
	return cache.getString(key)
}

func SetString(key string, value string, expireSeconds int) {
	cache.setString(key, value, expireSeconds)
}

func Inc(key string, by int) int {
	return cache.inc(key, by)
}

//
// Lists
// =============================

func ListReplaceFifo(key string, value string, maxSize int, expireSeconds int) []string {
	return nil
}

//
// Hashes
// =============================

func HashSet(key string, fields ...string) {
	cache.hashSet(key, fields...)
}

func HashGet(key string, field string) string {
	return cache.hashGet(key, field)
}

func HashGetAll(key string) map[string]string {
	return cache.hashGetAll(key)
}

func HashRemove(key string, fields ...string) {
	cache.hashRemove(key, fields...)
}

//
// Sets
// =============================

func SetGet(key string) []string {
	return cache.setGet(key)
}

func SetAdd(key string, value string) {
	cache.setAdd(key, value)
}

func SetRemove(key string, value string) {
	cache.setRemove(key, value)
}

func Remove(keys ...string) {
	cache.remove(keys...)
}

//
// Sorted sets
// =============================

func SortedSetAdd(key string, value string, score int64) {
	cache.sortedSetAdd(key, value, score)
}

func SortedSetPopRange(key string, min int64, max int64) []string {
	return cache.sortedSetPopRange(key, min, max)
}

func SortedSetRemove(key string, value string) {
	cache.sortedSetRemove(key, value)
}

//
// PubSub
// =============================

func Subscribe(channel string) chan []byte {
	return cache.subscribe(channel)
}

func Unsubscribe(ch chan []byte) {
	cache.unsubscribe(ch)
}

func Send(channel string, message []byte) {
	cache.send(channel, message)
}

//
// Locks
// =============================
func Lock(key string) Unlocker {
	return cache.lock(key)
}
