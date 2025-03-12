package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappacreate/templates/platformer-split/components"
	"github.com/TheBitDrifter/blueprint"

	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/tteokbokki/motion"
	"github.com/TheBitDrifter/tteokbokki/spatial"
	"github.com/TheBitDrifter/warehouse"
)

// PlayerPlatformCollisionSystem handles collisions between players and one-way platforms.
// It tracks historical player positions to determine if the player approached from above.
// This is necessary since collision detection at a discrete step doesn't provide approach direction.
//
// It is one of the more 'dense' bits of the template/example.
type PlayerPlatformCollisionSystem struct {
	// Map of player entity ID to their position history
	playerPositionHistory map[uint64][]vector.Two
	maxPositionsToTrack   int // Number of positions to track
}

// NewPlayerPlatformCollisionSystem creates a new collision system with initialized position tracking.
// It uses a pointer because the system is not pure and must retain its state.
func NewPlayerPlatformCollisionSystem() *PlayerPlatformCollisionSystem {
	trackCount := 15 // higher count == more tunneling protection == higher cost
	return &PlayerPlatformCollisionSystem{
		playerPositionHistory: make(map[uint64][]vector.Two),
		maxPositionsToTrack:   trackCount,
	}
}

func (s *PlayerPlatformCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	// Create cursors
	platformTerrainQuery := warehouse.Factory.NewQuery().And(components.PlatformTag)
	platformCursor := scene.NewCursor(platformTerrainQuery)
	playerCursor := scene.NewCursor(blueprint.Queries.InputBuffer)

	// Outer loop is platforms
	for platformCursor.Next() {
		// Inner is players
		for playerCursor.Next() {
			// Get player entity ID to track positions independently for each player
			playerEntity, err := playerCursor.CurrentEntity()
			if err != nil {
				return err
			}
			playerID := uint64(playerEntity.ID())

			// Delegate to helper
			err = s.resolve(scene, platformCursor, playerCursor, playerID)
			if err != nil {
				return err
			}

			// Track the full position (X and Y) for this specific player
			playerPos := blueprintspatial.Components.Position.GetFromCursor(playerCursor)
			s.trackPosition(playerID, playerPos.Two)
		}
	}
	return nil
}

