package clientsystems

import (
	"github.com/TheBitDrifter/coldbrew"
	coldbrew_clientsystems "github.com/TheBitDrifter/coldbrew/clientsystems"
)

var DefaultClientSystems = []coldbrew.ClientSystem{
	PlayerSoundSystem{},
	MusicSystem{},
	PlayerAnimationSystem{},
	&CameraFollowerSystem{},
	&coldbrew_clientsystems.BackgroundScrollSystem{},
	CollisionPlayerTransferSystem{},
}
