package keylock

import (
	"flag"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	log.Println("NumCPU:", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	os.Exit(m.Run())
}

func TestKeyLock(t *testing.T) {
	log.Println("TestKeyLock")
	var accesses int32
	N := 1000000
	NGOROUTINES := 10
	KEYS := []string{"hello", "world", "foo", "bar", "fiz", "buz"}
	var wait sync.WaitGroup
	wait.Add(NGOROUTINES * len(KEYS))
	keylock := NewKeyLock()

	adder := func(key string, n int) {
		for i := 0; i < n; i++ {
			keylock.Lock(key)
			atomic.AddInt32(&accesses, 1)
			keylock.Unlock(key)
		}
		wait.Done()
	}

	for i := 0; i < NGOROUTINES; i++ {
		for _, key := range KEYS {
			go adder(key, N)
		}
	}
	wait.Wait()
	for _, key := range KEYS {
		if accesses != int32(NGOROUTINES*N*len(KEYS)) {
			t.Errorf("[%s] Counter is %d, but should be %d", key, accesses, NGOROUTINES*N)
		}
	}
}

func TestKeyRWLock(t *testing.T) {
	log.Println("TestKeyRWLock")
	counters := make(map[string]*int)
	N := 100000
	NGOROUTINES := 3
	KEYS := []string{"hello", "world", "foo", "bar", "fiz", "buz"}
	var wait sync.WaitGroup
	keylock := NewKeyRWLock()

	// we can count use a buffered channel as an easy way to count how many times we've 'accessed' a key
	correctCounters := make(map[string]chan bool)
	for _, key := range KEYS {
		correctCounters[key] = make(chan bool, N*NGOROUTINES)
		zero := 0
		counters[key] = &zero
	}

	adder := func(key string, n int) {
		for i := 0; i < n; i++ {
			keylock.Lock(key)
			(*counters[key])++
			correctCounters[key] <- true
			keylock.Unlock(key)
		}
		wait.Done()
	}

	reader := func(key string, n int) {
		for i := 0; i < n; i++ {
			keylock.RLock(key)
			c := (*counters[key])
			if c != len(correctCounters[key]) {
				t.Errorf("Counter is %d, but should be %d", c, len(correctCounters[key]))
			}
			keylock.RUnlock(key)
		}
		wait.Done()
	}
	for i := 0; i < NGOROUTINES; i++ {
		for _, key := range KEYS {
			wait.Add(1)
			go adder(key, N)
		}
	}

	for _, key := range KEYS {
		wait.Add(1)
		go reader(key, N)
	}

	wait.Wait()

	for _, key := range KEYS {
		if (*counters[key]) != NGOROUTINES*N {
			t.Errorf("Counter is %d, but should be %d", (*counters[key]), NGOROUTINES*N)
		}
	}
}
