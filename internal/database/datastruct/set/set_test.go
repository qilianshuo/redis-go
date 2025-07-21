package set

import (
	"testing"
)

func TestMakeConcurrentSafe(t *testing.T) {
	set := NewConcurrentSet("a", "b", "c")
	if set.Len() != 3 {
		t.Errorf("MakeConcurrentSafe() failed, expected 3, got %d", set.Len())
	}
}

func TestSet_Add(t *testing.T) {
	set := NewConcurrentSet()
	result := set.Add("a")
	if result != 1 || !set.Has("a") {
		t.Errorf("Add() failed, expected 1 and 'a' to exist")
	}

	result = set.Add("a")
	if result != 0 {
		t.Errorf("Add() failed, expected 0 for duplicate add")
	}
}

func TestSet_Remove(t *testing.T) {
	set := NewConcurrentSet("a", "b")
	result := set.Remove("a")
	if result != 1 || set.Has("a") {
		t.Errorf("Remove() failed, expected 1 and 'a' to not exist")
	}

	result = set.Remove("c")
	if result != 0 {
		t.Errorf("Remove() failed, expected 0 for non-existent element")
	}
}

func TestSet_Has(t *testing.T) {
	set := NewConcurrentSet("a", "b")
	if !set.Has("a") {
		t.Errorf("Has() failed, expected 'a' to exist")
	}

	if set.Has("c") {
		t.Errorf("Has() failed, expected 'c' to not exist")
	}
}

func TestSet_Len(t *testing.T) {
	set := NewConcurrentSet("a", "b")
	if set.Len() != 2 {
		t.Errorf("Len() failed, expected 2, got %d", set.Len())
	}
}

func TestSet_ToSlice(t *testing.T) {
	set := NewConcurrentSet("a", "b")
	slice := set.ToSlice()
	if len(slice) != 2 {
		t.Errorf("ToSlice() failed, expected 2 elements, got %d", len(slice))
	}
}

func TestSet_ForEach(t *testing.T) {
	set := NewConcurrentSet("a", "b")
	count := 0
	set.ForEach(func(member string) bool {
		count++
		return true
	})
	if count != 2 {
		t.Errorf("ForEach() failed, expected 2 iterations, got %d", count)
	}
}

func TestSet_ShallowCopy(t *testing.T) {
	set := NewConcurrentSet("a", "b")
	setCopy := set.ShallowCopy()
	if setCopy.Len() != set.Len() || !setCopy.Has("a") || !setCopy.Has("b") {
		t.Errorf("ShallowCopy() failed, expected identical set")
	}
}

func TestIntersect(t *testing.T) {
	set1 := NewConcurrentSet("a", "b")
	set2 := NewConcurrentSet("b", "c")
	result := Intersect(set1, set2)
	if result.Len() != 1 || !result.Has("b") {
		t.Errorf("Intersect() failed, expected set with 'b'")
	}
}

func TestUnion(t *testing.T) {
	set1 := NewConcurrentSet("a", "b")
	set2 := NewConcurrentSet("b", "c")
	result := Union(set1, set2)
	if result.Len() != 3 || !result.Has("a") || !result.Has("b") || !result.Has("c") {
		t.Errorf("Union() failed, expected set with 'a', 'b', 'c'")
	}
}

func TestSet_RandomDistinctMembers(t *testing.T) {
	set := NewConcurrentSet("a", "b", "c")
	members := set.RandomDistinctMembers(2)
	if len(members) != 2 {
		t.Errorf("RandomDistinctMembers() failed, expected 2 distinct members, got %d", len(members))
	}
}
