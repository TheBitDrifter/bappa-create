package components

import (
	"github.com/TheBitDrifter/warehouse"
)

// Here we convert our types into components for the ECS
var (
	PlayerSceneTransferComponent = warehouse.FactoryNewComponent[PlayerSceneTransfer]()
	DirectionEightComponent      = warehouse.FactoryNewComponent[DirectionEight]()
	IsMovingComponent            = warehouse.FactoryNewComponent[struct{}]()
)
