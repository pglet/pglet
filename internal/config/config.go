package config

import (
	"os"
	"strconv"
)

const (

	// pages/sessions
	defaultPageLifetimeMinutes = 1440
	defaultAppLifetimeMinutes  = 60
	pageLifetimeMinutes        = "PAGE_LIFETIME_MINUTES"
	appLifetimeMinutes         = "APP_LIFETIME_MINUTES"

	// redis
	defaultRedisMaxIdle   = 5
	defaultRedisMaxActive = 10
	redisAddr             = "REDIS_ADDR"
	redisPassword         = "REDIS_PASSWORD"
	redisMaxIdle          = "REDIS_MAX_IDLE"
	redisMaxActive        = "REDIS_MAX_ACTIVE"
)

func RedisAddr() string {
	return os.Getenv(redisAddr)
}

func RedisPassword() string {
	return os.Getenv(redisPassword)
}

func RedisMaxIdle() int {
	if n, err := strconv.Atoi(os.Getenv(redisMaxIdle)); err == nil {
		return n
	}
	return defaultRedisMaxIdle
}

func RedisMaxActive() int {
	if n, err := strconv.Atoi(os.Getenv(redisMaxActive)); err == nil {
		return n
	}
	return defaultRedisMaxActive
}

func PageLifetimeMinutes() int {
	if n, err := strconv.Atoi(os.Getenv(pageLifetimeMinutes)); err == nil {
		return n
	}
	return defaultPageLifetimeMinutes
}

func AppLifetimeMinutes() int {
	if n, err := strconv.Atoi(os.Getenv(appLifetimeMinutes)); err == nil {
		return n
	}
	return defaultAppLifetimeMinutes
}
