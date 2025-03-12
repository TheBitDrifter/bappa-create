package components

import "github.com/TheBitDrifter/warehouse"

// Tags help us identify/categorize archetypes/entities when their
// composition alone isn't enough.
//
// For example its hard to tell the
// difference between a block and platform since they both have
// dynamics, shapes, sprites, etc
var (
	BlockTerrainTag = warehouse.FactoryNewComponent[struct{}]()
	PlatformTag     = warehouse.FactoryNewComponent[struct{}]()
	MusicTag        = warehouse.FactoryNewComponent[struct{}]()
)
