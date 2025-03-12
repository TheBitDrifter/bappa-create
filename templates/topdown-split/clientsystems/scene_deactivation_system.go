package clientsystems

import (
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/warehouse"
)

type SceneDeactivationSystem struct{}

// Global client system deactivates scenes with no players
// Watch out for this disabling scenes that are UI based only
// You may need to handle that (add a player or modify the check)!
func (SceneDeactivationSystem) Run(cli coldbrew.Client) error {
	for _, scene := range cli.ActiveScenes() {
		cursor := warehouse.Factory.NewCursor(blueprint.Queries.InputBuffer, scene.Storage())
		hasPlayers := cursor.TotalMatched() > 0
		if !hasPlayers {
			cli.DeactivateScene(scene)
		}
	}
	return nil
}
