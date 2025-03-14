package coresystems

import (
	"github.com/TheBitDrifter/bappacreate/templates/platformer-split/components"
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
			err := s.resolve(scene, blockTerrainCursor, playerCursor) // Now pass in the scene
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

	// Check grounded state
	playerAlreadyGrounded, onGround := components.OnGroundComponent.GetFromCursorSafe(playerCursor)

	// Check for a collision
	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *blockShape, playerPosition.Two, blockPosition.Two,
	); ok {

		playerOnTopOfBlock := collisionResult.IsTopB()

		// Prevents snapping onto terrain when the player is still jumping
		if playerAlreadyGrounded && onGround.LastTouch != scene.CurrentTick()-1 {
			playerAlreadyGrounded = false
		}
		if playerOnTopOfBlock && playerDynamics.Vel.Y < 0 && !playerAlreadyGrounded {
			return nil
		}

		// Determine if ground is sloped
		n := collisionResult.Normal
		horizontal := n.X == 0 && n.Y == 1 || n.X == 0 && n.Y == -1
		vertical := n.X == -1 && n.Y == 0 || n.X == 1 && n.Y == 0
		isSloped := !horizontal && !vertical

		// Vertical resolver to prevent positional sliding on slopes
		if isSloped {
			motion.VerticalResolver.Resolve(
				&playerPosition.Two,
				&blockPosition.Two,
				playerDynamics,
				blockDynamics,
				collisionResult,
			)
		} else {
			// Otherwise resolve as normal
			motion.Resolver.Resolve(
				&playerPosition.Two,
				&blockPosition.Two,
				playerDynamics,
				blockDynamics,
				collisionResult,
			)
		}

		// Ensure the player is on top of the terrain before marking them as grounded
		if !playerOnTopOfBlock {
			return nil
		}

		currentTick := scene.CurrentTick()

		// Update onGround accordingly (create or update)
		if !playerAlreadyGrounded {
			playerEntity, err := playerCursor.CurrentEntity()
			if err != nil {
				return err
			}
			// We cannot mutate during a cursor iteration, so we use the enqueue API
			err = playerEntity.EnqueueAddComponentWithValue(
				components.OnGroundComponent,
				components.OnGround{LastTouch: currentTick, Landed: currentTick, SlopeNormal: collisionResult.Normal},
			)
			if err != nil {
				return err
			}
		} else {
			onGround.LastTouch = scene.CurrentTick()
			onGround.SlopeNormal = collisionResult.Normal
		}

	}
	return nil
}
