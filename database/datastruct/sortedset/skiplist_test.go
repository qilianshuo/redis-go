package sortedset

import (
	"testing"
)

func TestSkipList_InsertAndGet(t *testing.T) {
	sl := makeSkipList()
	sl.insert("a", 1.0)
	sl.insert("b", 2.0)
	sl.insert("c", 3.0)

	if sl.length != 3 {
		t.Fatalf("expected length 3, got %d", sl.length)
	}

	node := sl.getByRank(1)
	if node == nil || node.Member != "a" || node.Score != 1.0 {
		t.Fatalf("expected node with member 'a' and score 1.0, got %+v", node)
	}

	node = sl.getByRank(2)
	if node == nil || node.Member != "b" || node.Score != 2.0 {
		t.Fatalf("expected node with member 'b' and score 2.0, got %+v", node)
	}
}

func TestSkipList_Remove(t *testing.T) {
	sl := makeSkipList()
	sl.insert("a", 1.0)
	sl.insert("b", 2.0)
	sl.insert("c", 3.0)

	removed := sl.remove("b", 2.0)
	if !removed {
		t.Fatalf("expected to remove member 'b'")
	}

	if sl.length != 2 {
		t.Fatalf("expected length 2, got %d", sl.length)
	}

	node := sl.getByRank(2)
	if node == nil || node.Member != "c" {
		t.Fatalf("expected node with member 'c', got %+v", node)
	}
}

func TestSkipList_GetRank(t *testing.T) {
	sl := makeSkipList()
	sl.insert("a", 1.0)
	sl.insert("b", 2.0)
	sl.insert("c", 3.0)

	rank := sl.getRank("b", 2.0)
	if rank != 2 {
		t.Fatalf("expected rank 2, got %d", rank)
	}

	rank = sl.getRank("d", 4.0)
	if rank != 0 {
		t.Fatalf("expected rank 0 for non-existent member, got %d", rank)
	}
}

func TestSkipList_RemoveRange(t *testing.T) {
	sl := makeSkipList()
	sl.insert("a", 1.0)
	sl.insert("b", 2.0)
	sl.insert("c", 3.0)
	sl.insert("d", 4.0)

	minScore, _ := ParseScoreBorder("2.0")
	maxScore, _ := ParseScoreBorder("(4.0")

	removed := sl.RemoveRange(minScore, maxScore, 0)
	if len(removed) != 2 {
		t.Fatalf("expected 2 elements removed, got %d", len(removed))
	}

	if sl.length != 2 {
		t.Fatalf("expected length 2, got %d", sl.length)
	}
}

func TestSkipList_RemoveRangeByRank(t *testing.T) {
	sl := makeSkipList()
	sl.insert("a", 1.0)
	sl.insert("b", 2.0)
	sl.insert("c", 3.0)
	sl.insert("d", 4.0)

	removed := sl.RemoveRangeByRank(2, 4)
	if len(removed) != 2 {
		t.Fatalf("expected 2 elements removed, got %d", len(removed))
	}

	if sl.length != 2 {
		t.Fatalf("expected length 2, got %d", sl.length)
	}
}

func BenchmarkSkipList_Insert(b *testing.B) {
	sl := makeSkipList()
	for i := 0; i < b.N; i++ {
		sl.insert(string(rune(i)), float64(i))
	}
}

func BenchmarkSkipList_Remove(b *testing.B) {
	sl := makeSkipList()
	for i := 0; i < 1000; i++ {
		sl.insert(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.remove(string(rune(i%1000)), float64(i%1000))
	}
}

func BenchmarkSkipList_GetRank(b *testing.B) {
	sl := makeSkipList()
	for i := 0; i < 1000; i++ {
		sl.insert(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.getRank(string(rune(i%1000)), float64(i%1000))
	}
}

func BenchmarkSkipList_GetByRank(b *testing.B) {
	sl := makeSkipList()
	for i := 0; i < 1000; i++ {
		sl.insert(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.getByRank(int64(i%1000 + 1))
	}
}

func BenchmarkSkipList_RemoveRange(b *testing.B) {
	sl := makeSkipList()
	for i := 0; i < 1000; i++ {
		sl.insert(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		minScore, _ := ParseScoreBorder("2.0")
		maxScore, _ := ParseScoreBorder("(4.0")
		sl.RemoveRange(minScore, maxScore, 0)
	}
}

func BenchmarkSkipList_RemoveRangeByRank(b *testing.B) {
	sl := makeSkipList()
	for i := 0; i < 1000; i++ {
		sl.insert(string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.RemoveRangeByRank(100, 900)
	}
}
