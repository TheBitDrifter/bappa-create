package coresystems

import (
	"github.com/TheBitDrifter/blueprint"
	tteo_coresystems "github.com/TheBitDrifter/tteokbokki/coresystems"
)

var DefaultCoreSystems = []blueprint.CoreSystem{
	// If you want default physics:
	tteo_coresystems.IntegrationSystem{},
	tteo_coresystems.TransformSystem{},
}
