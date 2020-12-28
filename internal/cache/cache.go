package cache

import (
	"encoding/json"
	"os"
)

type cacher interface {
	exists(key string) bool
	getString(key string) string
	setString(key string, value string, expireSeconds int)
	inc(key string, by int) int
	hashSet(key string, fields ...string)
	hashGet(key string, field string) string
	hashGetAll(key string) map[string]string
	hashRemove(key string, fields ...string)
	setGet(key string) []string
	setAdd(key string, value string)
	setRemove(key string, value string)
	remove(key string)
}

var cache cacher

func Init() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

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

func GetObject(key string, result interface{}) {
	s := cache.getString(key)
	json.Unmarshal([]byte(s), result)
}

func SetObject(key string, value interface{}, expireSeconds int) {
	payload, _ := json.Marshal(value)
	cache.setString(key, string(payload), expireSeconds)
}

func Inc(key string, by int) int {
	return cache.inc(key, by)
}

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

func SetGet(key string) []string {
	return cache.setGet(key)
}

func SetAdd(key string, value string) {
	cache.setAdd(key, value)
}

func SetRemove(key string, value string) {
	cache.setRemove(key, value)
}

func Remove(key string) {
	cache.remove(key)
}
