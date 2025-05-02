package main

import (
	_ "image/png"
	"log"

	"github/actor"

	"utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rendering"
)

type Game struct {
	scourgeActor *actor.Actor
	purgerActor  *actor.Actor
	NPCActors    []*actor.Actor
	State        string
	keys         []ebiten.Key
}

const (
	screenWidth  = 1000
	screenHeight = 400
)

func newGame(playerActor *actor.Actor, npcActors []*actor.Actor) *Game {
	return &Game{
		scourgeActor: playerActor,
		purgerActor:  playerActor,
		NPCActors:    npcActors,
	}
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	for npcActor := range g.NPCActors {
		g.NPCActors[npcActor].Patrol()
	}
	g.NPCActors[0].Patrol()
	if len(g.keys) > 0 {
		g.purgerActor.HandleInput(g.keys)
		// if g.purgerActor.CollidesWith(g.NPCActors[0]) {
		// 	g.purgerActor.RollbackPosition()
		// }
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: cache images?
	for _, npc := range g.NPCActors {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(npc.Position[0], npc.Position[1])
		rendering.DrawImageWithMatrix(screen, npc.Image, op)
	}

	opArthas := &ebiten.DrawImageOptions{}
	opArthas.GeoM.Translate(g.purgerActor.Position[0], g.purgerActor.Position[1])

	rendering.DrawImageWithMatrix(screen, g.purgerActor.Image, opArthas)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 400
}

func main() {
	// TODO: proper error handling
	scourgeTexture := rendering.CreateTexture(utils.LoadFile("./assets/pudge.png"))
	scourgeTexture1 := rendering.CreateTexture(utils.LoadFile("./assets/scourge.png"))
	purgerTexture := rendering.CreateTexture(utils.LoadFile("./assets/arthas.png"))

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")

	playerActor := actor.NewActor([2]float64{0, 0}, purgerTexture, 14)
	npcActor := actor.NewActor([2]float64{200, 200}, scourgeTexture, 0.5)
	npcActor1 := actor.NewActor([2]float64{200, 200}, scourgeTexture1, 4)
	npcActor2 := actor.NewActor([2]float64{100, 200}, scourgeTexture1, 1)
	npcActor3 := actor.NewActor([2]float64{150, 50}, scourgeTexture1, 1)
	npcActor4 := actor.NewActor([2]float64{200, 200}, scourgeTexture1, 2)
	npcActor5 := actor.NewActor([2]float64{200, 200}, scourgeTexture1, 0.5)

	if err := ebiten.RunGame(newGame(playerActor, []*actor.Actor{npcActor, npcActor1, npcActor2, npcActor3, npcActor4, npcActor5})); err != nil {
		log.Fatal(err)
	}
}
