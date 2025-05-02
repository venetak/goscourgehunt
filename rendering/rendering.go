package rendering

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
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
