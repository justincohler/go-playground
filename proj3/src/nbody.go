package main

import (
	"fmt"
	"math"
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

func simulate(bodies []*Body, steps int) [][]Vector {
	N := len(bodies)

	// Make a history of all body positions for each step
	history := make([][]Vector, steps)
	for i := range history {
		history[i] = make([]Vector, N)
	}

	for i := 0; i < steps; i++ {
		for j, body := range bodies {
			body.UpdateLocation(bodies)
			history[i][j] = body.location
		}
	}
	return history
}

func main() {

	sun := &Body{
		mass:         2e30,
		velocity:     Vector{},
		acceleration: Vector{},
		location:     Vector{}}
	mercury := &Body{
		mass:         3e23,
		velocity:     Vector{47000, 0, 0},
		acceleration: Vector{},
		location:     Vector{0, 6e10, 0}}
	venus := &Body{
		mass:         4e24,
		velocity:     Vector{35000, 0, 0},
		acceleration: Vector{},
		location:     Vector{0, 1e11, 0}}

	bodies := []*Body{sun, mercury, venus}

	history := simulate(bodies, 10)

	for i, point := range history {
		fmt.Println(i, ":", point)
	}
}
