package actor

import (
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Actor struct {
	Position          [2]float64
	prevFramePosition [2]float64
	targetPosition    [2]float64
	Image             *ebiten.Image
	Speed             float64
	MoveDirectionX    float64
	MoveDirectionY    float64
}

func NewActor(position [2]float64, image *ebiten.Image, speed float64) *Actor {
	return &Actor{
		Position:          position,
		prevFramePosition: position,
		Image:             image,
		Speed:             speed, // Default speed
		MoveDirectionX:    0.0,   // Default direction
		MoveDirectionY:    0.0,
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

func (actor *Actor) snapshotCurrentPosition() {
	actor.prevFramePosition = actor.Position
}

func (actor *Actor) RollbackPosition() {
	actor.Position = actor.prevFramePosition
}

func (actor *Actor) MoveIn(direction [2]float64) {
	actor.snapshotCurrentPosition()

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

// Utils ---- move to module?
func getRandomNumInRange(limit float64) float64 {
	return 0 + rand.Float64()*(limit-0)
}

func (actor *Actor) MoveTo(targetPosition [2]float64) {
	x := actor.Position[0]
	y := actor.Position[1]

	dx := targetPosition[0] - x
	dy := targetPosition[1] - y

	teta := math.Atan2(dy, dx)

	actor.Position[0] += math.Cos(teta) * actor.Speed
	actor.Position[1] += math.Sin(teta) * actor.Speed
}

func (actor *Actor) Patrol() {
	// TODO: propertly clamp the target position, do not use whole numbers with whole number step to avoid overstepping...
	if int(actor.Position[0]) == int(actor.targetPosition[0]) &&
		int(actor.Position[1]) == int(actor.targetPosition[1]) {
		actor.targetPosition[0] = getRandomNumInRange(400)
		actor.targetPosition[1] = getRandomNumInRange(300)
	} else {
		log.Print("Patrol------------------")
		log.Print("actor.Position[0]=======")
		log.Print(actor.Position[0])
		log.Print("actor.Position[1]=======")
		log.Print(actor.Position[1])
	}

	actor.MoveTo(actor.targetPosition)
}
