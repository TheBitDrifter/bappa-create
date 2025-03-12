package components

import (
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/warehouse"
)

var (
	// For tracking grounded state
	OnGroundComponent = warehouse.FactoryNewComponent[OnGround]()

	// For ignoring/dropping down one way platforms
	IgnorePlatformComponent = warehouse.FactoryNewComponent[IgnorePlatform]()

	// For scene transfers
	PlayerSceneTransferComponent = warehouse.FactoryNewComponent[PlayerSceneTransfer]()
)

// For tracking ground interactions
type OnGround struct {
	LastTouch   int
	Landed      int
	LastJump    int
	SlopeNormal vector.Two
}

// For dropping down platforms
type IgnorePlatform struct {
	Items [5]struct {
		LastActive int
		EntityID   int
		Recycled   int
	}
}

// For player/scene transfers/changes
type PlayerSceneTransfer struct {
	Dest string
	X, Y float64
}
