package list

import (
	"testing"
)

func TestQuickList_AddAndGet(t *testing.T) {
	ql := NewQuickList()
	ql.Add(1)
	ql.Add(2)
	ql.Add(3)

	if ql.Len() != 3 {
		t.Fatalf("expected length 3, got %d", ql.Len())
	}

	if val := ql.Get(0); val != 1 {
		t.Fatalf("expected 1, got %v", val)
	}
	if val := ql.Get(1); val != 2 {
		t.Fatalf("expected 2, got %v", val)
	}
	if val := ql.Get(2); val != 3 {
		t.Fatalf("expected 3, got %v", val)
	}
}

func TestQuickList_Set(t *testing.T) {
	ql := NewQuickList()
	ql.Add(1)
	ql.Add(2)
	ql.Add(3)

	ql.Set(1, 42)
	if val := ql.Get(1); val != 42 {
		t.Fatalf("expected 42, got %v", val)
	}
}

func TestQuickList_Insert(t *testing.T) {
	ql := NewQuickList()
	ql.Add(1)
	ql.Add(3)

	ql.Insert(1, 2)
	if ql.Len() != 3 {
		t.Fatalf("expected length 3, got %d", ql.Len())
	}

	if val := ql.Get(1); val != 2 {
		t.Fatalf("expected 2, got %v", val)
	}
}

func TestQuickList_Remove(t *testing.T) {
	ql := NewQuickList()
	ql.Add(1)
	ql.Add(2)
	ql.Add(3)

	removed := ql.Remove(1)
	if removed != 2 {
		t.Fatalf("expected 2, got %v", removed)
	}

	if ql.Len() != 2 {
		t.Fatalf("expected length 2, got %d", ql.Len())
	}

	if val := ql.Get(1); val != 3 {
		t.Fatalf("expected 3, got %v", val)
	}
}

func TestQuickList_Range(t *testing.T) {
	ql := NewQuickList()
	for i := 1; i <= 5; i++ {
		ql.Add(i)
	}

	result := ql.Range(1, 4)
	expected := []interface{}{2, 3, 4}

	for i, val := range result {
		if val != expected[i] {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	}
}

func TestQuickList_RemoveAllByVal(t *testing.T) {
	ql := NewQuickList()
	ql.Add(1)
	ql.Add(2)
	ql.Add(1)
	ql.Add(3)

	removed := ql.RemoveAllByVal(func(val interface{}) bool {
		return val == 1
	})

	if removed != 2 {
		t.Fatalf("expected 2 removals, got %d", removed)
	}

	if ql.Len() != 2 {
		t.Fatalf("expected length 2, got %d", ql.Len())
	}
}

func TestQuickList_Contains(t *testing.T) {
	ql := NewQuickList()
	ql.Add(1)
	ql.Add(2)
	ql.Add(3)

	if !ql.Contains(func(val interface{}) bool { return val == 2 }) {
		t.Fatalf("expected to contain 2")
	}

	if ql.Contains(func(val interface{}) bool { return val == 42 }) {
		t.Fatalf("expected not to contain 42")
	}
}
