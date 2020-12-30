package cache

import (
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/pglet/pglet/internal/config"
	log "github.com/sirupsen/logrus"
)

var errLockMismatch = errors.New("key is locked with a different secret")

const lockScript = `
local v = redis.call("GET", KEYS[1])
if v == false or v == ARGV[1]
then
	return redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2]) and 1
else
	return 0
end
`

const unlockScript = `
local v = redis.call("GET",KEYS[1])
if v == false then
	return 1
elseif v == ARGV[1] then
	return redis.call("DEL",KEYS[1])
else
	return 0
end
`

type redisCache struct {
	pool *redis.Pool
}

func newRedisCache() cacher {

	redisAddr := config.RedisAddr()
	redisPassword := config.RedisPassword()

	log.Println("Connecting to Redis server", redisAddr)

	pool := &redis.Pool{
		MaxIdle:   config.RedisMaxIdle(),
		MaxActive: config.RedisMaxActive(),

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisAddr)
			if err != nil {
				return nil, err
			}

			if redisPassword != "" {
				if _, err := conn.Do("AUTH", redisPassword); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, nil
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return &redisCache{
		pool: pool,
	}
}

func (c *redisCache) exists(key string) bool {
	conn := c.pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Fatal(err)
	}
	return exists
}

func (c *redisCache) getString(key string) string {
	conn := c.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func (c *redisCache) setString(key string, value string, expireSeconds int) {
	conn := c.pool.Get()
	defer conn.Close()

	args := []interface{}{key, value}
	if expireSeconds > 0 {
		args = append(args, "EX", expireSeconds)
	}
	_, err := conn.Do("SET", args...)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *redisCache) inc(key string, by int) int {
	conn := c.pool.Get()
	defer conn.Close()

	value, err := redis.Int(conn.Do("INCRBY", key, by))
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func (c *redisCache) remove(keys ...string) {
	conn := c.pool.Get()
	defer conn.Close()

	args := make([]interface{}, len(keys))
	for i, f := range keys {
		args[i] = f
	}

	_, err := conn.Do("DEL", args...)
	if err != nil {
		log.Fatal(err)
	}
}

//
// Hashes
// =============================

func (c *redisCache) hashSet(key string, fields ...string) {
	conn := c.pool.Get()
	defer conn.Close()

	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, f := range fields {
		args[i+1] = f
	}

	_, err := conn.Do("HSET", args...)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *redisCache) hashGet(key string, field string) string {
	conn := c.pool.Get()
	defer conn.Close()

	value, err := conn.Do("HGET", key, field)
	if err != nil {
		log.Fatal(err)
	}

	if value == nil {
		return ""
	}

	return string(value.([]byte))
}

func (c *redisCache) hashGetAll(key string) map[string]string {
	conn := c.pool.Get()
	defer conn.Close()

	value, err := redis.StringMap(conn.Do("HGETALL", key))
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func (c *redisCache) hashRemove(key string, fields ...string) {
	conn := c.pool.Get()
	defer conn.Close()

	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, f := range fields {
		args[i+1] = f
	}

	_, err := conn.Do("HDEL", args...)
	if err != nil {
		log.Fatal(err)
	}
}

//
// Sets
// =============================

func (c *redisCache) setGet(key string) []string {
	conn := c.pool.Get()
	defer conn.Close()

	value, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func (c *redisCache) setAdd(key string, value string) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *redisCache) setRemove(key string, value string) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SREM", key, value)
	if err != nil {
		log.Fatal(err)
	}
}

//
// Sorted Sets
// =============================

func (c *redisCache) sortedSetAdd(key string, value string, score int64) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("ZADD", key, score, value)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *redisCache) sortedSetPopRange(key string, min int64, max int64) []string {
	conn := c.pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("ZRANGEBYSCORE", key, min, max)
	conn.Send("ZREMRANGEBYSCORE", key, min, max)
	value, err := redis.Strings(conn.Do("EXEC"))

	if err != nil {
		log.Fatal(err)
	}
	return value
}

func (c *redisCache) sortedSetRemove(key string, value string) {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("ZREM", key, value)
	if err != nil {
		log.Fatal(err)
	}
}

//
// PubSub
// =============================

func (c *redisCache) subscribe(channel string) chan []byte {
	return nil
}

func (c *redisCache) unsubscribe(ch chan []byte) {

}

func (c *redisCache) send(channel string, message []byte) {

}

//
// Locks
// Source: https://gist.github.com/bgentry/6105288
// =============================
func (c *redisCache) lock(key string) Unlocker {
	attempts := 100
	lockTimeout := 10 * time.Second
	retryTimeout := 100 * time.Millisecond
	secret := uuid.New().String()

	for i := 0; i < attempts; i++ {
		conn := c.pool.Get()
		if c.writeLock(conn, key, secret, int64(lockTimeout)) {
			return &redisLock{
				conn:   conn,
				name:   key,
				secret: secret,
			}
		}
		time.Sleep(retryTimeout)
	}
	log.Fatalf("Cannot aquire lock %s in 10 seconds", key)
	return nil
}

// writeLock attempts to grab a redis lock. The error returned is safe to ignore
// if all you care about is whether or not the lock was acquired successfully.
func (c *redisCache) writeLock(conn redis.Conn, name, secret string, ttl int64) bool {

	script := redis.NewScript(1, lockScript)
	resp, err := redis.Int(script.Do(conn, name, secret, ttl))
	if err != nil || resp == 0 {
		conn.Close()
		return false
	}
	return true
}

type redisLock struct {
	conn   redis.Conn
	name   string
	secret string
}

func (rl *redisLock) Unlock() {
	defer rl.conn.Close()

	script := redis.NewScript(1, unlockScript)
	resp, err := redis.Int(script.Do(rl.conn, rl.name, rl.secret))
	if err != nil {
		log.Fatal(err)
	}
	if resp == 0 {
		log.Fatal(errLockMismatch)
	}
}
