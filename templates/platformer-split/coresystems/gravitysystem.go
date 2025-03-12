package coresystems

import (
	"github.com/TheBitDrifter/blueprint"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	"github.com/TheBitDrifter/tteokbokki/motion"
)

const (
	DEFAULT_GRAVITY  = 9.8
	PIXELS_PER_METER = 50.0
)

type GravitySystem struct{}

func (GravitySystem) Run(scene blueprint.Scene, dt float64) error {
	// Iterate through entities with dynamics components(physics)
	cursor := scene.NewCursor(blueprint.Queries.Dynamics)
	for cursor.Next() {
		// Get the dynamics
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)

		// Get the mass
		mass := 1 / dyn.InverseMass

		// Use the motion package to calc the gravity force
		gravity := motion.Forces.Generator.NewGravityForce(mass, DEFAULT_GRAVITY, PIXELS_PER_METER)

		// Apply the force
		motion.Forces.AddForce(dyn, gravity)
	}
	return nil
}
