package clientsystems

import (
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/warehouse"
)

type SortVerticalSystem struct{}

// This system assigns priority to sprites via their vertical position (y)
// This allows the GlobalRenderSystem to draw them in proper order
func (SortVerticalSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	// Query all non background spriteBundles
	query := warehouse.Factory.NewQuery().And(
		blueprintclient.Components.SpriteBundle,
		blueprintspatial.Components.Position,
		warehouse.Factory.NewQuery().Not(blueprintclient.Components.ParallaxBackground),
	)

	// Make a cursor from the query
	cursor := scene.NewCursor(query)

	// Iterate
	for cursor.Next() {
		// Get the position and sprite bundle
		pos := blueprintspatial.Components.Position.GetFromCursor(cursor)
		bundle := blueprintclient.Components.SpriteBundle.GetFromCursor(cursor)

		// Assign priority based on posY
		priority := pos.Y
		for i := range bundle.Blueprints {
			bp := &bundle.Blueprints[i]
			bp.Config.Priority = int(priority) + scene.Height()
		}
	}
	return nil
}
