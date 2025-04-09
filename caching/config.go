package caching

import "time"

type CacheConfig struct {
	//Addr is the cache DSN.
	// memcached default: "127.0.0.1:11211"
	// redis default: "127.0.0.1:6379"
	Addr string `json:"addr" xml:"addr" yaml:"addr" csv:"addr"`
	// Enabled caching should be enabled across the application
	Enabled *bool `json:"enabled" xml:"enabled" yaml:"enabled" csv:"enabled"`
	// Prefix is primarily used for namespacing your application.
	// ie. "myProject:api"
	Prefix string `json:"prefix" xml:"prefix" yaml:"prefix" csv:"prefix"`
	// HashKeys hashes the keys befora adding them.
	// it avoids issues related to long key names & hides the parameters used
	HashKeys *bool `json:"hash_keys" xml:"hash_keys" yaml:"hash_keys" csv:"hash_keys"`
	//DefaultExpiration specifies a default value for the expiration date of a key.
	//The default value is today at midnight.
	DefaultExpiration *time.Duration `json:"default_expiration" xml:"default_expiration" yaml:"default_expiration" csv:"default_expiration"`
}
