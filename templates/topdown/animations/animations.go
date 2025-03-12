package animations

import (
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
)

var Up = blueprintclient.AnimationData{
	Name:        "up",
	RowIndex:    3,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}

var Down = blueprintclient.AnimationData{
	Name:        "down",
	RowIndex:    0,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}

var UpSide = blueprintclient.AnimationData{
	Name:        "upside",
	RowIndex:    4,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}

var DownSide = blueprintclient.AnimationData{
	Name:        "downside",
	RowIndex:    5,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}

// Note: downside and side are the same because I could only find a 6
// directional sprite. If we had an 8 directional sprite we could use a 4th
// animation here
var Side = blueprintclient.AnimationData{
	Name:        "side",
	RowIndex:    5,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}
