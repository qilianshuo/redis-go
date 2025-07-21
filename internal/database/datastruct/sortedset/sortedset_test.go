package sortedset

import (
	"testing"
)

func TestSortedSet_AddAndGet(t *testing.T) {
	ss := Make()
	ss.Add("a", 1.0)
	ss.Add("b", 2.0)
	ss.Add("c", 3.0)

	if ss.Len() != 3 {
		t.Fatalf("expected length 3, got %d", ss.Len())
	}

	element, ok := ss.Get("b")
	if !ok || element.Score != 2.0 {
		t.Fatalf("expected element with score 2.0, got %+v", element)
	}
}

func TestSortedSet_Remove(t *testing.T) {
	ss := Make()
	ss.Add("a", 1.0)
	ss.Add("b", 2.0)
	ss.Add("c", 3.0)

	removed := ss.Remove("b")
	if !removed {
		t.Fatalf("expected to remove member 'b'")
	}

	if ss.Len() != 2 {
		t.Fatalf("expected length 2, got %d", ss.Len())
	}

	_, ok := ss.Get("b")
	if ok {
		t.Fatalf("expected member 'b' to be removed")
	}
}

func TestSortedSet_GetRank(t *testing.T) {
	ss := Make()
	ss.Add("a", 1.0)
	ss.Add("b", 2.0)
	ss.Add("c", 3.0)

	rank := ss.GetRank("b", false)
	if rank != 1 {
		t.Fatalf("expected rank 1, got %d", rank)
	}

	rank = ss.GetRank("b", true)
	if rank != 1 {
		t.Fatalf("expected rank 1 in descending order, got %d", rank)
	}
}

func TestSortedSet_RangeByRank(t *testing.T) {
	ss := Make()
	ss.Add("a", 1.0)
	ss.Add("b", 2.0)
	ss.Add("c", 3.0)

	result := ss.RangeByRank(0, 2, false)
	if len(result) != 2 || result[0].Member != "a" || result[1].Member != "b" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSortedSet_Range(t *testing.T) {
	ss := Make()
	ss.Add("a", 1.0)
	ss.Add("b", 2.0)
	ss.Add("c", 3.0)

	minBorder, _ := ParseScoreBorder("1.5")
	maxBorder, _ := ParseScoreBorder("3.0")
	result := ss.Range(minBorder, maxBorder, 0, -1, false)
	if len(result) != 2 || result[0].Member != "b" || result[1].Member != "c" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestSortedSet_RemoveRange(t *testing.T) {
	ss := Make()
	ss.Add("a", 1.0)
	ss.Add("b", 2.0)
	ss.Add("c", 3.0)

	minBorder, _ := ParseScoreBorder("1.5")
	maxBorder, _ := ParseScoreBorder("3.0")
	removed := ss.RemoveRange(minBorder, maxBorder)
	if removed != 2 {
		t.Fatalf("expected 2 elements removed, got %d", removed)
	}

	if ss.Len() != 1 {
		t.Fatalf("expected length 1, got %d", ss.Len())
	}
}

func BenchmarkSortedSet_Add(b *testing.B) {
	ss := Make()
	for i := 0; i < b.N; i++ {
		ss.Add(string(rune(i)), float64(i))
	}
}

func BenchmarkSortedSet_Remove(b *testing.B) {
	ss := Make()
	for i := 0; i < 1000; i++ {
		ss.Add(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ss.Remove(string(rune(i % 1000)))
	}
}

func BenchmarkSortedSet_GetRank(b *testing.B) {
	ss := Make()
	for i := 0; i < 1000; i++ {
		ss.Add(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ss.GetRank(string(rune(i%1000)), false)
	}
}

func BenchmarkSortedSet_RangeByRank(b *testing.B) {
	ss := Make()
	for i := 0; i < 1000; i++ {
		ss.Add(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ss.RangeByRank(100, 900, false)
	}
}

func BenchmarkSortedSet_RemoveRange(b *testing.B) {
	ss := Make()
	for i := 0; i < 1000; i++ {
		ss.Add(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		minBorder, _ := ParseScoreBorder("100.0")
		maxBorder, _ := ParseScoreBorder("900.0")
		ss.RemoveRange(minBorder, maxBorder)
	}
}
