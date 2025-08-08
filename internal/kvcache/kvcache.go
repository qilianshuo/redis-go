package kvcache

import (
	"time"

	"github.com/mirage208/redis-go/common/datastruct/dict"
)

const (
	dataDictSize = 1 << 16 // 64K
	ttlDictSize  = 1 << 16 // 64K
)

// KVCache is a simple key-value cache structure
type KVCache struct {
	data *dict.ConcurrentDict
	ttl  *dict.ConcurrentDict
}

func NewKVCache() *KVCache {
	return &KVCache{
		data: dict.NewConcurrentDict(dataDictSize),
		ttl:  dict.NewConcurrentDict(ttlDictSize),
	}
}

// GetEntity retrieves a value by key from the cache.
// It returns the value and a boolean indicating whether the key exists.
// If the key does not exist, it returns nil and false.
// If the key exists, it returns the value and true.
// If the key has an expiration time set, it checks if the current time is past that expiration time.
func (c *KVCache) GetEntity(key string) (value any, ok bool) {
	value, ok = c.data.Get(key)
	if !ok {
		return
	}
	if expireTime, ok := c.ttl.Get(key); ok {
		if time.Now().After(expireTime.(time.Time)) {
			c.data.Remove(key)
			c.ttl.Remove(key)
			return nil, false
		}
	}
	return
}

// PutEntity inserts or updates a key-value pair in the cache.
func (c *KVCache) PutEntity(key string, value any) (ok bool) {
	ret := c.data.Put(key, value)
	if ret > 0 {
		ok = true
	} else {
		ok = false
	}
	return
}

// PutIfAbsent inserts a key-value pair if the key does not already exist and returns true if the insertion was successful.
func (c *KVCache) PutIfAbsent(key string, value any) (ok bool) {
	ret := c.data.PutIfAbsent(key, value)
	if ret > 0 {
		ok = true
	} else {
		ok = false
	}
	return
}

// PutIfExists updates the value for an existing key and returns true if the key existed.
func (c *KVCache) PutIfExists(key string, value any) (ok bool) {
	ret := c.data.PutIfExists(key, value)
	if ret > 0 {
		ok = true
	} else {
		ok = false
	}
	return
}

// Expire sets the expiration time for a key.
func (c *KVCache) Expire(key string, expireTime time.Time) {
	c.ttl.Put(key, expireTime)
}

// Persist removes the expiration time for a key, making it persistent.
func (c *KVCache) Persist(key string) {
	c.ttl.Remove(key)
}
