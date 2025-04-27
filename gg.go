package main

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	count        int
	scourgeActor Actor
	purgerActor  Actor
}

type Actor struct {
	position       [2]float64
	targetPosition [2]float64
	image          *ebiten.Image
}

func (actor *Actor) Move(position [2]float64) {
	drawOptions := &ebiten.DrawImageOptions{}
	translationMatrix := drawOptions.GeoM
	actor.position = position
	translationMatrix.Translate(position[0], position[1])
}

var (
	scourgeImage *ebiten.Image
	purgerImage  *ebiten.Image
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
)

func getRandomNumInRange(limit float64) float64 {
	return 0 + rand.Float64()*(limit-0)
}

func (g *Game) Update() error {
	// if it hasn't reached the target position
	if g.purgerActor.targetPosition[0] == g.purgerActor.position[0] &&
		g.purgerActor.targetPosition[1] == g.purgerActor.position[1] {
		targetposX := getRandomNumInRange(screenWidth)
		targetposY := getRandomNumInRange(screenHeight)

		log.Print("posx", targetposX)
		log.Print("posy", targetposY)
		g.purgerActor.targetPosition = [2]float64{targetposX, targetposY}
	}

	// TODO: set move speed
	newPosition := [2]float64{g.purgerActor.position[0] + 0.1, g.purgerActor.position[1] + 0.1}
	g.purgerActor.Move(newPosition)

	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	opArthas := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(0.5, 0.5)
	op.GeoM.Translate(100, 100)

	opArthas.GeoM.Scale(0.2, 0.2)
	opArthas.GeoM.Translate(g.purgerActor.position[0], g.purgerActor.position[1])

	drawImageWithMatrix(screen, scourgeImage, op)
	drawImageWithMatrix(screen, purgerImage, opArthas)
}

func drawImageWithMatrix(screen *ebiten.Image, image *ebiten.Image, transformationM *ebiten.DrawImageOptions) {
	// with animation
	// i := (g.count / 20) % frameCount
	// i := 1
	// log.Print(i)
	// sx, sy := frameOX+i*frameWidth, frameOY
	// screen.DrawImage(scourgeImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	// screen.DrawImage(scourgeImage.SubImage(image.Rect(60, 120, 90, 150)).(*ebiten.Image), op)

	screen.DrawImage(image, transformationM)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func loadFile(path string) *os.File {
	file, err := os.Open(path) // Path to your image
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func main() {
	scourgeTexture := loadFile("pudge.png")
	purgerTexture := loadFile("purger9000.PNG")

	img, _, err := image.Decode(scourgeTexture)
	if err != nil {
		log.Fatal(err)
	}

	scourgeImage = ebiten.NewImageFromImage(img)

	art, _, err := image.Decode(purgerTexture)
	if err != nil {
		log.Fatal(err)
	}

	purgerImage = ebiten.NewImageFromImage(art)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
