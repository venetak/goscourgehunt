package main

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"

	"github/actor"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	scourgeActor *actor.Actor
	purgerActor  *actor.Actor
	NPCActors    []*actor.Actor
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

// Utils ---- move to module?
func getRandomNumInRange(limit float64) float64 {
	return 0 + rand.Float64()*(limit-0)
}

// Rendering utils --- module?
func drawImageWithMatrix(screen *ebiten.Image, image *ebiten.Image, transformationM *ebiten.DrawImageOptions) {
	screen.DrawImage(image, transformationM)
}

func createTexture(imageFile *os.File) *ebiten.Image {
	img, _, err := image.Decode(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	if len(g.keys) > 0 {
		g.purgerActor.HandleInput(g.keys)
		if g.purgerActor.CollidesWith(g.NPCActors[0]) {
			g.purgerActor.RollbackPosition()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: cache images?
	op := &ebiten.DrawImageOptions{}
	opArthas := &ebiten.DrawImageOptions{}

	// op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(g.NPCActors[0].Position[0], g.NPCActors[0].Position[1])

	// opArthas.GeoM.Scale(0.2, 0.2)
	opArthas.GeoM.Translate(g.purgerActor.Position[0], g.purgerActor.Position[1])

	drawImageWithMatrix(screen, g.purgerActor.Image, opArthas)
	for _, npc := range g.NPCActors {
		drawImageWithMatrix(screen, npc.Image, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 400
}

func loadFile(path string) *os.File {
	file, err := os.Open(path) // Path to your image
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func main() {
	// TODO: proper error handling
	scourgeTexture := createTexture(loadFile("pudge.png"))
	purgerTexture := createTexture(loadFile("pudge.png"))

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")

	playerActor := actor.NewActor([2]float64{0, 0}, purgerTexture, 0)
	npcActor := actor.NewActor([2]float64{200, 200}, scourgeTexture, 0)

	if err := ebiten.RunGame(newGame(playerActor, []*actor.Actor{npcActor})); err != nil {
		log.Fatal(err)
	}
}
