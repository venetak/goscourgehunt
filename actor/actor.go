package actor

import (
	"log"
	"math"

	"github/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Actor struct {
	Position        [2]float64
	initialPosition [2]float64
	targetPosition  [2]float64
	Image           *ebiten.Image
	Speed           float64
	MoveDirectionX  float64
	MoveDirectionY  float64
	moveRange       float64
}

func NewActor(position [2]float64, image *ebiten.Image, speed float64) *Actor {

	return &Actor{
		Position:        position,
		initialPosition: position,
		targetPosition:  position,
		Image:           image,
		Speed:           speed, // Default speed
		MoveDirectionX:  0.0,   // Default direction
		MoveDirectionY:  0.0,
		moveRange:       100.0, // Default move range
	}
}

// HandleInput processes the input from the user and updates the actor's movement direction accordingly.
// It takes a slice of ebiten.Key values representing the keys currently being pressed.
// The function resets the actor's movement direction, checks which keys are pressed, and updates
// the movement direction (MoveDirectionX and MoveDirectionY) based on the arrow keys.
// Finally, it calculates the new position and moves the actor in the specified direction.
func (actor *Actor) HandleInput(inputDeviceActionButtonNames []ebiten.Key) {
	actor.resetMoveDirection()

	for _, key := range inputDeviceActionButtonNames {
		if ebiten.IsKeyPressed(key) {
			switch key {
			case ebiten.KeyArrowLeft:
				actor.MoveDirectionX = -1
			case ebiten.KeyArrowRight:
				actor.MoveDirectionX = 1
			case ebiten.KeyArrowUp:
				actor.MoveDirectionY = -1
			case ebiten.KeyArrowDown:
				actor.MoveDirectionY = 1
			}
		}
	}

	newPosition := [2]float64{actor.MoveDirectionX, actor.MoveDirectionY}
	actor.MoveIn(newPosition)
}

func (actor *Actor) resetMoveDirection() {
	actor.MoveDirectionX = 0
	actor.MoveDirectionY = 0
}

func (actor *Actor) RollbackPosition() {
	actor.MoveIn([2]float64{-actor.MoveDirectionX, -actor.MoveDirectionY})
}

func (actor *Actor) MoveIn(direction [2]float64) {
	dx := direction[0]
	dy := direction[1]

	// Pythagorean theorem for the hypotenuse od a right triangle
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	if distance > 0 {
		// normalize the direction vector to get smooth movement on the diagonals
		dx = (dx / distance) * actor.Speed
		dy = (dy / distance) * actor.Speed

		actor.Position[0] += dx
		actor.Position[1] += dy
	}
}

type BoundingRect struct {
	PositionX float64
	PositionY float64
	Width     float64
	Height    float64
}

func calcBoundingRect(actor *Actor) *BoundingRect {
	playerRect := actor.Image.Bounds()

	return &BoundingRect{
		PositionX: actor.Position[0],
		PositionY: actor.Position[1],
		Width:     float64(playerRect.Dx()),
		Height:    float64(playerRect.Dy()),
	}
}

// AABB (Axis-Aligned Bounding Box) Collision detection
// Checks if two actors are colliding based on their bounding rectangles.
// - npc (non-player character) The second actor
// returns true if they are colliding, and false otherwise.
func (actor *Actor) CollidesWith(npc *Actor) bool {
	playerRect := calcBoundingRect(actor)
	npcRect := calcBoundingRect(npc)

	if playerRect.PositionX < npcRect.PositionX+npcRect.Width &&
		playerRect.PositionX+playerRect.Width > npcRect.PositionX &&
		playerRect.PositionY < npcRect.PositionY+npcRect.Height &&
		playerRect.PositionY+playerRect.Height > npcRect.PositionY {
		log.Print("Collision!!!!!--------")
		return true
	}

	return false
}

// Moves the actor towards the target position.
// Calculates the angle between the positive x-axis and the line connecting the actor's current position to the target position
// todetirmine the direction of movement.
func (actor *Actor) MoveTo(targetPosition [2]float64) {
	x := actor.Position[0]
	y := actor.Position[1]

	dx := targetPosition[0] - x
	dy := targetPosition[1] - y

	teta := math.Atan2(dy, dx)

	nextPositionX := actor.Position[0] + math.Cos(teta)*actor.Speed
	nextPositionY := actor.Position[1] + math.Sin(teta)*actor.Speed

	actor.Position[0] = nextPositionX
	actor.Position[1] = nextPositionY

	// if x or y oversteps the target position, set it to the target position (clamping)
	if dx > 0 && nextPositionX > actor.targetPosition[0] {
		actor.Position[0] = actor.targetPosition[0]
	}
	if dx < 0 && nextPositionX < actor.targetPosition[0] {
		actor.Position[0] = actor.targetPosition[0]
	}
	if dy > 0 && nextPositionY > actor.targetPosition[1] {
		actor.Position[1] = actor.targetPosition[1]
	}
	if dy < 0 && nextPositionY < actor.targetPosition[1] {
		actor.Position[1] = actor.targetPosition[1]
	}
}

// Initiates a patrol movement for the actor.
// It patrols between the initial position and a random target position within a specified move range.
// If the actor reaches the target position, it generates a new random target position within the move range.
func (actor *Actor) Patrol(moveRange float64) {
	if actor.Position[0] == actor.targetPosition[0] &&
		actor.Position[1] == actor.targetPosition[1] {
		actor.targetPosition[0] = utils.GetRandomNumInRange(actor.initialPosition[0]-actor.moveRange, actor.initialPosition[0]+actor.moveRange)
		actor.targetPosition[1] = utils.GetRandomNumInRange(actor.initialPosition[1]-actor.moveRange, actor.initialPosition[1]+actor.moveRange)
	}

	actor.MoveTo(actor.targetPosition)
}
