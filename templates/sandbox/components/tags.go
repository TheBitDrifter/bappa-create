package components

import "github.com/TheBitDrifter/warehouse"

// Tags help us identify/categorize archetypes/entities when their composition alone isn't enough
var (
	ExampleTag = warehouse.FactoryNewComponent[struct{}]()
)
