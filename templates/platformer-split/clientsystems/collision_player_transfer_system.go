package clientsystems

import (
	"github.com/TheBitDrifter/bappacreate/templates/platformer-split/components"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/coldbrew"
	"github.com/TheBitDrifter/tteokbokki/spatial"
	"github.com/TheBitDrifter/warehouse"
)

type CollisionPlayerTransferSystem struct{}

type playerTransfer struct {
	target       string
	playerEntity warehouse.Entity
}

// System handles transfering players when they collide with PlayerSceneTransfer entities
func (CollisionPlayerTransferSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	// We transfer after the loop
	// So we collect pending here
	var pending []playerTransfer

	// Query the transfer collision entities
	collisionTransferQuery := warehouse.Factory.NewQuery().And(
		blueprintspatial.Components.Shape,
		components.PlayerSceneTransferComponent,
	)

	// Query player(s)
	playerWithShapeQuery := warehouse.Factory.NewQuery().And(
		blueprintspatial.Components.Shape,
		blueprintinput.Components.InputBuffer,
	)

	collisionTransferCursor := scene.NewCursor(collisionTransferQuery)
	playerWithShapeCursor := scene.NewCursor(playerWithShapeQuery)

	// Iterate through pairs (outer loop)
	for range collisionTransferCursor.Next() {

		// Get pos and collider/shape
		transferPos := blueprintspatial.Components.Position.GetFromCursor(collisionTransferCursor)
		transferCollider := blueprintspatial.Components.Shape.GetFromCursor(collisionTransferCursor)

		// Inner loop
		for range playerWithShapeCursor.Next() {

			// Get pos and collider/shape
			playerPos := blueprintspatial.Components.Position.GetFromCursor(playerWithShapeCursor)
			playerCollider := blueprintspatial.Components.Shape.GetFromCursor(playerWithShapeCursor)

			// Check for collision
			if ok, _ := spatial.Detector.Check(*playerCollider, *transferCollider, playerPos, transferPos); ok {
				sceneTransfer := components.PlayerSceneTransferComponent.GetFromCursor(collisionTransferCursor)
				playerEn, err := playerWithShapeCursor.CurrentEntity()
				if err != nil {
					return err
				}

				// If collision enqueue transfers
				transfer := playerTransfer{
					target:       sceneTransfer.Dest,
					playerEntity: playerEn,
				}
				pending = append(pending, transfer)
				// Update the player pos
				playerPos := blueprintspatial.Components.Position.GetFromCursor(playerWithShapeCursor)
				playerPos.X = sceneTransfer.X
				playerPos.Y = sceneTransfer.Y

				// Update the camera pos
				camIndex := int(*blueprintclient.Components.CameraIndex.GetFromCursor(playerWithShapeCursor))
				cam := cli.Cameras()[camIndex]
				// Get the cameras local scene position
				_, cameraScenePosition := cam.Positions()
				centerX := float64(cam.Surface().Bounds().Dx()) / 2
				centerY := float64(cam.Surface().Bounds().Dy()) / 2

				// Set position on target player
				cameraScenePosition.X = playerPos.X - centerX
				cameraScenePosition.Y = playerPos.Y - centerY

			}
		}

	}

	// Process transfers after loop
	for _, transfer := range pending {
		_, err := cli.ActivateSceneByName(transfer.target, transfer.playerEntity)
		if err != nil {
			return err
		}
	}

	return nil
}
