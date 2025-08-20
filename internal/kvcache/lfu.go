package kvcache

import (
	"container/list"
	"time"
)

// LFUCache is a cache that evicts the least frequently used items first.
type LFUCache struct {
	capability int
	cache      map[string]*list.Element
	ttl        map[string]*time.Time
	freqList   *list.List
}

// lfuEntry represents a single entry in the cache.
type lfuEntry struct {
	key   string
	value any
	freq  int
}

// NewLFUCache creates a new LFUCache.
func NewLFUCache(capability int) *LFUCache {
	return &LFUCache{
		capability: capability,
		cache:      make(map[string]*list.Element),
		ttl:        make(map[string]*time.Time),
		freqList:   list.New(),
	}
}

func (c *LFUCache) Set(key string, val any) (ok bool) {
	if elem, ok := c.cache[key]; ok {
		// Update the value and frequency.
		elem.Value.(*lfuEntry).value = val
		elem.Value.(*lfuEntry).freq++
		return true
	}

	// If the cache is full, evict the least frequently used item.
	if len(c.cache) >= c.capability {
		c.evict()
	}

	// Create a new entry and add it to the cache.
	newElem := c.freqList.PushFront(&lfuEntry{key: key, value: val, freq: 1})
	c.cache[key] = newElem
	return true
}

func (c *LFUCache) Get(key string) (val any, ok bool) {
	if elem, ok := c.cache[key]; ok {
		// Move the accessed element to the front of the frequency list.
		c.freqList.MoveToFront(elem)
		return elem.Value.(*lfuEntry).value, true
	}
	return nil, false
}

func (c *LFUCache) Expire(key string, expireTime *time.Time) (ok bool) {
	if elem, ok := c.cache[key]; ok {
		c.ttl[key] = expireTime
		c.freqList.MoveToFront(elem)
		return true
	}
	return false
}

func (c *LFUCache) Remove(key string) (ok bool) {
	if elem, ok := c.cache[key]; ok {
		delete(c.cache, key)
		delete(c.ttl, key)
		c.freqList.Remove(elem)
		return true
	}
	return false
}

func (c *LFUCache) Purge() {
	// Clear the cache and frequency list.
	c.cache = make(map[string]*list.Element)
	c.freqList.Init()
	c.ttl = make(map[string]*time.Time)
}

func (c *LFUCache) Len() int {
	return c.freqList.Len()
}

func (c *LFUCache) Has(key string) bool {
	_, ok := c.cache[key]
	return ok
}

func (c *LFUCache) ForEach(iter func(key string, val any, expiration *time.Time) bool) {
	for key, elem := range c.cache {
		if !iter(key, elem.Value.(*lfuEntry).value, c.ttl[key]) {
			break
		}
	}
}

func (c *LFUCache) evict() {
	var lfuKey string
	var lfuFreq int
	for key, elem := range c.cache {
		if lfuFreq == 0 || elem.Value.(*lfuEntry).freq < lfuFreq {
			lfuFreq = elem.Value.(*lfuEntry).freq
			lfuKey = key
		}
	}
	if lfuKey != "" {
		delete(c.cache, lfuKey)
		delete(c.ttl, lfuKey)
	}
}
