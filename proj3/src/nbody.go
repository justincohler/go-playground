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
	"runtime"
	"strconv"
)

// Vector contains metadata associated with x, y
type Vector struct {
	x float64
	y float64
	z float64
}

// Body contains position, mass, and velocity associated with a body
type Body struct {
	mass         float64
	velocity     Vector
	acceleration Vector
	location     Vector
}

func (body *Body) updateAcceleration(otherBodies []*Body) {
	const G = 6.67408e-11
	for _, otherBody := range otherBodies {
		if body != otherBody {
			dist := math.Pow((body.location.x - otherBody.location.x), 2)
			dist += math.Pow((body.location.y - otherBody.location.y), 2)
			dist += math.Pow((body.location.z - otherBody.location.z), 2)
			dist = math.Sqrt(dist)
			pull := G * otherBody.mass / math.Pow(dist, 3)
			body.acceleration.x += pull * (otherBody.location.x - body.location.x)
			body.acceleration.y += pull * (otherBody.location.y - body.location.y)
			body.acceleration.z += pull * (otherBody.location.z - body.location.z)
		}
	}

}

func (body *Body) updateVelocity(otherBodies []*Body) {
	body.updateAcceleration(otherBodies)
	body.velocity.x += body.acceleration.x
	body.velocity.y += body.acceleration.y
	body.velocity.z += body.acceleration.z
}

// UpdateLocation finds the next location of a given body
func (body *Body) UpdateLocation(otherBodies []*Body) {
	body.updateVelocity(otherBodies)
	body.location.x += body.velocity.x
	body.location.y += body.velocity.y
	body.location.z += body.velocity.z
}

func simulate(bodies []*Body, nThreads int) <-chan []Vector {
	N := len(bodies)
	stepChan := make(chan []Vector)
	bodyCounter := make(chan interface{}, N)

	go func() {
		for {
			stepLocations := make([]Vector, N)
			for i := 0; i < nThreads; i++ {
				blockSize := math.Ceil(float64(N) / float64(nThreads))
				start := int(float64(i) * blockSize)
				end := int(math.Min(float64(N), (float64(i)+1)*blockSize))
				go func() { // functional parallelism
					for j := start; j < end; j++ {
						body := bodies[j]
						body.UpdateLocation(bodies)
						stepLocations[j] = body.location
						bodyCounter <- j
					}
				}()
			}

			for i := 0; i < nThreads; i++ { // wait for step to complete
				<-bodyCounter
			}
			stepChan <- stepLocations
		}
	}()
	return stepChan
}

func generateBodies(nBodies int) []*Body {
	bodies := make([]*Body, nBodies)
	for i := 0; i < nBodies; i++ {
		mass := 1e20 + rand.Float64()*(1e30-1e20)
		bodies[i] = &Body{
			mass:         mass,
			velocity:     Vector{rand.Float64() * 100000, rand.Float64() * 100000, rand.Float64() * 100000},
			acceleration: Vector{},
			location:     Vector{rand.Float64() * 1e20, rand.Float64() * 1e20, rand.Float64() * 1e20}}
	}
	return bodies
}

func main() {

	nBodies, _ := strconv.Atoi(os.Args[1])
	steps, _ := strconv.Atoi(os.Args[2])
	bodies := generateBodies(nBodies)

	nThreads := runtime.NumCPU()
	stepChan := simulate(bodies, nThreads)

	for stepLocations := range stepChan {
		if steps == 0 {
			break
		}
		fmt.Println(stepLocations)
		steps--
	}
}