func (s *PlayerPlatformCollisionSystem) resolve(scene blueprint.Scene, platformCursor, playerCursor *warehouse.Cursor, playerID uint64) error {
	// Get the player state
	playerShape := blueprintspatial.Components.Shape.GetFromCursor(playerCursor)
	playerPosition := blueprintspatial.Components.Position.GetFromCursor(playerCursor)
	playerDynamics := blueprintmotion.Components.Dynamics.GetFromCursor(playerCursor)

	// Get the platform state
	platformShape := blueprintspatial.Components.Shape.GetFromCursor(platformCursor)
	platformPosition := blueprintspatial.Components.Position.GetFromCursor(platformCursor)
	platformRotation := float64(*blueprintspatial.Components.Rotation.GetFromCursor(platformCursor))
	platformDynamics := blueprintmotion.Components.Dynamics.GetFromCursor(platformCursor)

	// Check for collision
	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *platformShape, playerPosition.Two, platformPosition.Two,
	); ok {

		// Check if were ignoring the current platform (dropping down)
		ignoringPlatforms, ignorePlatform := components.IgnorePlatformComponent.GetFromCursorSafe(playerCursor)

		platformEntity, err := platformCursor.CurrentEntity()
		if err != nil {
			return err
		}
		if ignoringPlatforms {
			for _, ignored := range ignorePlatform.Items {
				if ignored.EntityID == int(platformEntity.ID()) && ignored.Recycled == platformEntity.Recycled() {
					return nil
				}
			}
		}

		// Check if any of the past player positions indicate the player was above the platform
		platformTop := platformShape.Polygon.WorldVertices[0].Y
		var playerWasAbove bool

		// Checking for 'above' is much easier when the edge is flat (fixed y value)
		if platformRotation == 0 {
			playerWasAbove = s.checkAnyPlayerPositionWasAbove(playerID, platformTop, playerShape.LocalAAB.Height)

			// Rotation check is more complicated using vector math to determine if player 'cleared top'
		} else {
			playerWasAbove = s.checkAnyPlayerPositionWasAboveAdvanced(
				playerID,
				// The top edge for the triangle platforms is always 0,1
				[]vector.Two{
					platformShape.Polygon.WorldVertices[0],
					platformShape.Polygon.WorldVertices[1],
				},
				// Pass the AAB dimensions to calc the players bottom points along with their historical positions
				playerShape.LocalAAB.Width, playerShape.LocalAAB.Height,
				platformRotation,
			)
		}

		// We only want to resolve collisions when:
		// 1. The player is falling (vel.Y > 0)
		// 2. The collision is with the top of the platform
		// 3. The player was above the platform at some point (within n ticks)
		if playerDynamics.Vel.Y > 0 && collisionResult.IsTopB() && playerWasAbove {

			// Use a vertical resolver since we can't collide with the sides
			motion.VerticalResolver.Resolve(
				&playerPosition.Two,
				&platformPosition.Two,
				playerDynamics,
				platformDynamics,
				collisionResult,
			)

			// Standard onGround handling
			currentTick := scene.CurrentTick()

			// If not grounded, enqueue onGround with values
			playerAlreadyGrounded, onGround := components.OnGroundComponent.GetFromCursorSafe(playerCursor)

			if !playerAlreadyGrounded {
				playerEntity, _ := playerCursor.CurrentEntity()
				err := playerEntity.EnqueueAddComponentWithValue(
					components.OnGroundComponent,
					components.OnGround{LastTouch: currentTick, Landed: currentTick, SlopeNormal: collisionResult.Normal},
				)
				if err != nil {
					return err
				}
			} else {

				// Otherwise update the existing OnGround
				onGround.LastTouch = scene.CurrentTick()
				onGround.SlopeNormal = collisionResult.Normal
			}

			// If player is ignoring platforms and we have reached here, they aren't ignoring this one yet
			// So replace the oldest ignoredPlatform with this one
			if ignoringPlatforms {
				var oldestTick int64 = math.MaxInt64
				var oldestIndex int = -1

				for i, ignored := range ignorePlatform.Items {
					if ignored.EntityID == int(platformEntity.ID()) && ignored.Recycled == platformEntity.Recycled() {
						return nil
					}

					// Track the oldest tick
					if int64(ignored.LastActive) < oldestTick {
						oldestTick = int64(ignored.LastActive)
						oldestIndex = i
					}
				}

				// If we found an oldest entry, replace it with the current platform entity
				if oldestIndex != -1 {
					ignorePlatform.Items[oldestIndex].EntityID = int(platformEntity.ID())
					ignorePlatform.Items[oldestIndex].Recycled = platformEntity.Recycled()
					ignorePlatform.Items[oldestIndex].LastActive = currentTick
					return nil
				}
			}

		}
	}
	return nil
}

// trackPosition adds a position to the history and ensures only the last N are kept for a specific player
func (s *PlayerPlatformCollisionSystem) trackPosition(playerID uint64, pos vector.Two) {
	// Initialize the position history for this player if it doesn't exist
	if _, exists := s.playerPositionHistory[playerID]; !exists {
		s.playerPositionHistory[playerID] = make([]vector.Two, 0, s.maxPositionsToTrack)
	}

	// Add the new position to this player's history
	s.playerPositionHistory[playerID] = append(s.playerPositionHistory[playerID], pos)

	// If we've exceeded our max, remove the oldest position
	if len(s.playerPositionHistory[playerID]) > s.maxPositionsToTrack {
		s.playerPositionHistory[playerID] = s.playerPositionHistory[playerID][1:]
	}
}

