package kvcache

import (
	"testing"
	"time"
)

func TestLFUCacheSetGet(t *testing.T) {
	cache := NewLFUCache(2)
	ok := cache.Set("a", 1)
	if !ok {
		t.Errorf("Set failed")
	}
	v, ok := cache.Get("a")
	if !ok || v != 1 {
		t.Errorf("Get failed, got %v", v)
	}
}

func TestLFUCacheEvict(t *testing.T) {
	cache := NewLFUCache(2)
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("b", 3) // increase freq of b
	cache.Set("c", 4) // should evict a
	if cache.Has("a") {
		t.Errorf("Eviction failed, 'a' should be evicted")
	}
	if !cache.Has("b") || !cache.Has("c") {
		t.Errorf("Eviction failed, missing keys")
	}
}

func TestLFUCacheRemove(t *testing.T) {
	cache := NewLFUCache(2)
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

func TestLFUCacheExpire(t *testing.T) {
	cache := NewLFUCache(2)
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

func TestLFUCachePurge(t *testing.T) {
	cache := NewLFUCache(2)
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

func TestLFUCacheLen(t *testing.T) {
	cache := NewLFUCache(2)
	if cache.Len() != 0 {
		t.Errorf("Len should be 0 initially")
	}
	cache.Set("a", 1)
	if cache.Len() != 1 {
		t.Errorf("Len should be 1 after Set")
	}
}
