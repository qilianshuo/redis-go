package list

import (
	"testing"
)

func TestLinkedList_AddAndGet(t *testing.T) {
	list := Make()
	list.Add(1)
	list.Add(2)
	list.Add(3)

	if list.Len() != 3 {
		t.Fatalf("expected length 3, got %d", list.Len())
	}

	if val := list.Get(0); val != 1 {
		t.Fatalf("expected 1, got %v", val)
	}
	if val := list.Get(1); val != 2 {
		t.Fatalf("expected 2, got %v", val)
	}
	if val := list.Get(2); val != 3 {
		t.Fatalf("expected 3, got %v", val)
	}
}

func TestLinkedList_Set(t *testing.T) {
	list := Make(1, 2, 3)
	list.Set(1, 42)

	if val := list.Get(1); val != 42 {
		t.Fatalf("expected 42, got %v", val)
	}
}

func TestLinkedList_Insert(t *testing.T) {
	list := Make(1, 3)
	list.Insert(1, 2)

	if list.Len() != 3 {
		t.Fatalf("expected length 3, got %d", list.Len())
	}

	if val := list.Get(1); val != 2 {
		t.Fatalf("expected 2, got %v", val)
	}
}

func TestLinkedList_Remove(t *testing.T) {
	list := Make(1, 2, 3)
	removed := list.Remove(1)

	if removed != 2 {
		t.Fatalf("expected 2, got %v", removed)
	}

	if list.Len() != 2 {
		t.Fatalf("expected length 2, got %d", list.Len())
	}

	if val := list.Get(1); val != 3 {
		t.Fatalf("expected 3, got %v", val)
	}
}

func TestLinkedList_RemoveLast(t *testing.T) {
	list := Make(1, 2, 3)
	removed := list.RemoveLast()

	if removed != 3 {
		t.Fatalf("expected 3, got %v", removed)
	}

	if list.Len() != 2 {
		t.Fatalf("expected length 2, got %d", list.Len())
	}
}

func TestLinkedList_RemoveAllByVal(t *testing.T) {
	list := Make(1, 2, 1, 3)
	removed := list.RemoveAllByVal(func(val any) bool {
		return val == 1
	})

	if removed != 2 {
		t.Fatalf("expected 2 removals, got %d", removed)
	}

	if list.Len() != 2 {
		t.Fatalf("expected length 2, got %d", list.Len())
	}
}

func TestLinkedList_RemoveByVal(t *testing.T) {
	list := Make(1, 2, 1, 3)
	removed := list.RemoveByVal(func(val any) bool {
		return val == 1
	}, 1)

	if removed != 1 {
		t.Fatalf("expected 1 removal, got %d", removed)
	}

	if list.Len() != 3 {
		t.Fatalf("expected length 3, got %d", list.Len())
	}
}

func TestLinkedList_ReverseRemoveByVal(t *testing.T) {
	list := Make(1, 2, 1, 3)
	removed := list.ReverseRemoveByVal(func(val any) bool {
		return val == 1
	}, 1)

	if removed != 1 {
		t.Fatalf("expected 1 removal, got %d", removed)
	}

	if list.Len() != 3 {
		t.Fatalf("expected length 3, got %d", list.Len())
	}
}

func TestLinkedList_Len(t *testing.T) {
	list := Make(1, 2, 3)

	if list.Len() != 3 {
		t.Fatalf("expected length 3, got %d", list.Len())
	}
}

func TestLinkedList_ForEach(t *testing.T) {
	list := Make(1, 2, 3)
	sum := 0

	list.ForEach(func(i int, val any) bool {
		sum += val.(int)
		return true
	})

	if sum != 6 {
		t.Fatalf("expected sum 6, got %d", sum)
	}
}

func TestLinkedList_Contains(t *testing.T) {
	list := Make(1, 2, 3)

	if !list.Contains(func(val any) bool { return val == 2 }) {
		t.Fatalf("expected to contain 2")
	}

	if list.Contains(func(val any) bool { return val == 42 }) {
		t.Fatalf("expected not to contain 42")
	}
}

func TestLinkedList_Range(t *testing.T) {
	list := Make(1, 2, 3, 4, 5)
	result := list.Range(1, 4)
	expected := []any{2, 3, 4}

	for i, val := range result {
		if val != expected[i] {
			t.Fatalf("expected %v, got %v", expected, result)
		}
	}
}
