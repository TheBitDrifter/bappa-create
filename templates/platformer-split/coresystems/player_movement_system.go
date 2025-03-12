package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappacreate/templates/platformer-split/actions"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-split/components"
	"github.com/TheBitDrifter/blueprint"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/warehouse"
)

// PlayerMovementSystem handles all player movement mechanics including horizontal
// movement on flat ground and slopes, jumping with coyote time and input buffering,
// and platform drop-through functionality.
type PlayerMovementSystem struct{}

// Run executes all movement systems in the correct order:
// 1. Horizontal movement (including slope interactions)
// 2. Jumping (with coyote time and input buffering)
// 3. Down input handling (for platform drop-through)
func (sys PlayerMovementSystem) Run(scene blueprint.Scene, dt float64) error {
	sys.handleHorizontal(scene)
	sys.handleJump(scene)
	return sys.handleDown(scene)
}

// handleHorizontal processes left/right movement with different behaviors for:
// - Air movement
// - Flat ground movement
// - Uphill/downhill slope movement with proper tangent calculations
func (PlayerMovementSystem) handleHorizontal(scene blueprint.Scene) {
	// Constants
	const (
		speedX    = 120.0 // Player's horizontal movement speed
		snapForce = 40.0  // Downward force to keep player attached to slopes
	)

	// Query all entities with input buffers
	cursor := scene.NewCursor(blueprint.Queries.InputBuffer)
	for cursor.Next() {
		// --- Gather required components ---
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)              // Physics properties
		incomingInputs := blueprintinput.Components.InputBuffer.GetFromCursor(cursor) // User inputs
		direction := blueprintspatial.Components.Direction.GetFromCursor(cursor)      // Facing direction

		// --- Process left/right inputs ---
		// Check and consume directional inputs
		_, pressedLeft := incomingInputs.ConsumeInput(actions.Left)
		if pressedLeft {
			direction.SetLeft()
		}

		_, pressedRight := incomingInputs.ConsumeInput(actions.Right)
		if pressedRight {
			direction.SetRight()
		}

		// Track if player is attempting to move horizontally this frame
		isMovingHorizontal := pressedLeft || pressedRight

		// --- Check ground status ---
		// First check if the OnGroundComponent exists and get its value safely
		isGroundComponentPresent, onGround := components.OnGroundComponent.GetFromCursorSafe(cursor)

		// Default to airborne movement if no ground component exists
		if !isGroundComponentPresent {
			// Simple air movement - just move in pressed direction or stop
			if isMovingHorizontal {
				dyn.Vel.X = speedX * direction.AsFloat()
			} else {
				dyn.Vel.X = 0
			}
			continue
		}

		currentTick := scene.CurrentTick()
		// Check if player touched ground on the previous tick (doesn't guarantee current grounded state)
		touchedGroundLastTick := currentTick-1 == onGround.LastTouch
		ticksOnGround := currentTick - onGround.LastJump

		// --- Ground snapping logic ---
		// Apply small downward force to keep player attached to slopes when grounded
		// Only applies if player has been on ground for a while and touched ground last tick
		if isGroundComponentPresent && ticksOnGround > 10 && touchedGroundLastTick {
			dyn.Vel.Y = math.Max(dyn.Vel.Y, snapForce) // Apply downward force
		}

		// --- Handle in-air movement ---
		// When in air or just landed this tick, use simplified movement
		if !isGroundComponentPresent || !touchedGroundLastTick {
			if isMovingHorizontal {
				dyn.Vel.X = speedX * direction.AsFloat() // direction.AsFloat() returns -1 for left, 1 for right
			} else {
				dyn.Vel.X = 0 // No horizontal movement when no keys pressed
			}
			continue // Skip slope handling for airborne players
		}

		// --- Handle flat ground movement ---
		// Check if player is on a flat surface (normal pointing straight up)
		flat := onGround.SlopeNormal.X == 0 && onGround.SlopeNormal.Y == 1
		if flat {
			// Same as air movement on flat ground
			if isMovingHorizontal {
				dyn.Vel.X = speedX * direction.AsFloat()
			} else {
				dyn.Vel.X = 0
			}
			continue // Skip slope handling
		}

		// --- Handle slope movement ---
		if isMovingHorizontal {
			// Calculate tangent vector along the slope
			// The tangent is perpendicular to the normal, so we swap X/Y and negate Y
			tangent := vector.Two{X: onGround.SlopeNormal.Y, Y: -onGround.SlopeNormal.X}

			// Determine if player is moving uphill by checking if direction and normal X have same sign
			isUphill := (direction.AsFloat() * onGround.SlopeNormal.X) > 0

			// Scale tangent by movement direction for correct slope alignment
			slopeDir := tangent.Scale(direction.AsFloat())

			if isUphill {
				// When going uphill, only set X velocity and let physics handle Y
				dyn.Vel.X = slopeDir.X * speedX
			} else {
				// When going downhill, help player follow the slope with both X and Y velocities
				dyn.Vel.X = slopeDir.X * speedX

				// Only apply downward velocity after being on ground for a while
				// This prevents immediate sliding when just landing on a slope
				if ticksOnGround > 10 {
					dyn.Vel.Y = slopeDir.Y * speedX
				}
			}
		} else {
			// No movement if not pressing direction keys
			dyn.Vel.X = 0
		}
	}
}

