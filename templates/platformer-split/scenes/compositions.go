package scenes

import (
	"github.com/TheBitDrifter/bappacreate/templates/platformer-split/components"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/warehouse"
)

// These are slices of common component compositions for various archetypes.
// They only include/represent the initial and static components of archetype
// Components can still be added or removed dynamically at runtime
//
// These slices are especially useful for creating starting entities, via archetypes, inside plan functions

var PlayerComposition = []warehouse.Component{
	blueprintspatial.Components.Position,
	blueprintclient.Components.SpriteBundle,
	blueprintspatial.Components.Direction,
	blueprintinput.Components.InputBuffer,
	blueprintclient.Components.CameraIndex,
	blueprintspatial.Components.Shape,
	blueprintmotion.Components.Dynamics,
	blueprintclient.Components.SoundBundle,
}

var BlockTerrainComposition = []warehouse.Component{
	components.BlockTerrainTag,
	blueprintspatial.Components.Shape,
	blueprintspatial.Components.Position,
	blueprintmotion.Components.Dynamics,
}

var PlatformComposition = []warehouse.Component{
	components.PlatformTag,
	blueprintspatial.Components.Rotation,
	blueprintclient.Components.SpriteBundle,
	blueprintspatial.Components.Shape,
	blueprintspatial.Components.Position,
	blueprintmotion.Components.Dynamics,
}

var MusicComposition = []warehouse.Component{
	blueprintclient.Components.SoundBundle,
	components.MusicTag,
}

var CollisionPlayerTransferComposition = []warehouse.Component{
	blueprintspatial.Components.Position,
	blueprintspatial.Components.Shape,
	components.PlayerSceneTransferComponent,
}
