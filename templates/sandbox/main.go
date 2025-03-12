package main

import (
	"embed"
	"log"

	"github.com/TheBitDrifter/bappacreate/templates/sandbox/clientsystems"
	"github.com/TheBitDrifter/bappacreate/templates/sandbox/coresystems"
	"github.com/TheBitDrifter/bappacreate/templates/sandbox/rendersystems"
	"github.com/TheBitDrifter/bappacreate/templates/sandbox/scenes"
	"github.com/TheBitDrifter/coldbrew"
	coldbrew_clientsystems "github.com/TheBitDrifter/coldbrew/clientsystems"
	coldbrew_rendersystems "github.com/TheBitDrifter/coldbrew/rendersystems"
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

	client.SetTitle("Sandbox Template")
	client.SetResizable(true)
	client.SetMinimumLoadTime(8)
	client.SetEnforceMinOnActive(true)

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

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}
