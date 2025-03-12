package scenes

import (
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/warehouse"
)

const SCENE_TWO_NAME = "scene two"

var SceneTwo = Scene{
	Name:   SCENE_TWO_NAME,
	Plan:   sceneTwoPlan,
	Width:  1600,
	Height: 500,
}

func sceneTwoPlan(height, width int, sto warehouse.Storage) error {
	// Left bounds
	err := NewBlockTerrain(sto, 140, 200, 10, 400)
	if err != nil {
		return err
	}

	// Right bounds
	err = NewBlockTerrain(sto, 510, 200, 10, 400)
	if err != nil {
		return err
	}

	// Bottom bounds
	err = NewBlockTerrain(sto, 325, 385, 350, 10)
	if err != nil {
		return err
	}

	// Top left horizontal bounds
	err = NewBlockTerrain(sto, 170, 10, 38, 10)
	if err != nil {
		return err
	}

	// Top right horizontal bounds
	err = NewBlockTerrain(sto, 383, 10, 230, 10)
	if err != nil {
		return err
	}

	// Top left vertical bounds
	err = NewBlockTerrain(sto, 194, 25, 6, 60)
	if err != nil {
		return err
	}

	// Top right vertical bounds
	err = NewBlockTerrain(sto, 262, 25, 6, 60)
	if err != nil {
		return err
	}

	// Background
	err = blueprint.CreateStillBackground(sto, "backgrounds/scene_two.png", vector.Two{X: 140, Y: 0})
	if err != nil {
		return err
	}

	// Scene/Player transfer on collision
	colliderPosX := 228.0
	colliderPosY := -7.0

	colliderWidth := 58.0
	colliderHeight := 20.0

	targetLocationX := 317.0
	targetLocationY := 385.0
	targetSceneName := SCENE_ONE_NAME

	// Since its the the last one and has the same func sig as the parent, we can return it
	return NewCollisionPlayerTransfer(
		sto,
		colliderPosX,
		colliderPosY,
		colliderWidth,
		colliderHeight,
		targetLocationX,
		targetLocationY,
		targetSceneName,
	)
}
