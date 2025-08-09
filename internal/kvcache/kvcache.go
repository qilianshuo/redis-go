package kvcache

import (
	"time"
)

// KVCache is a sequential key-value cache structure
type KVCache struct {
	data map[string]*DataEntity
	ttl  map[string]*time.Time
}

func NewKVCache() *KVCache {
	return &KVCache{
		data: map[string]*DataEntity{},
		ttl:  map[string]*time.Time{},
	}
}

// DataEntity stores data bound to a key, including a string, list, hash, set and so on
type DataEntity struct {
	Data any
}

// GetEntity retrieves a value by key from the cache.
func (c *KVCache) GetEntity(key string) (entity *DataEntity, ok bool) {
	entity, ok = c.data[key]
	if !ok {
		return nil, false
	}

	if expireTime, ok := c.ttl[key]; ok {
		if time.Now().After(*expireTime) {
			delete(c.data, key)
			delete(c.ttl, key)
			return nil, false
		}
	}
	return entity, true
}

// PutEntity inserts or updates a key-value pair in the cache.
func (c *KVCache) PutEntity(key string, entity *DataEntity) (ok bool) {
	c.data[key] = entity
	return true
}

// PutIfAbsent inserts a key-value pair if the key does not already exist and returns true if the insertion was successful.
func (c *KVCache) PutIfAbsent(key string, entity *DataEntity) (ok bool) {
	if _, exists := c.data[key]; !exists {
		c.data[key] = entity
		return true
	}
	return false
}

// PutIfExists updates the value for an existing key and returns true if the key existed.
func (c *KVCache) PutIfExists(key string, entity *DataEntity) (ok bool) {
	if _, exists := c.data[key]; exists {
		c.data[key] = entity
		return true
	}
	return false
}

// Expire sets the expiration time for a key.
func (c *KVCache) Expire(key string, expireTime *time.Time) {
	c.ttl[key] = expireTime
}

// Persist removes the expiration time for a key, making it persistent.
func (c *KVCache) Persist(key string) {
	delete(c.ttl, key)
}

// ForEach iterates over all key-value pairs in the cache, applying the provided function.
func (c *KVCache) ForEach(f func(key string, entity *DataEntity, expiration *time.Time) bool) {
	for key, entity := range c.data {
		if expiration, ok := c.ttl[key]; ok {
			if !f(key, entity, expiration) {
				break
			}
		} else {
			if !f(key, entity, nil) {
				break
			}
		}
	}
}
