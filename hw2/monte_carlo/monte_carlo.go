package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// SafeCounter offers an atomic unsigned int64 counter.
type SafeCounter struct {
	val uint64
	mux sync.Mutex
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

// CountCirclePoints increments a thread-safe counter for
// Monte Carlo circle simulations.
func CountCirclePoints(num uint64, circlePoints *SafeCounter, wg *sync.WaitGroup) {
	defer wg.Done()
	randomSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSource)

	for i := uint64(0); i < num; i++ {
		x := random.Float64()
		y := random.Float64()

		polarDistance := x*x + y*y

		if polarDistance <= 1 {
			circlePoints.Inc()
		}
	}
}

func main() {
	POINTS, _ := strconv.ParseUint(os.Args[1], 10, 64)
	THREADS, _ := strconv.Atoi(os.Args[2])

	var wg sync.WaitGroup
	var circlePoints SafeCounter

	for i := 0; i < THREADS; i++ {
		wg.Add(1)
		go CountCirclePoints(POINTS/uint64(THREADS), &circlePoints, &wg)
	}

	wg.Wait()
	pi := float64(4*circlePoints.Value()) / float64(POINTS)
	fmt.Println(pi)
}
