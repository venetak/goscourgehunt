package main

import (
	_ "image/png"
	"log"

	"github/actor"

	"github/utils"

	"github.com/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rendering"
)

func main() {
	// TODO: proper error handling
	// move this to init function of game?
	scourgeTexture := rendering.CreateTexture(utils.LoadFile("./assets/pudge.png"))
	scourgeTexture1 := rendering.CreateTexture(utils.LoadFile("./assets/scourge.png"))
	purgerTexture := rendering.CreateTexture(utils.LoadFile("./assets/arthas.png"))

	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")

	playerActor := actor.NewActor([2]float64{0, 0}, purgerTexture, 14, "Purger")
	npcActor := actor.NewActor([2]float64{200, 200}, scourgeTexture, 4, "Scourge")
	npcActor1 := actor.NewActor([2]float64{400, 200}, scourgeTexture1, 1, "Undead1")
	npcActor2 := actor.NewActor([2]float64{500, 300}, scourgeTexture1, 1, "Undead2")
	npcActor3 := actor.NewActor([2]float64{250, 50}, scourgeTexture1, 1, "Undead3")
	// npcActor4 := actor.NewActor([2]float64{300, 100}, scourgeTexture1, 1)
	// npcActor5 := actor.NewActor([2]float64{300, 300}, scourgeTexture1, 0.5)

	if err := ebiten.RunGame(game.NewGame(playerActor, []*actor.Actor{npcActor, npcActor1, npcActor2, npcActor3})); err != nil {
		log.Fatal(err)
	}
}
