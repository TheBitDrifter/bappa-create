package rendersystems

import "github.com/TheBitDrifter/coldbrew"

var DefaultRenderSystems = []coldbrew.RenderSystem{
	PlayerCameraPriorityRenderer{},
}
