package stop

import (
	"sync"
	"testing"
	"time"
)

func TestDoneOpenBeforeStop(t *testing.T) {
	s := New()

	select {
	case <-s.Done():
		t.Fatal("Done() closed before Stop() was called")
	default:
	}
}

func TestStopClosesDone(t *testing.T) {
	s := New()
	s.Stop()

	select {
	case <-s.Done():
	case <-time.After(time.Second):
		t.Fatal("Done() did not unblock after Stop()")
	}
}

func TestStopIsIdempotent(t *testing.T) {
	s := New()
	s.Stop()
	s.Stop()
	s.Stop()

	select {
	case <-s.Done():
	case <-time.After(time.Second):
		t.Fatal("Done() did not unblock after repeated Stop()")
	}
}

func TestConcurrentStop(t *testing.T) {
	s := New()

	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Stop()
		}()
	}
	wg.Wait()

	select {
	case <-s.Done():
	case <-time.After(time.Second):
		t.Fatal("Done() did not unblock after concurrent Stop()")
	}
}
