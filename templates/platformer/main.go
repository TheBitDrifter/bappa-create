package main

import (
	"embed"
	"log"

	"github.com/TheBitDrifter/bappacreate/templates/platformer/actions"
	"github.com/TheBitDrifter/bappacreate/templates/platformer/clientsystems"
	"github.com/TheBitDrifter/bappacreate/templates/platformer/coresystems"
	"github.com/TheBitDrifter/bappacreate/templates/platformer/rendersystems"
	"github.com/TheBitDrifter/bappacreate/templates/platformer/scenes"
	"github.com/TheBitDrifter/coldbrew"
	coldbrew_clientsystems "github.com/TheBitDrifter/coldbrew/clientsystems"
	coldbrew_rendersystems "github.com/TheBitDrifter/coldbrew/rendersystems"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

const (
	RESOLUTION_X       = 640
	RESOLUTION_Y       = 360
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
	client.SetTitle("Platformer Template")
	client.SetResizable(true)
	client.SetMinimumLoadTime(30)

	// Register scene One
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
	)

	// Activate camera
	client.ActivateCamera()

	// Register receiver/actions
	receiver1, _ := client.ActivateReceiver()
	receiver1.RegisterKey(ebiten.KeySpace, actions.Jump)
	receiver1.RegisterKey(ebiten.KeyW, actions.Jump)
	receiver1.RegisterKey(ebiten.KeyA, actions.Left)
	receiver1.RegisterKey(ebiten.KeyD, actions.Right)
	receiver1.RegisterKey(ebiten.KeyS, actions.Down)

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}
