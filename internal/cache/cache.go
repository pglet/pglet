package cache

import (
	"encoding/json"
	"os"
)

type cacher interface {
	getString(key string) string
	setString(key string, value string, expireSeconds int)
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
	return 0
}

func HashSet(key string, args ...interface{}) {

}

func HashGetString(key string, field string) string {
	// TODO
	return ""
}

func HashGetAll(key string) map[string]string {
	return nil
}

func HashRemove(key string, args ...string) {

}

func SetGet(key string) []string {
	// TODO
	return nil
}

func SetAdd(key string, value string) {
	// TODO
}

func SetRemove(key string, value string) {
	// TODO
}

func Remove(key string) {
	// TODO
}
