package coresystems

import (
	"github.com/TheBitDrifter/blueprint"
	tteo_coresystems "github.com/TheBitDrifter/tteokbokki/coresystems"
)

var DefaultCoreSystems = []blueprint.CoreSystem{
	PlayerMovementSystem{},               // Apply player input forces
	tteo_coresystems.IntegrationSystem{}, // Update velocities and positions
	tteo_coresystems.TransformSystem{},   // Update collision shapes
	PlayerBlockCollisionSystem{},         // Handle  collisions
}
