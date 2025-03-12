package main

import (
	"embed"
	"log"

	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/actions"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/clientsystems"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/coresystems"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/rendersystems"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/scenes"
	"github.com/TheBitDrifter/coldbrew"
	coldbrew_clientsystems "github.com/TheBitDrifter/coldbrew/clientsystems"
	coldbrew_rendersystems "github.com/TheBitDrifter/coldbrew/rendersystems"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

const (
	RESOLUTION_X       = 640
	RESOLUTION_Y       = 720
	MAX_SPRITES_CACHED = 100
	MAX_SOUNDS_CACHED  = 100
	MAX_SCENES_CACHED  = 12
)

func main() {
	// Create the client
	client := coldbrew.NewClient(
		RESOLUTION_X,
		RESOLUTION_Y,
		MAX_SPRITES_CACHED,
		MAX_SOUNDS_CACHED,
		MAX_SCENES_CACHED,
		assets,
	)

	// Settings
	client.SetTitle("Topdown Template Split Screen")
	client.SetResizable(true)
	client.SetMinimumLoadTime(8)
	client.SetCameraBorderSize(5)

	// Register scene one
	err := client.RegisterScene(
		scenes.SceneOne.Name,
		scenes.SceneOne.Width,
		scenes.SceneOne.Height,
		scenes.SceneOne.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register scene two
	err = client.RegisterScene(
		scenes.SceneTwo.Name,
		scenes.SceneTwo.Width,
		scenes.SceneTwo.Height,
		scenes.SceneTwo.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystems,
		coresystems.DefaultCoreSystems,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Register global systems
	client.RegisterGlobalRenderSystem(
		coldbrew_rendersystems.GlobalRenderer{},
		&coldbrew_rendersystems.DebugRenderer{},
	)
	client.RegisterGlobalClientSystem(
		coldbrew_clientsystems.InputBufferSystem{},
		&coldbrew_clientsystems.CameraSceneAssignerSystem{},
		clientsystems.SceneDeactivationSystem{},
	)

	// Activate camera
	cameraOne, err := client.ActivateCamera()
	if err != nil {
		log.Fatal(err)
	}
	cameraOne.SetDimensions(RESOLUTION_X, RESOLUTION_Y/2)

	cameraTwo, err := client.ActivateCamera()
	if err != nil {
		log.Fatal(err)
	}
	cameraTwo.SetDimensions(RESOLUTION_X, RESOLUTION_Y/2)
	camTwoScreenPos, _ := cameraTwo.Positions()
	camTwoScreenPos.Y = RESOLUTION_Y / 2

	// Register receiver/actions
	receiver1, _ := client.ActivateReceiver()
	receiver1.RegisterKey(ebiten.KeyW, actions.Up)
	receiver1.RegisterKey(ebiten.KeyA, actions.Left)
	receiver1.RegisterKey(ebiten.KeyD, actions.Right)
	receiver1.RegisterKey(ebiten.KeyS, actions.Down)

	receiver2, _ := client.ActivateReceiver()
	receiver2.RegisterKey(ebiten.KeyUp, actions.Up)
	receiver2.RegisterKey(ebiten.KeyLeft, actions.Left)
	receiver2.RegisterKey(ebiten.KeyRight, actions.Right)
	receiver2.RegisterKey(ebiten.KeyDown, actions.Down)

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}
