package main

import (
	"fmt"
	_ "image/png"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	debugEnabled := os.Getenv("DEBUG")
	debug := false
	if debugEnabled == "TRUE" {
		debug = true
	}

	ebiten.SetWindowSize(game.ScreenWidth*2, game.ScreenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	// TODO: proper error handling

	// npcActor4 := actor.NewActor([2]float64{300, 100}, scourgeTexture1, 1)
	// npcActor5 := actor.NewActor([2]float64{300, 300}, scourgeTexture1, 0.5)

	if err := ebiten.RunGame(game.NewGame(debug)); err != nil {
		log.Fatal(err)
	}
}