// handleJump processes jump inputs with coyote time and input buffering features
// Coyote time: Player can jump shortly after leaving a platform
// Input buffering: Jump inputs are remembered and applied when landing
func (PlayerMovementSystem) handleJump(scene blueprint.Scene) {
	// Constants
	const (
		jumpForce        = 320.0 // Upward force applied when jumping
		coyoteTimeTicks  = 10    // Ticks after leaving ground where jump is still allowed
		inputBufferTicks = 5     // Ticks before landing where a jump input is remembered
	)

	// Create query for players eligible to jump (have ground and input components)
	playersEligibleToJumpQuery := warehouse.Factory.NewQuery()
	playersEligibleToJumpQuery.And(components.OnGroundComponent, blueprintinput.Components.InputBuffer)

	// Get all entities that match the query
	cursor := scene.NewCursor(playersEligibleToJumpQuery)
	currentTick := scene.CurrentTick()

	for cursor.Next() {
		// Get required components
		dyn := blueprintmotion.Components.Dynamics.GetFromCursor(cursor)
		incomingInputs := blueprintinput.Components.InputBuffer.GetFromCursor(cursor)

		// OnGroundComponent is guaranteed to exist because of our query
		onGround := components.OnGroundComponent.GetFromCursor(cursor)

		// Check for jump input
		if stampedInput, inputReceived := incomingInputs.ConsumeInput(actions.Jump); inputReceived {
			// --- Jump Eligibility Checks ---

			// Coyote time: Allow jumping within certain ticks of leaving ground
			playerGroundedWithinCoyoteTime := currentTick-onGround.LastTouch <= coyoteTimeTicks

			// Input buffering checks:
			// 1. Was input received before touching ground?
			jumpInputIsBeforeGroundTouch := stampedInput.Tick <= onGround.LastTouch
			// 2. Was input within the buffer window?
			jumpInputWithinBufferWindow := onGround.LastTouch-stampedInput.Tick <= inputBufferTicks
			// Combined buffer condition
			validBufferedJumpInput := jumpInputIsBeforeGroundTouch && jumpInputWithinBufferWindow

			// Direct jump: Input received while already on ground
			directJumpInput := stampedInput.Tick >= onGround.LastTouch

			// Prevent double jumps: Make sure player hasn't jumped since last touching ground
			playerHasNotJumpedSinceGroundTouch := onGround.LastJump < onGround.LastTouch

			// Player can jump if:
			// 1. They haven't jumped since touching ground, AND
			// 2a. They are in coyote time with a direct input, OR
			// 2b. They have a valid buffered input from before landing
			canJump := playerHasNotJumpedSinceGroundTouch &&
				((playerGroundedWithinCoyoteTime && directJumpInput) || validBufferedJumpInput)

			if canJump {
				// Apply upward velocity and acceleration for jump
				dyn.Vel.Y = -jumpForce
				dyn.Accel.Y = -jumpForce
				// Record jump time
				onGround.LastJump = currentTick
			}
		}
	}
}

// handleDown processes down input for platform drop-through functionality
// This allows players to press down to fall through one-way platforms
func (PlayerMovementSystem) handleDown(scene blueprint.Scene) error {
	// Create query for players eligible to drop (have ground and input components)
	playersEligibleToDropQuery := warehouse.Factory.NewQuery()
	playersEligibleToDropQuery.And(components.OnGroundComponent, blueprintinput.Components.InputBuffer)

	cursor := scene.NewCursor(playersEligibleToDropQuery)
	currentTick := scene.CurrentTick()

	for cursor.Next() {
		// OnGroundComponent is guaranteed to exist because of our query
		onGround := components.OnGroundComponent.GetFromCursor(cursor)

		// The component exists but were not actually grounded (will be removed soon)
		// Likely being held onto for coyote time tracking
		if onGround.LastTouch != currentTick-1 {
			continue
		}

		// Don't allow platform drop if player just jumped this tick
		if onGround.LastJump == currentTick {
			continue
		}

		// Get player entity
		playerEntity, err := cursor.CurrentEntity()
		if err != nil {
			return err
		}

		// Get input component
		incomingInputs := blueprintinput.Components.InputBuffer.GetFromCursor(cursor)

		// Check for down action
		if stampedInput, inputReceived := incomingInputs.ConsumeInput(actions.Down); inputReceived {
			// Only process inputs from this tick (ignore buffered inputs)
			if stampedInput.Tick != currentTick {
				continue
			}

			// Otherwise, add the IgnorePlatform component to allow dropping through platforms
			err := playerEntity.EnqueueAddComponent(components.IgnorePlatformComponent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
