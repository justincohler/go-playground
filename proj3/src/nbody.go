// Attribution for the initial body math & structure goes to:
//
// http://www.cyber-omelette.com/2016/11/python-n-body-orbital-simulation.html#theprogram
//
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"physics"
	"runtime"
	"strconv"
)

func blockUpdateLocations(start int, end int, bodies []*physics.Body, bodyCounter chan<- int, nDaysPerStep int) {
	for j := start; j < end; j++ {
		bodies[j].UpdateLocation(bodies, nDaysPerStep)
		bodyCounter <- j
	}
}

func simulate(bodies []*physics.Body, nThreads int, nDaysPerStep int) <-chan interface{} {
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
			stepDone <- true
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

	nBodies, _ := strconv.Atoi(os.Args[1])
	steps, _ := strconv.Atoi(os.Args[2])
	nDaysPerStep, _ := strconv.Atoi(os.Args[3])

	if nBodies == 1 {
		panic("Simulator requires at least two bodies.")
	}

	bodies := generateBodies(nBodies)

	nThreads := runtime.NumCPU()
	stepDone := simulate(bodies, nThreads, nDaysPerStep)

	for range stepDone {
		if steps == 0 {
			break
		}
		for _, body := range bodies {
			if steps%1000 == 0 {
				fmt.Print(body.Location, ";")
			}
		}
		steps--
	}
}
