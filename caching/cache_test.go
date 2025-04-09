package caching

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCache(t *testing.T) {
	tests := []struct {
		cacheType   string
		expectPanic bool
	}{
		{"redis", false},
		{"memcached", false},
		{"unknown", true},
	}

	for _, test := range tests {
		cacheType := test.cacheType

		if test.expectPanic {
			assert.Panics(t, func() {
				GetCache(test.cacheType, &CacheConfig{Prefix: "tmp:test"})
			}, "Expected panic for cache type: %s", cacheType)
		} else {
			cache := assert.NotPanics(t, func() {
				cache := GetCache(test.cacheType, &CacheConfig{Prefix: "tmp:test"})
				require.NotNil(t, cache, "Cache instance should not be nil for type: %s", cacheType)
			}, "Unexpected panic for cache type: %s", cacheType)
			_ = cache
		}
	}
}

func TestSetKey(t *testing.T) {
	testName := "TestSetKey"
	tests := []struct {
		cacheType   string
		expectPanic bool
		prefix      string
		key         string
	}{
		{
			cacheType:   "redis",
			expectPanic: false,
			prefix:      "test:" + testName,
			key:         "hello",
		},
		{
			cacheType:   "memcached",
			expectPanic: false,
			prefix:      "test:" + testName,
			key:         "hello",
		},
	}

	for _, test := range tests {
		cacheType := test.cacheType

		cache := assert.NotPanics(t, func() {
			cache := GetCache(test.cacheType, &CacheConfig{Prefix: test.prefix})
			err := cache.SetKey(test.key, test.prefix, nil)
			assert.Nil(t, err)
			require.NotNil(t, cache, "Cache instance should not be nil for type: %s", cacheType)
		}, "Unexpected panic for cache type: %s", cacheType)
		_ = cache

	}
}

func TestSetKeyIndexAndDeleteCacheIndex(t *testing.T) {
	testName := "TestSetKeyIndexAndDeleteCacheIndex"
	tests := []struct {
		cacheType   string
		expectPanic bool
		prefix      string
		key         string
	}{
		{
			cacheType:   "redis",
			expectPanic: false,
			prefix:      "test:" + testName,
			key:         "hello",
		},
		{
			cacheType:   "memcached",
			expectPanic: false,
			prefix:      "test:" + testName,
			key:         "hello",
		},
	}

	for _, test := range tests {
		cacheType := test.cacheType

		cache := assert.NotPanics(t, func() {
			cache := GetCache(test.cacheType, &CacheConfig{Prefix: test.prefix})
			err := cache.SetKeyIndex(test.key, test.prefix)
			assert.Nil(t, err)

			evictedKeys, err := cache.DeleteCacheIndex(test.key)
			assert.Nil(t, err)
			assert.NotZero(t, evictedKeys)

		}, "Unexpected panic for cache type: %s", cacheType)
		_ = cache

	}
}

func TestGetKey(t *testing.T) {
	testName := "TestGetKey"
	tests := []struct {
		cacheType   string
		expectPanic bool
		prefix      string
		key         string
	}{
		{
			cacheType:   "redis",
			expectPanic: false,
			prefix:      "test:" + testName,
			key:         "hello",
		},
		{
			cacheType:   "memcached",
			expectPanic: false,
			prefix:      "test:" + testName,
			key:         "hello",
		},
	}

	for _, test := range tests {
		cacheType := test.cacheType

		cache := assert.NotPanics(t, func() {
			cache := GetCache(test.cacheType, &CacheConfig{Prefix: test.prefix})
			err := cache.SetKey(test.key, test.prefix, nil)
			assert.Nil(t, err)

			_, err = cache.Get(test.key)

			assert.Nil(t, err)

		}, "Unexpected panic for cache type: %s", cacheType)
		_ = cache

	}
}

func TestCacheIncrement(t *testing.T) {
	testName := "TestCacheIncrement"
	tests := []struct {
		cacheType   string
		expectPanic bool
		prefix      string
		key         string
	}{
		{
			cacheType:   "redis",
			expectPanic: false,
			prefix:      "test:" + testName,
			key:         "ip:123",
		},
		{
			cacheType:   "memcached",
			expectPanic: false,
			prefix:      "test:" + testName,
			key:         "ip:123",
		},
	}

	for _, test := range tests {
		cacheType := test.cacheType

		cache := assert.NotPanics(t, func() {
			cache := GetCache(test.cacheType, &CacheConfig{Prefix: test.prefix})
			err := cache.CacheIncrement(test.key, 360*time.Second)
			assert.Nil(t, err)

		}, "Unexpected panic for cache type: %s", cacheType)
		_ = cache

	}
}

func TestDefaultPrefix(t *testing.T) {
	tests := []struct {
		cacheType   string
		expectPanic bool
		prefix      string
		key         string
		HashKeys    bool
	}{
		{
			cacheType:   "redis",
			expectPanic: false,
			key:         "hello",
			HashKeys:    false,
		},
		{
			cacheType:   "memcached",
			expectPanic: false,
			key:         "hello",
			HashKeys:    true,
		},
	}

	for _, test := range tests {
		cacheType := test.cacheType

		cache := assert.NotPanics(t, func() {
			cache := GetCache(test.cacheType, &CacheConfig{HashKeys: &test.HashKeys})
			err := cache.SetKey(test.key, test.key, nil)
			assert.Nil(t, err)
			require.NotNil(t, cache, "Cache instance should not be nil for type: %s", cacheType)
		}, "Unexpected panic for cache type: %s", cacheType)
		_ = cache

	}
}
