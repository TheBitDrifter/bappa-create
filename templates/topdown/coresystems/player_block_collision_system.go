package coresystems

import (
	"github.com/TheBitDrifter/bappacreate/templates/topdown/components"
	"github.com/TheBitDrifter/blueprint"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/tteokbokki/motion"
	"github.com/TheBitDrifter/tteokbokki/spatial"
	"github.com/TheBitDrifter/warehouse"
)

type PlayerBlockCollisionSystem struct{}

func (s PlayerBlockCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	// Create cursors
	blockTerrainQuery := warehouse.Factory.NewQuery().And(components.BlockTerrainTag)
	blockTerrainCursor := scene.NewCursor(blockTerrainQuery)
	playerCursor := scene.NewCursor(blueprint.Queries.InputBuffer)

	// Outer loop is blocks
	for range blockTerrainCursor.Next() {
		// Inner is players
		for range playerCursor.Next() {
			// Delegate to helper
			err := s.resolve(scene, blockTerrainCursor, playerCursor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Main collision logic
func (PlayerBlockCollisionSystem) resolve(scene blueprint.Scene, blockCursor, playerCursor *warehouse.Cursor) error {
	// Get the player pos, shape, and dynamics
	playerPosition := blueprintspatial.Components.Position.GetFromCursor(playerCursor)
	playerShape := blueprintspatial.Components.Shape.GetFromCursor(playerCursor)
	playerDynamics := blueprintmotion.Components.Dynamics.GetFromCursor(playerCursor)

	// Get the block pos, shape, and dynamics
	blockPosition := blueprintspatial.Components.Position.GetFromCursor(blockCursor)
	blockShape := blueprintspatial.Components.Shape.GetFromCursor(blockCursor)
	blockDynamics := blueprintmotion.Components.Dynamics.GetFromCursor(blockCursor)

	// Check for a collision
	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *blockShape, playerPosition.Two, blockPosition.Two,
	); ok {
		// Resolve a collision
		motion.Resolver.Resolve(
			&playerPosition.Two,
			&blockPosition.Two,
			playerDynamics,
			blockDynamics,
			collisionResult,
		)
	}
	return nil
}
