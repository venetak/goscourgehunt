package actor

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Actor struct {
	Position       [2]float64
	TargetPosition [2]float64
	image          *ebiten.Image
	inputMap       map[string]string
	Speed          float64
	MoveDirectionX float64
	MoveDirectionY float64
}

func (actor *Actor) Move(position [2]float64) {
	drawOptions := &ebiten.DrawImageOptions{}
	translationMatrix := drawOptions.GeoM
	actor.Position = position
	translationMatrix.Translate(position[0], position[1])
}

func (actor *Actor) HandleInput(inputDeviceActionButtonNames []ebiten.Key) {
	// TODO: create maps for other input devices and merge them
	//
	actor.Speed = 1

	actor.MoveDirectionX = 0
	actor.MoveDirectionY = 0

	actionHandler := inputDeviceActionButtonNames[0].String()
	if actionHandler == "ArrowRight" {
		actor.MoveDirectionX = 1
	}
	if actionHandler == "ArrowUp" {
		actor.MoveDirectionY = -1
	}
	if actionHandler == "ArrowDown" {
		actor.MoveDirectionY = 1
	}
	if actionHandler == "ArrowLeft" {
		actor.MoveDirectionX = -1
	}

	if len(inputDeviceActionButtonNames) == 2 {
		secondPressedButton := inputDeviceActionButtonNames[1].String()
		if secondPressedButton == "ArrowRight" {
			actor.MoveDirectionX = 1
		}
		if secondPressedButton == "ArrowUp" {
			actor.MoveDirectionY = -1
		}
		if secondPressedButton == "ArrowDown" {
			actor.MoveDirectionY = 1
		}
		if secondPressedButton == "ArrowLeft" {
			actor.MoveDirectionX = -1
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

		log.Print("-----------------------")
		log.Print(dx)
		log.Print(dy)
	}
}

func (actor *Actor) MoveRight() {
	drawOptions := &ebiten.DrawImageOptions{}
	translationMatrix := drawOptions.GeoM
	newPosition := [2]float64{actor.Position[0] + 0.1, actor.Position[1] + 0.1}
	actor.Position = newPosition
	translationMatrix.Translate(newPosition[0], newPosition[1])
}
