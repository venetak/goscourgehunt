package actor

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Actor struct {
	Position       [2]float64
	Image          *ebiten.Image
	Speed          float64
	MoveDirectionX float64
	MoveDirectionY float64
}

func NewActor(position [2]float64, image *ebiten.Image, speed float64) *Actor {
	return &Actor{
		Position:       position,
		Image:          image,
		Speed:          1.0, // Default speed
		MoveDirectionX: 0.0, // Default direction
		MoveDirectionY: 0.0,
	}
}

func (actor *Actor) HandleInput(inputDeviceActionButtonNames []ebiten.Key) {
	// TODO: create maps for other input devices and merge them
	//
	actor.Speed = 1

	actor.MoveDirectionX = 0
	actor.MoveDirectionY = 0

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
