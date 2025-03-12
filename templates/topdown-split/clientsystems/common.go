package clientsystems

import (
	"github.com/TheBitDrifter/coldbrew"
)

var DefaultClientSystems = []coldbrew.ClientSystem{
	PlayerSoundSystem{},             // Player Sounds
	MusicSystem{},                   // Music
	PlayerAnimationSystem{},         // Player Animations
	CameraFollowerSystem{},          // Camera follows player
	CollisionPlayerTransferSystem{}, // Handles scene transfers
	SortVerticalSystem{},            // Sorts sprites based on Y positions
}
