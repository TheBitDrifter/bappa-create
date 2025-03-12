package components

// Direction constants using iota for automatic incrementation
const (
	DirectionUp uint8 = iota
	DirectionRight
	DirectionDown
	DirectionLeft
	DirectionRightUp
	DirectionRightDown
	DirectionLeftDown
	DirectionLeftUp
)

// DirectionEight represents 8-directions
type DirectionEight struct {
	val uint8
}

// Setter methods to change direction
func (d *DirectionEight) SetUp() {
	d.val = DirectionUp
}

func (d *DirectionEight) SetRight() {
	d.val = DirectionRight
}

func (d *DirectionEight) SetDown() {
	d.val = DirectionDown
}

func (d *DirectionEight) SetLeft() {
	d.val = DirectionLeft
}

func (d *DirectionEight) SetRightUp() {
	d.val = DirectionRightUp
}

func (d *DirectionEight) SetRightDown() {
	d.val = DirectionRightDown
}

func (d *DirectionEight) SetLeftDown() {
	d.val = DirectionLeftDown
}

func (d *DirectionEight) SetLeftUp() {
	d.val = DirectionLeftUp
}

// Constructor functions for each direction
func NewDirectionUp() DirectionEight {
	return DirectionEight{val: DirectionUp}
}

func NewDirectionRight() DirectionEight {
	return DirectionEight{val: DirectionRight}
}

func NewDirectionDown() DirectionEight {
	return DirectionEight{val: DirectionDown}
}

func NewDirectionLeft() DirectionEight {
	return DirectionEight{val: DirectionLeft}
}

func NewDirectionRightUp() DirectionEight {
	return DirectionEight{val: DirectionRightUp}
}

func NewDirectionRightDown() DirectionEight {
	return DirectionEight{val: DirectionRightDown}
}

func NewDirectionLeftDown() DirectionEight {
	return DirectionEight{val: DirectionLeftDown}
}

func NewDirectionLeftUp() DirectionEight {
	return DirectionEight{val: DirectionLeftUp}
}

// Helper methods to check direction
func (d DirectionEight) IsUp() bool {
	return d.val == DirectionUp
}

func (d DirectionEight) IsRight() bool {
	return d.val == DirectionRight
}

func (d DirectionEight) IsDown() bool {
	return d.val == DirectionDown
}

func (d DirectionEight) IsLeft() bool {
	return d.val == DirectionLeft
}

func (d DirectionEight) IsRightUp() bool {
	return d.val == DirectionRightUp
}

func (d DirectionEight) IsRightDown() bool {
	return d.val == DirectionRightDown
}

func (d DirectionEight) IsLeftDown() bool {
	return d.val == DirectionLeftDown
}

func (d DirectionEight) IsLeftUp() bool {
	return d.val == DirectionLeftUp
}

// Get string representation of direction
func (d DirectionEight) String() string {
	switch d.val {
	case DirectionUp:
		return "Up"
	case DirectionRight:
		return "Right"
	case DirectionDown:
		return "Down"
	case DirectionLeft:
		return "Left"
	case DirectionRightUp:
		return "RightUp"
	case DirectionRightDown:
		return "RightDown"
	case DirectionLeftDown:
		return "LeftDown"
	case DirectionLeftUp:
		return "LeftUp"
	default:
		return "Unknown"
	}
}
