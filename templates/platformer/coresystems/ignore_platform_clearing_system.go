package coresystems

import (
	"github.com/TheBitDrifter/bappacreate/templates/platformer/components"
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/warehouse"
)

// IgnorePlatformClearingSystem clears out individual platform ignores that have expired
type IgnorePlatformClearingSystem struct{}

func (IgnorePlatformClearingSystem) Run(scene blueprint.Scene, dt float64) error {
	// Query any entity that has IgnorePlatform
	ignorePlatformQuery := warehouse.Factory.NewQuery().And(components.IgnorePlatformComponent)
	ignorePlatformCursor := scene.NewCursor(ignorePlatformQuery)

	// Define the expiration time in ticks
	const expirationTicks = 15

	// Iterate through matched entities
	for ignorePlatformCursor.Next() {

		// Get the IgnorePlatform component state
		ignorePlatform := components.IgnorePlatformComponent.GetFromCursor(ignorePlatformCursor)
		currentTick := scene.CurrentTick()

		// Track if we have any active ignores left
		anyActive := false

		// Check each ignore entry
		for i := range ignorePlatform.Items {
			// Skip already cleared entries
			if ignorePlatform.Items[i].EntityID == 0 {
				continue
			}

			// Check if this entry has expired
			if currentTick-ignorePlatform.Items[i].LastActive > expirationTicks {

				// Clear this specific entry by setting its EntityID to 0
				ignorePlatform.Items[i].EntityID = 0
				ignorePlatform.Items[i].Recycled = 0
				ignorePlatform.Items[i].LastActive = 0

			} else {
				anyActive = true
			}
		}
		// If we don't have any active ignores left, remove the entire component
		if !anyActive {
			ignoringEntity, _ := ignorePlatformCursor.CurrentEntity()
			err := ignoringEntity.EnqueueRemoveComponent(components.IgnorePlatformComponent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
