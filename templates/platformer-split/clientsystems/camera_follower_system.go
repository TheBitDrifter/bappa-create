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
	for range playerCursor.Next() {
		// Get the players position
		playerPos := blueprintspatial.Components.Position.GetFromCursor(playerCursor)
		// Get the players camera
		camIndex := int(*blueprintclient.Components.CameraIndex.GetFromCursor(playerCursor))
		cam := cli.Cameras()[camIndex]
		// Get the cameras local scene position
		_, cameraScenePosition := cam.Positions()
		centerX := float64(cam.Surface().Bounds().Dx()) / 2
		centerY := float64(cam.Surface().Bounds().Dy()) / 2

		// The key change: calculate the centered positions for BOTH player and camera
		// for proper deadzone comparison
		centeredPlayerX := playerPos.X
		centeredPlayerY := playerPos.Y
		centeredCameraX := cameraScenePosition.X + centerX
		centeredCameraY := cameraScenePosition.Y + centerY

		// Calculate distance between camera center and player position
		diffX := centeredPlayerX - centeredCameraX
		diffY := centeredPlayerY - centeredCameraY

		// Apply deadzone - camera only moves when player is outside of deadzone
		deadzoneX := 60.0 // horizontal deadzone in pixels
		deadzoneY := 60.0 // vertical deadzone in pixels

		// Target position starts at current camera position
		targetX := cameraScenePosition.X
		targetY := cameraScenePosition.Y

		// Only move camera if player outside deadzone
		if math.Abs(diffX) > deadzoneX {
			// Adjust target position to keep player at edge of deadzone
			if diffX > 0 {
				targetX = centeredPlayerX - centerX - deadzoneX
			} else {
				targetX = centeredPlayerX - centerX + deadzoneX
			}
		}

		if math.Abs(diffY) > deadzoneY {
			// Adjust target position to keep player at edge of deadzone
			if diffY > 0 {
				targetY = centeredPlayerY - centerY - deadzoneY
			} else {
				targetY = centeredPlayerY - centerY + deadzoneY
			}
		}

		// Apply smooth lerping to camera movement
		cameraScenePosition.X = lerp(cameraScenePosition.X, targetX, 0.02)
		cameraScenePosition.Y = lerp(cameraScenePosition.Y, targetY, 0.04)

		// Lock the camera to the scene boundaries
		lockCameraToSceneBoundaries(cam, scene, cameraScenePosition)
	}
	return nil
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
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
