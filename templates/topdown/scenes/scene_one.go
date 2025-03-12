package scenes

import (
	"github.com/TheBitDrifter/blueprint"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/warehouse"
)

const SCENE_ONE_NAME = "scene one"

var SceneOne = Scene{
	Name:   SCENE_ONE_NAME,
	Plan:   sceneOnePlan,
	Width:  640,
	Height: 400,
}

func sceneOnePlan(height, width int, sto warehouse.Storage) error {
	// Create the player
	err := NewPlayer(sto)
	if err != nil {
		return err
	}

	// Background
	err = blueprint.CreateStillBackground(sto, "backgrounds/scene_one.png", vector.Two{X: 140, Y: 0})
	if err != nil {
		return err
	}
	// Music
	err = NewFantasyMusic(sto)
	if err != nil {
		return err
	}

	// Trees
	err = NewTreeProp(sto, 200, 100)
	if err != nil {
		return err
	}

	err = NewTreeProp(sto, 250, 368)
	if err != nil {
		return err
	}

	err = NewTreeProp(sto, 400, 300)
	if err != nil {
		return err
	}

	err = NewTreeProp(sto, 370, 150)
	if err != nil {
		return err
	}

	err = NewTreeProp(sto, 450, 100)
	if err != nil {
		return err
	}

	// Moveable statue
	err = NewMoveableStatueProp(sto, 340, 180)
	if err != nil {
		return err
	}

	// Left bounds
	err = NewBlockTerrain(sto, 140, 200, 10, 400)
	if err != nil {
		return err
	}
	// Right bounds
	err = NewBlockTerrain(sto, 510, 200, 10, 400)
	if err != nil {
		return err
	}
	// Top Bounds
	err = NewBlockTerrain(sto, 325, 0, 350, 10)
	if err != nil {
		return err
	}
	// Bottom horizontal left bounds
	err = NewBlockTerrain(sto, 220, 390, 147, 20)
	if err != nil {
		return err
	}
	// Bottom horizontal right bounds
	err = NewBlockTerrain(sto, 430, 390, 147, 20)
	if err != nil {
		return err
	}

	// Scene/Player transfer on collision
	colliderPosX := 325.0
	colliderPosY := 405.0

	colliderWidth := 60.0
	colliderHeight := 10.0

	targetLocationX := 225.0
	targetLocationY := 15.0
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
