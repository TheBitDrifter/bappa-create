package animations

import (
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	"github.com/TheBitDrifter/blueprint/vector"
)

var IdleAnimation = blueprintclient.AnimationData{
	Name:        "idle",
	RowIndex:    0,
	FrameCount:  6,
	FrameWidth:  144,
	FrameHeight: 116,
	Speed:       8,
}

var RunAnimation = blueprintclient.AnimationData{
	Name:        "run",
	RowIndex:    1,
	FrameCount:  8,
	FrameWidth:  144,
	FrameHeight: 116,
	Speed:       5,
}

var JumpAnimation = blueprintclient.AnimationData{
	Name:           "jump",
	RowIndex:       2,
	FrameCount:     3,
	FrameWidth:     144,
	FrameHeight:    116,
	Speed:          5,
	Freeze:         true,
	PositionOffset: vector.Two{X: 0, Y: 10},
}

var FallAnimation = blueprintclient.AnimationData{
	Name:           "fall",
	RowIndex:       3,
	FrameCount:     3,
	FrameWidth:     144,
	FrameHeight:    116,
	Speed:          5,
	Freeze:         true,
	PositionOffset: vector.Two{X: 0, Y: 10},
}
