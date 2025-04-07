package caching

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	//RedisDB is the redis client implementation of cache
	RedisDB struct {
		database *redis.Client
		config   *cacheConfig
	}
)

// SetUpRedisDB initializes the redis instance and makes a ping to check the connection
func SetUpRedisDB(c *cacheConfig) Cache {
	r := &RedisDB{config: c}
	defaultOpts := getRedisDefaultOpt()

	if r.config.Enabled == nil {
		r.config.Enabled = defaultOpts.Enabled
	}
	if r.config.HashKeys == nil {
		r.config.HashKeys = defaultOpts.HashKeys
	}

	if r.config.Addr == "" {
		r.config.Addr = defaultOpts.Addr
	}
	if r.config.Prefix == "" {
		r.config.Prefix = defaultOpts.Prefix
	}

	if r.config.DefaultExpiration == nil {
		r.config.DefaultExpiration = defaultOpts.DefaultExpiration
	}

	r.database = redis.NewClient(&redis.Options{Addr: r.config.Addr})
	status := r.database.Ping(context.Background())
	if status.Err() != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Failed to open cache connection", slog.String("source", "Cache"), slog.String("error", status.Err().Error()))
		panic("Failed to open cache connection at %s\n")
	}
	return r
}

// SetKey sets a string type in redis with a provided value and a ttl.
func (r *RedisDB) SetKey(key string, value string, ttl *time.Duration) error {

	var expiration time.Duration

	if ttl != nil {
		expiration = *ttl
	} else {
		expiration = *r.config.DefaultExpiration
	}

	status := r.database.Set(context.Background(), fmt.Sprintf("%s:%s", r.config.Prefix, key), value, time.Duration(expiration))
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *RedisDB) SetKeyIndex(indexKey, member string) error {
	ctx := context.Background()
	indexKey = fmt.Sprintf("%s:%s:keys", r.config.Prefix, indexKey)
	member = fmt.Sprintf("%s:%s", r.config.Prefix, member)

	if err := r.database.SAdd(ctx, indexKey, member).Err(); err != nil {
		return fmt.Errorf("failed to add member to set %q: %w", indexKey, err)
	}

	return nil
}

// DeleteCacheIndex clears the cache indexes for a provided route
func (r *RedisDB) DeleteCacheIndex(indexKey string) (int, error) {
	indexKey = fmt.Sprintf("%s:%s:keys", r.config.Prefix, indexKey)
	evictedKeys, err := r.database.Del(context.Background(), indexKey).Result()
	if err != nil {
		return int(evictedKeys), err
	}

	return int(evictedKeys), nil
}

func getRedisDefaultOpt() cacheConfig {
	midnight := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Now().Location()).Sub(time.Now())
	enabled := true
	hashKeys := false
	defaultOpt := cacheConfig{
		Addr:              "127.0.0.1:6379",
		Enabled:           &enabled,
		Prefix:            "default",
		HashKeys:          &hashKeys,
		DefaultExpiration: &midnight,
	}
	return defaultOpt
}
