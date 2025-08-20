package kvcache

import (
	"container/list"
	"time"
)

// LRUCache is a cache that evicts the least recently used items first.
type LRUCache struct {
	capability int
	cache      map[string]*list.Element
	ttl        map[string]*time.Time
	evictList  *list.List
}

// entry represents a single entry in the cache.
type entry struct {
	key   string
	value any
}

// NewLRUCache creates a new LRUCache.
func NewLRUCache(capability int) *LRUCache {
	return &LRUCache{
		capability: capability,
		cache:      make(map[string]*list.Element),
		evictList:  list.New(),
		ttl:        make(map[string]*time.Time),
	}
}

func (c *LRUCache) Set(key string, val any) (ok bool) {
	// Check if the key is already in the cache.
	if elem, ok := c.cache[key]; ok {
		// Update the value and move the element to the front of the eviction list.
		elem.Value.(*entry).value = val
		c.evictList.MoveToFront(elem)
		return true
	}

	// If the cache is full, evict the least recently used item.
	if c.evictList.Len() >= c.capability {
		c.evict()
	}

	// Create a new entry and add it to the cache and eviction list.
	newElem := c.evictList.PushFront(&entry{key: key, value: val})
	c.cache[key] = newElem
	return true
}

func (c *LRUCache) Get(key string) (val any, ok bool) {
	// Check if the key is in the cache.
	if elem, ok := c.cache[key]; ok {
		// Move the accessed element to the front of the eviction list.
		c.evictList.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (c *LRUCache) Expire(key string, expireTime *time.Time) (ok bool) {
	if elem, ok := c.cache[key]; ok {
		c.ttl[key] = expireTime
		c.evictList.MoveToFront(elem)
		return true
	}
	return false
}

func (c *LRUCache) Remove(key string) (ok bool) {
	if elem, ok := c.cache[key]; ok {
		delete(c.cache, key)
		c.evictList.Remove(elem)
		return true
	}
	return false
}

func (c *LRUCache) Purge() {
	// Clear the cache and eviction list.
	c.cache = make(map[string]*list.Element)
	c.evictList.Init()
	c.ttl = make(map[string]*time.Time)
}

func (c *LRUCache) Len() int {
	return c.evictList.Len()
}

func (c *LRUCache) Has(key string) bool {
	_, ok := c.cache[key]
	return ok
}

// ForEach iterates over all key-value pairs in the cache.
func (c *LRUCache) ForEach(iter func(key string, val any, expiration *time.Time) bool) {
	// TODO: Implement the ForEach method.
}

// evict removes the least recently used item from the cache.
func (c *LRUCache) evict() {
	if c.evictList.Len() == 0 {
		return
	}

	// Remove the least recently used item from the cache and eviction list.
	elem := c.evictList.Back()
	if elem != nil {
		c.evictList.Remove(elem)
		entry := elem.Value.(*entry)
		delete(c.cache, entry.key)
	}
}
