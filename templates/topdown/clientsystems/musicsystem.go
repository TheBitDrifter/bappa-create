package clientsystems

import (
	"github.com/TheBitDrifter/bappacreate/templates/topdown/components"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/sounds"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/warehouse"
)

type MusicSystem struct{}

// Note: this a very simple music system that does not account for multiple scenes (will only loop on scene one)
// Adjust accordingly
func (sys MusicSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	// Setup query and cursor for music
	musicQuery := warehouse.Factory.NewQuery().And(components.MusicTag)
	cursor := scene.NewCursor(musicQuery)

	// There's only one but iterate nonetheless
	for range cursor.Next() {
		soundBundle := blueprintclient.Components.SoundBundle.GetFromCursor(cursor)

		sound, err := coldbrew.MaterializeSound(soundBundle, sounds.Music)
		if err != nil {
			return err
		}
		player := sound.GetAny()

		// Loop if needed
		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
		}
	}
	return nil
}
