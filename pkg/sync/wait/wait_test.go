package wait

import (
	"testing"
	"time"
)

func TestWaitBasic(t *testing.T) {
	var w Wait
	w.Add(1)
	done := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Millisecond)
		w.Done()
		close(done)
	}()
	w.Wait()
	<-done
}

func TestWaitWithTimeoutTimeout(t *testing.T) {
	var w Wait
	w.Add(1)
	timedOut := w.WaitWithTimeout(20 * time.Millisecond)
	if !timedOut {
		t.Errorf("Expected timeout, got false")
	}
	w.Done()
}

func TestWaitWithTimeoutNoTimeout(t *testing.T) {
	var w Wait
	w.Add(1)
	go func() {
		time.Sleep(10 * time.Millisecond)
		w.Done()
	}()
	timedOut := w.WaitWithTimeout(100 * time.Millisecond)
	if timedOut {
		t.Errorf("Expected no timeout, got true")
	}
}
