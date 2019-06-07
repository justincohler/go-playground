// Attribution for the initial body math & structure goes to:
//
// http://www.cyber-omelette.com/2016/11/python-n-body-orbital-simulation.html#theprogram
//
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"physics"
	"runtime"
)

func blockUpdateLocations(start int, end int, bodies []*physics.Body, bodyCounter chan<- int, nDaysPerStep int) {
	for j := start; j < end; j++ {
		bodies[j].UpdateLocation(bodies, nDaysPerStep)
		bodyCounter <- j
	}
}

func simulate(bodies []*physics.Body, nThreads int, nDaysPerStep int, done <-chan interface{}) <-chan interface{} {
	N := len(bodies)
	stepDone := make(chan interface{})
	bodyCounter := make(chan int, N)

	parallelismFactor := int(math.Min(float64(nThreads), float64(N)))

	go func() {
		for {
			for i := 0; i < parallelismFactor; i++ {

				blockSize := math.Ceil(float64(N) / float64(parallelismFactor))
				start := int(float64(i) * blockSize)
				end := int(math.Min(float64(N), (float64(i)+1)*blockSize))

				go blockUpdateLocations(start, end, bodies, bodyCounter, nDaysPerStep) // functional parallelism
			}

			for i := 0; i < parallelismFactor; i++ { // wait for step to complete
				<-bodyCounter
			}
			select {
			case <-done:
				return
			case stepDone <- true:
			}

		}
	}()
	return stepDone
}

func generateBodies(nBodies int) []*physics.Body {
	bodies := make([]*physics.Body, nBodies)
	for i := 0; i < nBodies; i++ {
		mass := 1e20 + rand.Float64()*(1e30-1e20)
		bodies[i] = &physics.Body{
			ID:           i,
			Mass:         mass,
			Velocity:     physics.Vector{},
			Acceleration: physics.Vector{},
			Location:     physics.Vector{rand.Float64() * 1e10, rand.Float64() * 1e10, rand.Float64() * 1e10}}
	}
	return bodies
}

func main() {

	nBodies := flag.Int("bodies", 8, "number of bodies to simulate")
	steps := flag.Int("steps", 1000, "number of steps to simulate")
	nDaysPerStep := flag.Int("daysPerStep", 5, "number of days per simulation step")
	nThreads := flag.Int("threads", runtime.NumCPU(), "number of threads to parallelize")

	flag.Parse()

	if *nBodies == 1 {
		panic("Simulator requires at least two bodies.")
	}

	bodies := generateBodies(*nBodies)

	fmt.Println("Initial Locations:")
	for _, body := range bodies {
		fmt.Println("Body", body.ID, ":", body.Location)
	}

	done := make(chan interface{})
	defer close(done)

	stepDone := simulate(bodies, *nThreads, *nDaysPerStep, done)

	for range stepDone {
		if *steps == 0 {
			break
		}

		*steps--
	}
	fmt.Println("Final Locations:")
	for _, body := range bodies {
		fmt.Println("Body", body.ID, ":", body.Location)
	}

}
