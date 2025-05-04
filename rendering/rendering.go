package rendering

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

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
func DrawCenteredText(screen *ebiten.Image, textToDraw string, x, y float64, faceSource *text.GoTextFaceSource, fontSize float64) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.LayoutOptions.PrimaryAlign = text.AlignCenter
	op.LayoutOptions.SecondaryAlign = text.AlignCenter

	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: faceSource,
		Size:   fontSize,
	}, op)
}
