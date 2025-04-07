// Package caching provides an easy abstraction layer
// for adding caching to an application
package caching

import "time"

type (
	//Cache All methods required to be implemented for a technology to be considered a cache
	Cache interface {
		//SetKey sets a key value pair. It prefixes the key provided
		//with the specified prefix in the cacheConfig.
		//in case the ttl is not specified it uses the default ttl specified in the config.
		//if the ttl is not specified in the config. Then the expiration is at midnight
		SetKey(key string, value string, ttl *time.Duration) error
		SetKeyIndex(route string, key string) error
		DeleteCacheIndex(indexKey string) (int, error)
	}
)

// GetCache is a factory that returns a caching layer for all
// methods that are commonly used for adding a caching layer to your API.
// gives flexibility when building and deploying the API.
// it allows for easy benchmarks and swapping technologies
// it is also easier for development to use redis for example
// in order to access the cache data with an app like RedisInsight
func GetCache(cacheType string, c *cacheConfig) Cache {
	var cache Cache
	switch cacheType {
	case "memcached":
		cache = SetUpMemcachedDB(c)
	case "redis":
		cache = SetUpRedisDB(c)
	default:
		panic("Unsupported cache type: " + cacheType)
	}
	return cache
}
