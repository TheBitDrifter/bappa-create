package scenes

import "github.com/TheBitDrifter/warehouse"

const SCENE_TWO_NAME = "scene two"

var SceneTwo = Scene{
	Name:   SCENE_TWO_NAME,
	Plan:   sceneTwoPlan,
	Width:  1600,
	Height: 500,
}

// Scene two is a simple night sky and floor
func sceneTwoPlan(height, width int, sto warehouse.Storage) error {
	err := NewInvisibleWalls(sto, width, height)
	if err != nil {
		return err
	}

	// Floor
	err = NewFloor(sto, 460)
	if err != nil {
		return err
	}

	// Background
	err = NewSkyBackground(sto)
	if err != nil {
		return err
	}
	// Scene/Player transfer on collision
	colliderPosX := 0.0
	colliderPosY := 150.0

	colliderWidth := 11.0
	colliderHeight := float64(height)

	targetLocationX := float64(width - 20)
	targetLocationY := 400.0
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
