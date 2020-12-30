package config

const (

	// default values
	DefaultPageLifetimeMinutes    = 1440
	DefaultAppLifetimeMinutes     = 60
	DefaultSessionLifetimeMinutes = 20

	// environment variables
	RedisAddr              = "REDIS_ADDR"
	RedisPassword          = "REDIS_PASSWORD"
	PageLifetimeMinutes    = "PAGE_LIFETIME_MINUTES"
	AppLifetimeMinutes     = "APP_LIFETIME_MINUTES"
	SessionLifetimeMinutes = "SESSION_LIFETIME_MINUTES"
)
