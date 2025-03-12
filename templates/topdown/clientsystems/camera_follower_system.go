package clientsystems

import (
	"math"

	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/warehouse"
)

type CameraFollowerSystem struct{}

func (CameraFollowerSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	// Query players who have a camera (camera index component)
	playersWithCamera := warehouse.Factory.NewQuery()
	playersWithCamera.And(
		blueprintspatial.Components.Position,
		blueprintinput.Components.InputBuffer,
		blueprintclient.Components.CameraIndex,
	)

	// Iterate
	playerCursor := scene.NewCursor(playersWithCamera)
	for playerCursor.Next() {

		// Get the players position
		playerPos := blueprintspatial.Components.Position.GetFromCursor(playerCursor)

		// Get the players camera
		camIndex := int(*blueprintclient.Components.CameraIndex.GetFromCursor(playerCursor))
		cam := cli.Cameras()[camIndex]

		// Get the cameras local scene position
		_, cameraScenePosition := cam.Positions()
		centerX := float64(cam.Surface().Bounds().Dx()) / 2
		centerY := float64(cam.Surface().Bounds().Dy()) / 2

		// Get target position (player) and adjust for centering
		targetX := playerPos.X - centerX
		targetY := playerPos.Y - centerY

		// Set the camera position towards the matched player position with lerp (optional)
		smoothness := 0.04
		cameraScenePosition.X, cameraScenePosition.Y = lerpVec(
			cameraScenePosition.X, cameraScenePosition.Y, targetX, targetY, smoothness,
		)

		// Rounding to avoid dealing with subpixels (inbound shovel's youtube has some good stuff on this)
		cameraScenePosition.X = math.Round(cameraScenePosition.X)
		cameraScenePosition.Y = math.Round(cameraScenePosition.Y)

		// Lock the camera to the scene boundaries
		lockCameraToSceneBoundaries(cam, scene, cameraScenePosition)
	}
	return nil
}

// lockCameraToSceneBoundaries constrains camera position within scene boundaries
func lockCameraToSceneBoundaries(cam coldbrew.Camera, scene coldbrew.Scene, cameraPos *vector.Two) {
	sceneWidth := scene.Width()
	sceneHeight := scene.Height()
	camWidth, camHeight := cam.Dimensions()

	// Calculate maximum positions to keep camera within scene bounds
	maxX := sceneWidth - camWidth
	maxY := sceneHeight - camHeight

	// Constrain camera X position
	if cameraPos.X > float64(maxX) {
		cameraPos.X = float64(maxX)
	}
	if cameraPos.X < 0 {
		cameraPos.X = 0
	}

	// Constrain camera Y position
	if cameraPos.Y > float64(maxY) {
		cameraPos.Y = float64(maxY)
	}
	if cameraPos.Y < 0 {
		cameraPos.Y = 0
	}
}

// Lerp func for smoothish (not perfect) camera
func lerpVec(startX, startY, endX, endY, t float64) (float64, float64) {
	dx := endX - startX
	dy := endY - startY
	return startX + dx*t, startY + dy*t
}
