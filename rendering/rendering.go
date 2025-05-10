package rendering

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	faceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	faceSource = s
}

// Rendering utils --- module?
func DrawImageWithMatrix(screen *ebiten.Image, image *ebiten.Image, transformationM *ebiten.DrawImageOptions) {
	screen.DrawImage(image, transformationM)
}

func CreateTexture(imageFile *os.File) *ebiten.Image {
	img, _, err := image.Decode(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(img)
}

func DrawBox(screen *ebiten.Image, x, y, width, height float32) {
	bgColor := color.RGBA{0xFF, 0x00, 0x00, 0xFF} // Red background (like background-color)
	borderColor := color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	borderWidth := float32(2) // Border thickness (like border-width)

	// Draw filled rectangle for background.
	vector.DrawFilledRect(screen, x, y, width, height, bgColor, false)

	// Draw border (stroked rectangle) around the box.
	// Adjust position and size to account for border thickness if you want the border to be outside.
	vector.StrokeRect(screen, x, y, width, height, borderWidth, borderColor, false)
}

// Do I need this?
func DrawCenteredText(screen *ebiten.Image, textToDraw string, x, y float64, fontSize float64) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	op.LayoutOptions.SecondaryAlign = text.AlignCenter

	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: faceSource,
		Size:   fontSize,
	}, op)
}

// TODO: consider reducing the number of parameters
func DrawPlayerPromptAtActorPos(screen *ebiten.Image, textToDraw string, actorPos [2]float64, fontSize float64) {
	// Draw the text at the actor's position
	posX := actorPos[0]
	posY := actorPos[1] - 10
	// DrawBox(screen, float32(posX), float32(posY-10/2), 350, 20)

	DrawText(screen, textToDraw, posX, posY, fontSize)
}

func DrawText(screen *ebiten.Image, textToDraw string, x, y float64, fontSize float64) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)

	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: faceSource,
		Size:   fontSize,
	}, op)
}
