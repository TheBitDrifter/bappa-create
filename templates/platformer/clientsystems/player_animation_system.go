package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappacreate/templates/platformer/animations"
	"github.com/TheBitDrifter/bappacreate/templates/platformer/components"
	"github.com/TheBitDrifter/blueprint"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	"github.com/TheBitDrifter/coldbrew"
)

type PlayerAnimationSystem struct{}

func (PlayerAnimationSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	cursor := scene.NewCursor(blueprint.Queries.InputBuffer)

	for cursor.Next() {
		// Get state
		bundle := blueprintclient.Components.SpriteBundle.GetFromCursor(cursor)
		spriteBlueprint := &bundle.Blueprints[0]
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)
		grounded, onGround := components.OnGroundComponent.GetFromCursorSafe(cursor)
		if grounded {
			grounded = scene.CurrentTick() == onGround.LastTouch
		}

		// Player is moving horizontal and grounded (running)
		if math.Abs(dyn.Vel.X) > 0 && grounded {
			spriteBlueprint.TryAnimation(animations.RunAnimation)

			// Player is moving down and not grounded (falling)
		} else if dyn.Vel.Y > 0 && !grounded {
			spriteBlueprint.TryAnimation(animations.FallAnimation)

			// Player is moving up and not grounded (jumping)
		} else if dyn.Vel.Y <= 0 && !grounded {
			spriteBlueprint.TryAnimation(animations.JumpAnimation)

			// Default: player is idle
		} else {
			spriteBlueprint.TryAnimation(animations.IdleAnimation)
		}
	}
	return nil
}
