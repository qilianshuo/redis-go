package kvcache

import (
	"testing"
	"time"
)

func TestLRUCacheSetGet(t *testing.T) {
	cache := NewLRUCache(2)
	ok := cache.Set("a", 1)
	if !ok {
		t.Errorf("Set failed")
	}
	v, ok := cache.Get("a")
	if !ok || v != 1 {
		t.Errorf("Get failed, got %v", v)
	}
}

func TestLRUCacheEvict(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3) // should evict "a"
	if cache.Has("a") {
		t.Errorf("Eviction failed, 'a' should be evicted")
	}
	if !cache.Has("b") || !cache.Has("c") {
		t.Errorf("Eviction failed, missing keys")
	}
}

func TestLRUCacheRemove(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("a", 1)
	cache.Set("b", 2)
	ok := cache.Remove("a")
	if !ok || cache.Has("a") {
		t.Errorf("Remove failed")
	}
	ok = cache.Remove("notfound")
	if ok {
		t.Errorf("Remove should fail for missing key")
	}
}

func TestLRUCacheExpire(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("a", 1)
	expire := time.Now().Add(time.Hour)
	ok := cache.Expire("a", &expire)
	if !ok {
		t.Errorf("Expire failed")
	}
	ok = cache.Expire("notfound", &expire)
	if ok {
		t.Errorf("Expire should fail for missing key")
	}
}

func TestLRUCachePurge(t *testing.T) {
	cache := NewLRUCache(2)
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Purge()
	if cache.Len() != 0 {
		t.Errorf("Purge failed")
	}
	if cache.Has("a") || cache.Has("b") {
		t.Errorf("Purge did not clear keys")
	}
}

func TestLRUCacheLen(t *testing.T) {
	cache := NewLRUCache(2)
	if cache.Len() != 0 {
		t.Errorf("Len should be 0 initially")
	}
	cache.Set("a", 1)
	if cache.Len() != 1 {
		t.Errorf("Len should be 1 after Set")
	}
}
