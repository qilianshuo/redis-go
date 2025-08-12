package atomic

import (
	"sync"
	"testing"
)

func TestBooleanBasic(t *testing.T) {
	var b Boolean
	if b.Get() {
		t.Errorf("Default value should be false")
	}
	b.Set(true)
	if !b.Get() {
		t.Errorf("Set(true) failed")
	}
	b.Set(false)
	if b.Get() {
		t.Errorf("Set(false) failed")
	}
}

func TestBooleanConcurrent(t *testing.T) {
	var b Boolean
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			b.Set(i%2 == 0)
		}(i)
	}
	wg.Wait()
	// 最终值不确定，但应无竞态
	_ = b.Get()
}
