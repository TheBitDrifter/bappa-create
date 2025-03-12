package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappacreate/templates/platformer/components"
	"github.com/TheBitDrifter/bappacreate/templates/platformer/sounds"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/warehouse"
)

type PlayerSoundSystem struct{}

func (sys PlayerSoundSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	// Create a query for players that can play sounds and are on the ground
	playersWithSoundsOnTheGround := warehouse.Factory.NewQuery()
	playersWithSoundsOnTheGround.And(
		blueprintclient.Components.SoundBundle, // Has sounds
		blueprintinput.Components.InputBuffer,  // Can receive input
		blueprintmotion.Components.Dynamics,    // Has physics properties
		components.OnGroundComponent,           // Is on the ground
	)

	// Get all entities that match the query
	cursor := scene.NewCursor(playersWithSoundsOnTheGround)

	// Iterate over all matching players
	for cursor.Next() {

		// Get state
		soundBundle := blueprintclient.Components.SoundBundle.GetFromCursor(cursor)
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)
		onGround := components.OnGroundComponent.GetFromCursor(cursor)
		currentTick := scene.CurrentTick()

		// Landed sound
		if onGround.Landed == currentTick {
			landingSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Land)
			player := landingSound.GetAny()

			// A hack to prevent landing sound artifacts between scenes
			// In a more robust setup, we might track if a player has recently changed scenes via a component
			// Such a component would be helpful here
			sceneRecentlySelected := scene.CurrentTick()-scene.LastSelectedTick() < 30

			if !player.IsPlaying() && !sceneRecentlySelected {
				player.Rewind()
				player.Play()
			}
		}

		// Jump sound

		// Gotta have y velocity
		if dyn.Vel.Y < 5 && onGround.LastJump == currentTick {
			jumpSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Jump)
			player := jumpSound.GetAny()

			if !player.IsPlaying() {
				player.Rewind()
				player.Play()
			}
		}

		// Run Sound

		// Must be moving horizontally
		const minMovementSpeed = 5.0
		if math.Abs(dyn.Vel.X) <= minMovementSpeed {
			continue
		}
		// Ensure onGround is not just available for coyote timer
		touchedGroundThisTick := onGround.LastTouch == currentTick
		if !touchedGroundThisTick {
			continue
		}

		runSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Run)
		if err != nil {
			return err
		}
		player := runSound.GetAny()

		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
		}
	}

	return nil
}
