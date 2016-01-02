package keylock

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestKeyLock(t *testing.T) {
	counters := make(map[string]int)
	N := 1000000
	NGOROUTINES := 10
	KEY := "hello"
	var wait sync.WaitGroup
	wait.Add(NGOROUTINES)
	keylock := NewKeyLock()

	adder := func(key string, n int) {
		for i := 0; i < n; i++ {
			keylock.Lock(key)
			counters[key]++
			keylock.Unlock(key)
		}
		wait.Done()
	}

	for i := 0; i < NGOROUTINES; i++ {
		go adder(KEY, N)
	}
	wait.Wait()
	if counters[KEY] != NGOROUTINES*N {
		t.Errorf("Counter is %d, but should be %d", counters[KEY], NGOROUTINES*N)
	}
}

func TestKeyRWLock(t *testing.T) {
	counters := make(map[string]int)
	N := 100000
	NGOROUTINES := 3
	KEY := "hello"
	var wait sync.WaitGroup
	wait.Add(NGOROUTINES + 1) // 1 for the reader
	keylock := NewRWKeyLock()

	var correctCounter int64

	adder := func(key string, n int) {
		for i := 0; i < n; i++ {
			keylock.Lock(key)
			counters[key]++
			atomic.AddInt64(&correctCounter, 1)
			keylock.Unlock(key)
		}
		wait.Done()
	}

	reader := func(key string, n int) {
		for i := 0; i < n; i++ {
			keylock.RLock(key)
			c := counters[key]
			if c != int(correctCounter) {
				t.Errorf("Counter is %d, but should be %d", c, correctCounter)
			}
			keylock.RUnlock(key)
		}
		wait.Done()
	}
	for i := 0; i < NGOROUTINES; i++ {
		go adder(KEY, N)
	}
	go reader(KEY, N)

	wait.Wait()
	if counters[KEY] != NGOROUTINES*N {
		t.Errorf("Counter is %d, but should be %d", counters[KEY], NGOROUTINES*N)
	}
}
