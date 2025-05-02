package main

import (
	"bytes"
	_ "image/png"
	"log"

	"github/actor"

	"utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rendering"
)

type GameState string

var (
	arcadeFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s
}

// Constants defining the allowed game states
const (
	GameMenu    GameState = "Menu"
	GameStarted GameState = "Started"
	GamePaused  GameState = "Paused"
	GameEnded   GameState = "Ended"
)

const sampleText = "Press space key to start"

type Game struct {
	scourgeActor *actor.Actor
	purgerActor  *actor.Actor
	NPCActors    []*actor.Actor
	State        GameState
	keys         []ebiten.Key
}

const (
	screenWidth  = 900
	screenHeight = 400
	lineSpacing  = 0
	fontSize     = 10
)

func newGame(playerActor *actor.Actor, npcActors []*actor.Actor) *Game {
	return &Game{
		State:        GameMenu,
		scourgeActor: playerActor,
		purgerActor:  playerActor,
		NPCActors:    npcActors,
	}
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	switch g.State {
	case GameMenu:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State = GameStarted
		}
	case GamePaused:

	case GameStarted:
		for npcActor := range g.NPCActors {
			g.NPCActors[npcActor].Patrol(10)
		}
		// g.NPCActors[0].Patrol(10)
		if len(g.keys) > 0 {
			g.purgerActor.HandleInput(g.keys)
			if g.purgerActor.CollidesWith(g.NPCActors[0]) {
				g.purgerActor.RollbackPosition()
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: cache images?
	switch g.State {
	case GameMenu:
		op := &text.DrawOptions{}
		op.LineSpacing = lineSpacing

		width, height := text.Measure(sampleText, &text.GoTextFace{
			Source: arcadeFaceSource,
			Size:   fontSize,
		}, lineSpacing)

		op.GeoM.Reset()
		op.GeoM.Translate(screenWidth/2-(width/2), height*5)
		text.Draw(screen, sampleText, &text.GoTextFace{
			Source: arcadeFaceSource,
			Size:   fontSize,
		}, op)
	case GameStarted:
		for _, npc := range g.NPCActors {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(npc.Position[0], npc.Position[1])
			rendering.DrawImageWithMatrix(screen, npc.Image, op)
		}

		opArthas := &ebiten.DrawImageOptions{}
		opArthas.GeoM.Translate(g.purgerActor.Position[0], g.purgerActor.Position[1])

		rendering.DrawImageWithMatrix(screen, g.purgerActor.Image, opArthas)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 400
}

func main() {
	// TODO: proper error handling
	scourgeTexture := rendering.CreateTexture(utils.LoadFile("./assets/pudge.png"))
	scourgeTexture1 := rendering.CreateTexture(utils.LoadFile("./assets/scourge.png"))
	purgerTexture := rendering.CreateTexture(utils.LoadFile("./assets/arthas.png"))

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")

	playerActor := actor.NewActor([2]float64{0, 0}, purgerTexture, 14)
	npcActor := actor.NewActor([2]float64{200, 200}, scourgeTexture, 4)
	npcActor1 := actor.NewActor([2]float64{300, 200}, scourgeTexture1, 1)
	npcActor2 := actor.NewActor([2]float64{100, 200}, scourgeTexture1, 1)
	npcActor3 := actor.NewActor([2]float64{250, 50}, scourgeTexture1, 1)
	npcActor4 := actor.NewActor([2]float64{300, 100}, scourgeTexture1, 1)
	npcActor5 := actor.NewActor([2]float64{300, 300}, scourgeTexture1, 0.5)

	if err := ebiten.RunGame(newGame(playerActor, []*actor.Actor{npcActor, npcActor1, npcActor2, npcActor3, npcActor4, npcActor5})); err != nil {
		log.Fatal(err)
	}
}
