package game

import (
	"bytes"
	_ "image/png"
	"log"
	"strconv"

	"github/actor"

	"github.com/gameplay"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rendering"
)

var (
	arcadeFaceSource *text.GoTextFaceSource
	StatusMap        = gameplay.StatusMap
	GameMenu         = gameplay.GameMenu
	GameStarted      = gameplay.GameStarted
	GamePaused       = gameplay.GamePaused
	GameEnded        = gameplay.GameEnded
	AwaitingUser     = gameplay.AwaitingUser
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 550
	LineSpacing  = 0
	FontSize     = 10
)

const ScreenWidthFloat = float64(ScreenWidth)
const ScreenHeightFloat = float64(ScreenHeight)

var GameModeMap = map[int]string{
	1: "Invincible",
	2: "Frostmourne Hungers",
}

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	arcadeFaceSource = s
}

const sampleText = "Press space key to start"

type Game struct {
	scourgeActor *actor.Actor
	purgerActor  *actor.Actor
	NPCActors    []*actor.Actor
	State        *gameplay.GameState
	keys         []ebiten.Key
	GameMode     int
	PlayMode     gameplay.PlayMode
	PromptPlayer bool
}

func NewGame(playerActor *actor.Actor, npcActors []*actor.Actor) *Game {
	return &Game{
		State:        &gameplay.GameState{Status: gameplay.StatusMap[gameplay.GameMenu]},
		scourgeActor: playerActor,
		purgerActor:  playerActor,
		NPCActors:    npcActors,
	}
}

func (g *Game) SpawnActors(screen *ebiten.Image) {
	for _, npc := range g.NPCActors {
		if !npc.Draw {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(npc.Position[0], npc.Position[1])
		rendering.DrawImageWithMatrix(screen, npc.Image, op)
	}

	opArthas := &ebiten.DrawImageOptions{}
	opArthas.GeoM.Translate(g.purgerActor.Position[0], g.purgerActor.Position[1])

	rendering.DrawImageWithMatrix(screen, g.purgerActor.Image, opArthas)
}

func (g *Game) PurgedCountStr() string {
	return "Purged: " + strconv.Itoa(g.State.PurgedCount)
}
func (g *Game) SparedCountStr() string {
	return "Spared: " + strconv.Itoa(g.State.SparedCount)
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	switch g.State.Status {
	case StatusMap[GameMenu]:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State.Status = StatusMap[GameStarted]
			g.GameMode = 1
			g.PlayMode = gameplay.NewPlayMode(g.GameMode)
		}
	case StatusMap[GamePaused]:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State.Status = StatusMap[GameStarted]
			return nil
		}
	case StatusMap[GameStarted]:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State.Status = StatusMap[GamePaused]
			return nil
		}
		// starts patrolling
		g.PlayMode.InitActors(g.purgerActor, g.NPCActors, g.keys)
		if len(g.keys) > 0 {
			g.PlayMode.HandleKeyboardInput(g.State, g.purgerActor, g.NPCActors, g.keys)
			g.purgerActor.SetLimitBounds(ScreenWidthFloat, ScreenHeightFloat)
		}
	case StatusMap[AwaitingUser]:
		g.PlayMode.HandlePlayerInput(g.State, g.NPCActors, g.State.Target)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Image at: (0, 0)"))
	// TODO: cache images?
	switch g.State.Status {
	case StatusMap[GameMenu]:
		rendering.DrawCenteredText(screen, sampleText, ScreenWidth/2, ScreenHeight/2, arcadeFaceSource, FontSize)
	case StatusMap[GameStarted]:
		g.SpawnActors(screen)
		rendering.DrawText(screen, g.PurgedCountStr(), ScreenWidth-100, 0, arcadeFaceSource, FontSize)
		rendering.DrawText(screen, g.SparedCountStr(), ScreenWidth-100, 20, arcadeFaceSource, FontSize)
	case StatusMap[GamePaused]:
		rendering.DrawText(screen, g.PurgedCountStr(), ScreenWidth-100, 0, arcadeFaceSource, FontSize)

		rendering.DrawText(screen, g.SparedCountStr(), ScreenWidth-100, 20, arcadeFaceSource, FontSize)
		g.SpawnActors(screen)
		choiceText := "Game Paused"
		rendering.DrawBox(screen, ScreenWidth/2-100, ScreenHeight/2-50, 200, 100)
		rendering.DrawCenteredText(screen, choiceText, ScreenWidth/2, ScreenHeight/2, arcadeFaceSource, FontSize)

		rendering.DrawText(screen, g.SparedCountStr(), ScreenWidth-100, 20, arcadeFaceSource, FontSize)
	case StatusMap[AwaitingUser]:
		rendering.DrawText(screen, g.PurgedCountStr(), ScreenWidth-100, 0, arcadeFaceSource, FontSize)

		rendering.DrawText(screen, g.SparedCountStr(), ScreenWidth-100, 20, arcadeFaceSource, FontSize)
		g.SpawnActors(screen)
		g.PlayMode.PauseGame(g.State, g.purgerActor, screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return ScreenWidth, ScreenHeight
}

// PlayMode interface
