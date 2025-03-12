package scenes

import (
	"log"

	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/warehouse"
)

// NewPlayer creates a 'player' entity for the scene
func NewPlayer(sto warehouse.Storage) error {
	playerArchetype, err := sto.NewOrExistingArchetype(
		ExamplePlayerComposition...,
	)
	log.Println("yoooo")
	err = playerArchetype.Generate(1,
		blueprintspatial.NewPosition(180, 180),
		blueprintspatial.NewRectangle(18, 58),
		blueprintmotion.NewDynamics(10),
		blueprintspatial.NewDirectionRight(),
		blueprintinput.InputBuffer{ReceiverIndex: 0},
		blueprintclient.CameraIndex(0),
	)
	if err != nil {
		return err
	}
	return nil
}
