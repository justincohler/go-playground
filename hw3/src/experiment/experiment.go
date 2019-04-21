package main

import (
	"flag"
	"fmt"
	"mpcs52060/justincohler/hw3/src/ppsync"
	"sync"
)

// SafeCounter offers an atomic unsigned int64 counter.
type SafeCounter struct {
	val uint64
	mux ppsync.Lock
}

// Inc thread-safely increments the counter.
func (c *SafeCounter) Inc() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.val++
}

// Value returns the current value of the counter.
func (c *SafeCounter) Value() uint64 {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.val
}

// IncNum is a threaded counter.
func IncNum(wg *sync.WaitGroup, c *SafeCounter, num int) {
	defer wg.Done()
	for i := 0; i < 10e6; i++ {
		c.Inc()
	}
}

func main() {
	// Read Flags
	tasAction := flag.Bool("tas", false, "TAS Flag")
	ttasAction := flag.Bool("ttas", false, "TTAS Flag")
	// ebAction := flag.Bool("eb", false, "EB Flag")
	// aAction := flag.Bool("a", false, "A Flag")

	var wg sync.WaitGroup
	var lock ppsync.Lock

	switch {
	case *tasAction == true:
		lock = ppsync.TASLock{}
	case *ttasAction == true:
		lock = ppsync.TTASLock{}
	}

	c := SafeCounter{mux: lock}

	THREADS := 4
	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go IncNum(&wg, &c, 10e6/THREADS)
	}
	wg.Wait()
	fmt.Println(c.val)
}
