package main

import (
	"flag"
	"fmt"
    "mpcs52060/justincohler/hw3/src/ppsync"
	"runtime"
	"sync"
)

// SafeCounter offers an atomic unsigned int64 counter.
type SafeCounter struct {
	val    uint64
	locker sync.Locker
}

// Inc thread-safely increments the counter.
func (c *SafeCounter) Inc() {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.val++
}

// Value returns the current value of the counter.
func (c *SafeCounter) Value() uint64 {
	c.locker.Lock()
	defer c.locker.Unlock()
	return c.val
}

// IncNum is a threaded counter.
func IncNum(wg *sync.WaitGroup, c *SafeCounter, num int) {
	defer wg.Done()
	for i := 0; i < num; i++ {
		c.Inc()
	}
}

func main() {
	// Read Flags
	tasAction := flag.Bool("tas", false, "TAS Flag")
	ttasAction := flag.Bool("ttas", false, "TTAS Flag")
	ebAction := flag.Bool("eb", false, "EB Flag")
	aAction := flag.Bool("a", false, "A Flag")

	flag.Parse()

	var wg sync.WaitGroup
	var locker sync.Locker

	THREADS := runtime.NumCPU()

	switch {
	case *tasAction:
		locker = &ppsync.TASLock{}
	case *ttasAction:
		locker = &ppsync.TTASLock{}
	case *ebAction:
		locker = &ppsync.EBLock{}
	case *aAction:
		locker = ppsync.NewALock(THREADS)
	}

	c := SafeCounter{locker: locker}

	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go IncNum(&wg, &c, 1e6/THREADS)
	}
	wg.Wait()
	fmt.Println(c.val)
}
