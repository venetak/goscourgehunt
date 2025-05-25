package main

import (
	_ "image/png"
	"log"

	"github.com/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	// TODO: proper error handling

	// npcActor4 := actor.NewActor([2]float64{300, 100}, scourgeTexture1, 1)
	// npcActor5 := actor.NewActor([2]float64{300, 300}, scourgeTexture1, 0.5)

	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
