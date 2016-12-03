package sessions

import (
	"time"

	redis "gopkg.in/redis.v5"
)

// RedisStorer is a session storer implementation for saving sessions
// to a Redis database.
type RedisStorer struct {
	// How long sessions take to expire in Redis
	maxAge time.Duration
	client *redis.Client
}

// NewDefaultRedisStorer takes a bind address of the Redis server host:port and
// returns a RedisStorer object with default values.
// The default values are:
// Addr: localhost:6379
// Password: no password
// DB: First database (0) to be selected after connecting to Redis
// maxAge: 1 week (clear session stored in Redis after 1 week)
func NewDefaultRedisStorer(addr, password string, db int) (*RedisStorer, error) {
	if addr == "" {
		addr = "localhost:6379"
	}
	opts := redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
	return NewRedisStorer(opts, time.Hour*24*7)
}

// NewRedisStorer initializes and returns a new RedisStorer object.
// It takes a bind address of the Redis server host:port and the maxAge of how
// long each session should live in the Redis server.
// Persistent storage can be attained by setting maxAge to zero.
func NewRedisStorer(opts redis.Options, maxAge time.Duration) (*RedisStorer, error) {
	r := &RedisStorer{
		maxAge: maxAge,
		client: redis.NewClient(&opts),
	}

	return r, nil
}

// All keys in the redis store
func (r *RedisStorer) All() ([]string, error) {
	var sessions []string

	iter := r.client.Scan(0, "", 0).Iterator()
	for iter.Next() {
		sessions = append(sessions, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return sessions, err
	}

	return sessions, nil
}

// Get returns the value string saved in the session pointed to by the
// session id key.
func (r *RedisStorer) Get(key string) (value string, err error) {
	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return "", errNoSession{}
	} else if err != nil {
		return "", err
	}

	return val, nil
}

// Set saves the value string to the session pointed to by the session id key.
func (r *RedisStorer) Set(key, value string) error {
	return r.client.Set(key, value, r.maxAge).Err()
}

// Del the session pointed to by the session id key and remove it.
func (r *RedisStorer) Del(key string) error {
	return r.client.Del(key).Err()
}

// ResetExpiry resets the expiry of the key
func (r *RedisStorer) ResetExpiry(key string) error {
	return nil
}
