package kvcache

import "time"

// KVCache is the interface for a key-value cache
type KVCache interface {
	// Set inserts or updates the specified key-value pair.
	Set(key string, val any) (ok bool)
	// Get returns the value for the specified key if it is present in the cache.
	Get(key string) (val any, ok bool)
	// Expire sets the expiration time for the specified key.
	Expire(key string, expireTime *time.Time) (ok bool)
	// Remove removes the specified key from the cache if the key is present.
	// Returns true if the key was present and the key has been deleted.
	Remove(key string) (ok bool)
	// Purge removes all key-value pairs from the cache.
	Purge()
	// Len returns the number of items in the cache.
	Len() int
	// Has returns true if the key exists in the cache.
	Has(key string) bool
	// ForEach iterates over all key-value pairs in the cache.
	ForEach(iter func(key string, val any, expiration *time.Time) bool)
}
