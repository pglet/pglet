package cache

import "os"

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