// checkAnyPlayerPositionWasAbove checks if the player was above a non-rotated platform in any historical position
func (s *PlayerPlatformCollisionSystem) checkAnyPlayerPositionWasAbove(playerID uint64, platformTop float64, playerHeight float64) bool {
	positions, exists := s.playerPositionHistory[playerID]
	if !exists || len(positions) == 0 {
		return false
	}

	// Check all stored positions to see if the player was above in any of them
	for _, pos := range positions {
		playerBottom := pos.Y + playerHeight/2
		if playerBottom <= platformTop {
			return true // Found at least one position where player was above
		}
	}

	return false // No positions found where player was above
}

// checkAnyPlayerPositionWasAboveAdvanced checks if the player was above a rotated platform's top edge
func (s *PlayerPlatformCollisionSystem) checkAnyPlayerPositionWasAboveAdvanced(
	playerID uint64,
	platformTopVerts []vector.Two,
	playerWidth, playerHeight, rotation float64,
) bool {
	positions, exists := s.playerPositionHistory[playerID]
	if !exists || len(positions) == 0 {
		return false
	}

	// Get the vertices of the platform top edge
	v1 := platformTopVerts[0]
	v2 := platformTopVerts[1]

	// Calculate the edge vector and its length
	edgeVector := v2.Sub(v1)
	edgeLength := edgeVector.Mag()
	if edgeLength < 0.001 {
		return false // Avoid division by zero
	}

	// Normalize the edge vector
	edgeNormalized := edgeVector.Norm()

	// Calculate the normal vector perpendicular to the edge
	// Make sure it points upward (negative Y in this coordinate system)
	edgeNormal := vector.Two{X: -edgeNormalized.Y, Y: edgeNormalized.X}

	// Ensure the normal points upward (negative Y)
	if edgeNormal.Y > 0 {
		edgeNormal = edgeNormal.Scale(-1)
	}

	// For each historical position
	for _, historicalPos := range positions {
		// Calculate player points to check (bottom center, bottom left, and bottom right)
		checkPoints := []vector.Two{
			// Bottom center
			{
				X: historicalPos.X,
				Y: historicalPos.Y + playerHeight/2,
			},
		}

		// Add appropriate side points based on rotation direction
		// For negative rotation, the platform slopes down to the right,
		// so we need to check the right side more carefully
		if rotation < 0 {
			checkPoints = append(checkPoints, vector.Two{
				X: historicalPos.X + playerWidth/2 - 5, // Right side with small inset
				Y: historicalPos.Y + playerHeight/2,
			})
		}

		// For positive rotation, the platform slopes down to the left,
		// so we need to check the left side more carefully
		if rotation > 0 {
			checkPoints = append(checkPoints, vector.Two{
				X: historicalPos.X - playerWidth/2 + 5, // Left side with small inset
				Y: historicalPos.Y + playerHeight/2,
			})
		}

		// Check each point of the player's bottom
		for _, point := range checkPoints {
			// Vector from edge start to point
			v1ToPoint := point.Sub(v1)

			// Distance along the normal (positive means above the edge)
			distanceAlongNormal := v1ToPoint.ScalarProduct(edgeNormal)

			// Project onto the edge to check if within bounds
			projectionOnEdge := v1ToPoint.ScalarProduct(edgeNormalized)

			// Constants for checks
			const margin = 10.0   // Margin for edge projection
			const minAbove = 1.0  // Minimum distance to be considered "above"
			const maxAbove = 50.0 // Maximum distance to be considered relevant

			// A positive distance along normal means the point is in the direction of the normal
			// The normal points upward, so positive distance means "above" the edge
			isAbove := distanceAlongNormal > minAbove &&
				distanceAlongNormal < maxAbove &&
				projectionOnEdge >= -margin &&
				projectionOnEdge <= edgeLength+margin

			if isAbove {
				return true
			}
		}
	}

	return false
}
