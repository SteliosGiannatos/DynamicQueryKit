package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type (
	//MemcachedDB provides a struct for the memcache implementation of caching
	MemcachedDB struct {
		database *memcache.Client
		config   *cacheConfig
	}
)

// SetUp initializes the memcache connection
func SetUpMemcachedDB(opts *cacheConfig) *MemcachedDB {
	m := &MemcachedDB{config: opts}
	defaultOpts := getMemcachedDefaultOpt()

	if m.config.Enabled == nil {
		m.config.Enabled = defaultOpts.Enabled
	}
	if m.config.HashKeys == nil {
		m.config.HashKeys = defaultOpts.HashKeys
	}

	if m.config.Addr == "" {
		m.config.Addr = defaultOpts.Addr
	}
	if m.config.Prefix == "" {
		m.config.Prefix = defaultOpts.Prefix
	}
	if m.config.DefaultExpiration == nil {
		m.config.DefaultExpiration = defaultOpts.DefaultExpiration
	}

	m.database = memcache.New(m.config.Addr)
	if m.database == nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Failed to open cache connection", slog.String("source", "Cache"))
		panic("Failed to open cache connection")
	}

	err := m.database.Ping()
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Failed to open cache connection", slog.String("source", "Cache"), slog.String("error", err.Error()))
		panic("Failed to open cache connection")
	}
	return m
}

// SetKey easier way to set a cache key
// Default expiration date is today at midnight
func (m *MemcachedDB) SetKey(key string, value string, ttl *time.Duration) error {
	var expiration int32

	if ttl != nil {
		expiration = int32(ttl.Seconds())
	} else {
		expiration = int32(m.config.DefaultExpiration.Seconds())
	}

	err := m.database.Set(&memcache.Item{Key: fmt.Sprintf("%s:%s", m.config.Prefix, key), Value: []byte(value), Expiration: expiration})
	if err != nil {
		return err
	}
	return nil
}

// SetKeyIndex appends to the list of keys
// A list of the cached keys is maintained in the cache with no expiration
// so when it comes to invalidating routes with dynamic filters you know all the cached keys
func (m *MemcachedDB) SetKeyIndex(indexKey string, member string) error {
	member = fmt.Sprintf("%s:%s", m.config.Prefix, member)
	indexKey = fmt.Sprintf("%s:%s:keys", m.config.Prefix, indexKey)
	var members []string

	item, err := m.database.Get(indexKey)
	if err == memcache.ErrCacheMiss {
		members = []string{}

	} else if err != nil {
		return err
	} else {
		err = json.Unmarshal(item.Value, &members)
		if err != nil {
			return err
		}

	}

	// make sure the key does not already exist
	for _, existingKey := range members {
		if existingKey == member {
			slog.LogAttrs(context.Background(), slog.LevelDebug, "index key already exists", slog.String("source", indexKey), slog.String("existingKey", existingKey), slog.String("cacheKey", member), slog.String("index key", indexKey))
			return nil
		}
	}

	members = append(members, member)

	jsonMembers, _ := json.Marshal(members)

	err = m.database.Set(&memcache.Item{Key: indexKey, Value: jsonMembers, Expiration: 0})
	if err != nil {
		return err
	}
	return nil
}

// DeleteCacheIndex clears the cache indexes for a provided route
func (m *MemcachedDB) DeleteCacheIndex(indexKey string) (int, error) {
	var members []string
	evictedKeys := 0
	indexKey = fmt.Sprintf("%s:%s:keys", m.config.Prefix, indexKey)

	item, err := m.database.Get(indexKey)
	if err == memcache.ErrCacheMiss {
		slog.LogAttrs(context.Background(), slog.LevelDebug, "no cache for provided key", slog.String("key", indexKey), slog.String("error", err.Error()))
		return evictedKeys, nil
	}

	err = json.Unmarshal(item.Value, &members)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "error deleting cache index", slog.String("key", indexKey), slog.String("error", err.Error()))
		return evictedKeys, err
	}

	if len(members) == 0 {
		slog.LogAttrs(context.Background(), slog.LevelDebug, "no keys under", slog.String("key", indexKey))
		return evictedKeys, nil
	}

	for _, value := range members {
		slog.LogAttrs(context.Background(), slog.LevelDebug, "deleting cache", slog.String("key", value))
		if err = m.database.Delete(value); err == nil || err == memcache.ErrCacheMiss {
			evictedKeys++
		}
	}

	slog.LogAttrs(context.Background(), slog.LevelDebug, "deleting index key", slog.String("key", indexKey), slog.Int("members evicted", evictedKeys))
	err = m.database.Delete(indexKey)
	evictedKeys++
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelWarn, "error deleting cache index", slog.String("key", indexKey))
	}

	return evictedKeys, nil
}

func getMemcachedDefaultOpt() cacheConfig {
	midnight := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 0, 0, 0, 0, time.Now().Location()).Sub(time.Now())
	enabled := true
	hashKeys := false
	defaultOpt := cacheConfig{
		Addr:              "127.0.0.1:11211",
		Enabled:           &enabled,
		Prefix:            "default",
		HashKeys:          &hashKeys,
		DefaultExpiration: &midnight,
	}
	return defaultOpt
}
