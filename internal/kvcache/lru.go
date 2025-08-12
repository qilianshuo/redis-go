package kvcache

import "time"

// LRUCache is a cache that evicts the least recently used items first.
type LRUCache struct {
}

// NewLRUCache creates a new LRUCache.
func NewLRUCache() *LRUCache {
	return &LRUCache{}
}

func (c *LRUCache) Set(key string, val any) (ok bool) {
	// TODO: Implement the Set method.
	return
}

func (c *LRUCache) Get(key string) (val any, ok bool) {
	// TODO: Implement the Get method.
	return
}

func (c *LRUCache) Expire(key string, expireTime *time.Time) (ok bool) {
	// TODO: Implement the Expire method.
	return
}

func (c *LRUCache) Remove(key string) (ok bool) {
	// TODO: Implement the Remove method.
	return
}

func (c *LRUCache) Purge() {
	// TODO: Implement the Purge method.
}

func (c *LRUCache) Len() int {
	// TODO: Implement the Len method.
	return 0
}

func (c *LRUCache) Has(key string) bool {
	// TODO: Implement the Has method.
	return false
}

func (c *LRUCache) ForEach(iter func(key string, val any, expiration *time.Time) bool) {
	// TODO: Implement the ForEach method.
}
