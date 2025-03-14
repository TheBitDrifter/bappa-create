package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/actions"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/components"
	"github.com/TheBitDrifter/blueprint"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
)

type PlayerMovementSystem struct{}

func (sys PlayerMovementSystem) Run(scene blueprint.Scene, dt float64) error {
	// Get all entities with input buffers (players)
	cursor := scene.NewCursor(blueprint.Queries.InputBuffer)
	for range cursor.Next() {
		incomingInputs := blueprintinput.Components.InputBuffer.GetFromCursor(cursor)
		pos := blueprintspatial.Components.Position.GetFromCursor(cursor)
		direction := blueprintspatial.Components.Direction.GetFromCursor(cursor)
		direction8 := components.DirectionEightComponent.GetFromCursor(cursor)

		// Process horizontal movement
		_, pressedLeft := incomingInputs.ConsumeInput(actions.Left)
		_, pressedRight := incomingInputs.ConsumeInput(actions.Right)

		// Process vertical movement
		_, pressedUp := incomingInputs.ConsumeInput(actions.Up)
		_, pressedDown := incomingInputs.ConsumeInput(actions.Down)

		// Calculate movement vector
		moveX := 0.0
		moveY := 0.0

		if pressedLeft {
			moveX -= 1
		}
		if pressedRight {
			moveX += 1
		}
		if pressedUp {
			moveY -= 1
		}
		if pressedDown {
			moveY += 1
		}

		// Normalize diagonal movement
		if moveX != 0 && moveY != 0 {
			// Calculate vector length
			length := math.Sqrt(moveX*moveX + moveY*moveY)

			// Normalize the vector
			moveX = moveX / length
			moveY = moveY / length
		}

		// Apply movement
		pos.X += moveX
		pos.Y += moveY

		// Set direction components based on movement
		movingHorizontal := moveX != 0
		movingVertical := moveY != 0

		// Update cardinal direction component
		if moveX < 0 {
			direction.SetLeft()
		} else if moveX > 0 {
			direction.SetRight()
		}
		// Set the eight-direction component based on combined input
		if movingVertical || movingHorizontal {
			if pressedUp && pressedRight {
				direction8.SetRightUp()
			} else if pressedUp && pressedLeft {
				direction8.SetLeftUp()
			} else if pressedDown && pressedRight {
				direction8.SetRightDown()
			} else if pressedDown && pressedLeft {
				direction8.SetLeftDown()
			} else if pressedUp {
				direction8.SetUp()
			} else if pressedDown {
				direction8.SetDown()
			} else if pressedLeft {
				direction8.SetLeft()
			} else if pressedRight {
				direction8.SetRight()
			}
		}

		// Get player entity
		playerEn, err := cursor.CurrentEntity()
		if err != nil {
			return err
		}

		// Track movement
		if movingVertical || movingHorizontal {
			// Avoid/Cannot mutate during cursor iteration (locked storage)
			// Use Enqueue API
			playerEn.EnqueueAddComponent(components.IsMovingComponent)
		} else {
			playerEn.EnqueueRemoveComponent(components.IsMovingComponent)
		}
	}
	return nil
}
