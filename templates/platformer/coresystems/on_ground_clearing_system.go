package coresystems

import (
	"github.com/TheBitDrifter/bappacreate/templates/platformer/components"
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/warehouse"
)

type OnGroundClearingSystem struct{}

func (OnGroundClearingSystem) Run(scene blueprint.Scene, dt float64) error {
	// Define the expiration time in ticks
	// 15 is a bit of a magic number— what matters is that the value is greater than coyote timer
	// Systems should also check lastTouch alongside presence to avoid odd behavior
	const expirationTicks = 15

	// Query any entity that has onGround
	onGroundQuery := warehouse.Factory.NewQuery().And(components.OnGroundComponent)
	onGroundCursor := scene.NewCursor(onGroundQuery)

	// Iterate through matched entities
	for onGroundCursor.Next() {
		// Get the onGround component state
		onGround := components.OnGroundComponent.GetFromCursor(onGroundCursor)

		// If it's expired, remove it
		if scene.CurrentTick()-onGround.LastTouch > expirationTicks {
			groundedEntity, _ := onGroundCursor.CurrentEntity()

			// We can't mutate while iterating so we enqueue the changes instead
			err := groundedEntity.EnqueueRemoveComponent(components.OnGroundComponent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
