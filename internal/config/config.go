package config

import (
	"os"
	"strconv"
)

const (

	// general settings
	forceSSL = "FORCE_SSL"

	// pages/sessions
	defaultPageLifetimeMinutes = 1440
	defaultAppLifetimeMinutes  = 60
	pageLifetimeMinutes        = "PAGE_LIFETIME_MINUTES"
	appLifetimeMinutes         = "APP_LIFETIME_MINUTES"
	checkPageIP                = "CHECK_PAGE_IP" // unauthenticated clients only
	limitPagesPerHour          = "LIMIT_PAGES_PER_HOUR"
	limitSessionsPerHour       = "LIMIT_SESSIONS_PER_HOUR"
	limitSessionSizeBytes      = "LIMIT_SESSION_SIZE_BYTES"
	checkReservedPages         = "CHECK_RESERVED_PAGES"

	// redis
	defaultRedisMaxIdle   = 5
	defaultRedisMaxActive = 10
	redisAddr             = "REDIS_ADDR"
	redisPassword         = "REDIS_PASSWORD"
	redisMaxIdle          = "REDIS_MAX_IDLE"
	redisMaxActive        = "REDIS_MAX_ACTIVE"
)

func ForceSSL() bool {
	if n, err := strconv.ParseBool(os.Getenv(forceSSL)); err == nil {
		return n
	}
	return false
}

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

func CheckPageIP() bool {
	return os.Getenv(checkPageIP) == "true"
}

func CheckReservedPages() bool {
	return os.Getenv(checkReservedPages) == "true"
}

func LimitPagesPerHour() int {
	if n, err := strconv.Atoi(os.Getenv(limitPagesPerHour)); err == nil {
		return n
	}
	return 0
}

func LimitSessionsPerHour() int {
	if n, err := strconv.Atoi(os.Getenv(limitSessionsPerHour)); err == nil {
		return n
	}
	return 0
}

func LimitSessionSizeBytes() int {
	if n, err := strconv.Atoi(os.Getenv(limitSessionSizeBytes)); err == nil {
		return n
	}
	return 0
}
