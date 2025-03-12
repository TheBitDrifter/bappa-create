package scenes

import (
	"github.com/TheBitDrifter/warehouse"
)

const EXAMPLE_SCENE_NAME = "scene one"

var SceneOne = Scene{
	Name:   EXAMPLE_SCENE_NAME,
	Plan:   examplePlan,
	Width:  640,
	Height: 360,
}

func examplePlan(height, width int, sto warehouse.Storage) error {
	return NewPlayer(sto)
}
