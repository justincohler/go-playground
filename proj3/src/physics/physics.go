// Attribution for the body math & structure inspiration goes to:
//
// http://www.cyber-omelette.com/2016/11/python-n-body-orbital-simulation.html#theprogram
//
package physics

import (
	"math"
)

// Vector contains metadata associated with X, Y
type Vector struct {
	X float64
	Y float64
	Z float64
}

// Body contains position, Mass, and Velocity associated with a body
type Body struct {
	ID           int
	Mass         float64
	Velocity     Vector
	Acceleration Vector
	Location     Vector
}

// updateAcceleration performs a calculation to determine new acceleration for a body
func (body *Body) updateAcceleration(otherBodies []*Body) {
	const G = 6.67408e-11
	for _, otherBody := range otherBodies {
		if body != otherBody {
			dist := math.Pow((body.Location.X - otherBody.Location.X), 2)
			dist += math.Pow((body.Location.Y - otherBody.Location.Y), 2)
			dist += math.Pow((body.Location.Z - otherBody.Location.Z), 2)
			dist = math.Sqrt(dist)

			pull := G * otherBody.Mass / math.Pow(dist, 3)
			body.Acceleration.X += pull * (otherBody.Location.X - body.Location.X)
			body.Acceleration.Y += pull * (otherBody.Location.Y - body.Location.Y)
			body.Acceleration.Z += pull * (otherBody.Location.Z - body.Location.Z)
		}
	}

}

// updateVelocity performs a calcualtion to determine new velocity for a body
func (body *Body) updateVelocity(otherBodies []*Body, nSecondsPerStep float64) {
	body.updateAcceleration(otherBodies)
	body.Velocity.X += body.Acceleration.X * nSecondsPerStep
	body.Velocity.Y += body.Acceleration.Y * nSecondsPerStep
	body.Velocity.Z += body.Acceleration.Z * nSecondsPerStep
}

// UpdateLocation finds the neXt Location of a given body
func (body *Body) UpdateLocation(otherBodies []*Body, nDaysPerStep int) {
	nSecPerDay := 86400.0
	nSecondsPerStep := nSecPerDay * float64(nDaysPerStep)
	body.updateVelocity(otherBodies, nSecondsPerStep)
	body.Location.X += body.Velocity.X * nSecondsPerStep
	body.Location.Y += body.Velocity.Y * nSecondsPerStep
	body.Location.Z += body.Velocity.Z * nSecondsPerStep
}
