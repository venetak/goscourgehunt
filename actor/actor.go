package actor

import "github.com/hajimehoshi/ebiten/v2"

type Actor struct {
	Position       [2]float64
	TargetPosition [2]float64
	image          *ebiten.Image
}

func (actor *Actor) Move(position [2]float64) {
	drawOptions := &ebiten.DrawImageOptions{}
	translationMatrix := drawOptions.GeoM
	actor.Position = position
	translationMatrix.Translate(position[0], position[1])
}
