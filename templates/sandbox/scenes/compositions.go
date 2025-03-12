package scenes

import (
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/warehouse"
)

var ExamplePlayerComposition = []warehouse.Component{
	blueprintspatial.Components.Position,
	blueprintclient.Components.SpriteBundle,
	blueprintspatial.Components.Direction,
	blueprintinput.Components.InputBuffer,
	blueprintclient.Components.CameraIndex,
	blueprintspatial.Components.Shape,
	blueprintmotion.Components.Dynamics,
	blueprintclient.Components.SoundBundle,
}
