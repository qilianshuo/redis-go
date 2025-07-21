package dict

import (
	"testing"
)

func TestMakeConcurrent(t *testing.T) {
	dict := MakeConcurrent(4)
	if dict == nil {
		t.Errorf("MakeConcurrent() failed, expected non-nil ConcurrentDict")
	}
}

func TestConcurrentDict_PutAndGet(t *testing.T) {
	dict := MakeConcurrent(4)
	dict.Put("key1", "value1")

	val, exists := dict.Get("key1")
	if !exists || val != "value1" {
		t.Errorf("Get() failed, expected value1, got %v", val)
	}

	_, exists = dict.Get("key2")
	if exists {
		t.Errorf("Get() failed, expected key2 to not exist")
	}
}

func TestConcurrentDict_PutIfAbsent(t *testing.T) {
	dict := MakeConcurrent(4)
	result := dict.PutIfAbsent("key1", "value1")
	if result != 1 {
		t.Errorf("PutIfAbsent() failed, expected 1, got %d", result)
	}

	result = dict.PutIfAbsent("key1", "value2")
	if result != 0 {
		t.Errorf("PutIfAbsent() failed, expected 0, got %d", result)
	}
}

func TestConcurrentDict_PutIfExists(t *testing.T) {
	dict := MakeConcurrent(4)
	result := dict.PutIfExists("key1", "value1")
	if result != 0 {
		t.Errorf("PutIfExists() failed, expected 0, got %d", result)
	}

	dict.Put("key1", "value1")
	result = dict.PutIfExists("key1", "value2")
	if result != 1 {
		t.Errorf("PutIfExists() failed, expected 1, got %d", result)
	}
}

func TestConcurrentDict_Remove(t *testing.T) {
	dict := MakeConcurrent(4)
	dict.Put("key1", "value1")

	val, result := dict.Remove("key1")
	if result != 1 || val != "value1" {
		t.Errorf("Remove() failed, expected value1 and 1, got %v and %d", val, result)
	}

	_, result = dict.Remove("key2")
	if result != 0 {
		t.Errorf("Remove() failed, expected 0, got %d", result)
	}
}

func TestConcurrentDict_Len(t *testing.T) {
	dict := MakeConcurrent(4)
	if dict.Len() != 0 {
		t.Errorf("Len() failed, expected 0, got %d", dict.Len())
	}

	dict.Put("key1", "value1")
	if dict.Len() != 1 {
		t.Errorf("Len() failed, expected 1, got %d", dict.Len())
	}
}

func TestConcurrentDict_Keys(t *testing.T) {
	dict := MakeConcurrent(4)
	dict.Put("key1", "value1")
	dict.Put("key2", "value2")

	keys := dict.Keys()
	if len(keys) != 2 {
		t.Errorf("Keys() failed, expected 2 keys, got %d", len(keys))
	}
}

func TestConcurrentDict_ForEach(t *testing.T) {
	dict := MakeConcurrent(4)
	dict.Put("key1", "value1")
	dict.Put("key2", "value2")

	count := 0
	dict.ForEach(func(key string, val interface{}) bool {
		count++
		return true
	})

	if count != 2 {
		t.Errorf("ForEach() failed, expected 2 iterations, got %d", count)
	}
}

func TestConcurrentDict_RandomKeys(t *testing.T) {
	dict := MakeConcurrent(4)
	dict.Put("key1", "value1")
	dict.Put("key2", "value2")

	keys := dict.RandomKeys(1)
	if len(keys) != 1 {
		t.Errorf("RandomKeys() failed, expected 1 key, got %d", len(keys))
	}
}

func TestConcurrentDict_RandomDistinctKeys(t *testing.T) {
	dict := MakeConcurrent(4)
	dict.Put("key1", "value1")
	dict.Put("key2", "value2")

	keys := dict.RandomDistinctKeys(2)
	if len(keys) != 2 {
		t.Errorf("RandomDistinctKeys() failed, expected 2 keys, got %d", len(keys))
	}
}

func TestConcurrentDict_Clear(t *testing.T) {
	dict := MakeConcurrent(4)
	dict.Put("key1", "value1")
	dict.Clear()

	if dict.Len() != 0 {
		t.Errorf("Clear() failed, expected 0, got %d", dict.Len())
	}
}
