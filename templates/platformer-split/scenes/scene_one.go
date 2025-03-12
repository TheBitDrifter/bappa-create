package scenes

import (
	"github.com/TheBitDrifter/warehouse"
)

const SCENE_ONE_NAME = "scene one"

var SceneOne = Scene{
	Name:   SCENE_ONE_NAME,
	Plan:   sceneOnePlan,
	Width:  1600,
	Height: 500,
}

// Scene one is a city scape
func sceneOnePlan(height, width int, sto warehouse.Storage) error {
	err := NewPlayers(sto, 2)
	if err != nil {
		return err
	}

	// Left/Right bounds
	err = NewInvisibleWalls(sto, width, height)
	if err != nil {
		return err
	}
	// Some platforms
	err = sceneOnePlatforms(sto)
	if err != nil {
		return err
	}

	// Block obstacle
	err = NewBlock(sto, 285, 390)
	if err != nil {
		return err
	}

	// Rampe
	err = NewRamp(sto, 470, 412)
	if err != nil {
		return err
	}

	// Floor
	err = NewFloor(sto, 460)
	if err != nil {
		return err
	}

	// Background
	err = NewCityBackground(sto)
	if err != nil {
		return err
	}

	// Music
	err = NewJazzMusic(sto)
	if err != nil {
		return err
	}

	// Scene/Player transfer on collision
	colliderPosX := float64(width)
	colliderPosY := 150.0

	colliderWidth := 11.0
	colliderHeight := float64(height)

	targetLocationX := 20.0
	targetLocationY := 400.0
	targetSceneName := SCENE_TWO_NAME

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

// Some platforms
func sceneOnePlatforms(sto warehouse.Storage) error {
	err := NewPlatform(sto, 130, 350)
	if err != nil {
		return err
	}
	err = NewPlatform(sto, 220, 270)
	if err != nil {
		return err
	}
	err = NewPlatform(sto, 320, 170)
	if err != nil {
		return err
	}
	err = NewPlatform(sto, 420, 300)
	if err != nil {
		return err
	}

	err = NewPlatformRotated(sto, 570, 278, 0.4)
	if err != nil {
		return err
	}
	err = NewPlatformRotated(sto, 700, 170, 0.7)
	if err != nil {
		return err
	}
	err = NewPlatformRotated(sto, 750, 278, -0.2)
	if err != nil {
		return err
	}
	return nil
}
