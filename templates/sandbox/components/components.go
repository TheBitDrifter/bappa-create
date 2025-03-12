package components

import (
	"github.com/TheBitDrifter/warehouse"
)

var ExampleComponent = warehouse.FactoryNewComponent[Example]()

type Example struct {
	SomeState int
}
