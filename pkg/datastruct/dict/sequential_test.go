package dict

import (
	"testing"
)

func TestMakeSimple(t *testing.T) {
	dict := NewSequentialDict()
	if dict == nil || dict.m == nil {
		t.Errorf("MakeSimple() failed, expected non-nil SimpleDict")
	}
}

func TestSimpleDict_Get(t *testing.T) {
	dict := NewSequentialDict()
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

func TestSimpleDict_Put(t *testing.T) {
	dict := NewSequentialDict()
	result := dict.Put("key1", "value1")
	if result != 1 {
		t.Errorf("Put() failed, expected 1, got %d", result)
	}

	result = dict.Put("key1", "value2")
	if result != 0 {
		t.Errorf("Put() failed, expected 0, got %d", result)
	}
}

func TestSimpleDict_PutIfAbsent(t *testing.T) {
	dict := NewSequentialDict()
	result := dict.PutIfAbsent("key1", "value1")
	if result != 1 {
		t.Errorf("PutIfAbsent() failed, expected 1, got %d", result)
	}

	result = dict.PutIfAbsent("key1", "value2")
	if result != 0 {
		t.Errorf("PutIfAbsent() failed, expected 0, got %d", result)
	}
}

func TestSimpleDict_PutIfExists(t *testing.T) {
	dict := NewSequentialDict()
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

func TestSimpleDict_Remove(t *testing.T) {
	dict := NewSequentialDict()
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

func TestSimpleDict_Len(t *testing.T) {
	dict := NewSequentialDict()
	if dict.Len() != 0 {
		t.Errorf("Len() failed, expected 0, got %d", dict.Len())
	}

	dict.Put("key1", "value1")
	if dict.Len() != 1 {
		t.Errorf("Len() failed, expected 1, got %d", dict.Len())
	}
}

func TestSimpleDict_Keys(t *testing.T) {
	dict := NewSequentialDict()
	dict.Put("key1", "value1")
	dict.Put("key2", "value2")

	keys := dict.Keys()
	if len(keys) != 2 {
		t.Errorf("Keys() failed, expected 2 keys, got %d", len(keys))
	}
}

func TestSimpleDict_ForEach(t *testing.T) {
	dict := NewSequentialDict()
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

func TestSimpleDict_RandomKeys(t *testing.T) {
	dict := NewSequentialDict()
	dict.Put("key1", "value1")
	dict.Put("key2", "value2")

	keys := dict.RandomKeys(1)
	if len(keys) != 1 {
		t.Errorf("RandomKeys() failed, expected 1 key, got %d", len(keys))
	}
}

func TestSimpleDict_RandomDistinctKeys(t *testing.T) {
	dict := NewSequentialDict()
	dict.Put("key1", "value1")
	dict.Put("key2", "value2")

	keys := dict.RandomDistinctKeys(2)
	if len(keys) != 2 {
		t.Errorf("RandomDistinctKeys() failed, expected 2 keys, got %d", len(keys))
	}
}

func TestSimpleDict_Clear(t *testing.T) {
	dict := NewSequentialDict()
	dict.Put("key1", "value1")
	dict.Clear()

	if dict.Len() != 0 {
		t.Errorf("Clear() failed, expected 0, got %d", dict.Len())
	}
}
