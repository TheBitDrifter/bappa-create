package scenes

import "github.com/TheBitDrifter/blueprint"

// Local scene object makes it easier to organize
type Scene struct {
	Name          string
	Plan          blueprint.Plan
	Width, Height int
}
