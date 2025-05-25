package game

import (
	_ "image/png"
	"strconv"

	"github.com/actor"

	"github.com/gameplay"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/player"
	"github.com/rendering"
)

var (
	StatusMap    = gameplay.StatusMap
	GameMenu     = gameplay.GameMenu
	GameStarted  = gameplay.GameStarted
	GamePaused   = gameplay.GamePaused
	GameEnded    = gameplay.GameEnded
	AwaitingUser = gameplay.AwaitingUser
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 550
)

const ScreenWidthFloat = float64(ScreenWidth)
const ScreenHeightFloat = float64(ScreenHeight)

var GameModeMap = map[int]string{
	1: "Invincible",
	2: "Frostmourne Hungers",
}

const sampleText = "Press space key to start"

type Game struct {
	player       *player.Player
	purgerActor  *actor.Actor
	NPCActors    []*actor.Actor
	State        *gameplay.GameState
	keys         []ebiten.Key
	GameMode     int
	PlayMode     gameplay.PlayMode
	PromptPlayer bool
}

func NewGame() *Game {
	return &Game{
		State: &gameplay.GameState{Status: gameplay.StatusMap[gameplay.GameMenu]},
	}
}

func (g *Game) DrawActor(screen *ebiten.Image, actor *actor.Actor) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(actor.Position[0], actor.Position[1])
	rendering.DrawImageWithMatrix(screen, actor.Image, op)
}

func (g *Game) SpawnActors(screen *ebiten.Image) {
	for _, npc := range g.NPCActors {
		// if !npc.Draw {
		// 	continue
		// }
		g.DrawActor(screen, npc)
	}
	g.DrawActor(screen, g.purgerActor)
}

func (g *Game) PurgedCountStr() string {
	return "Purged: " + strconv.Itoa(g.State.PurgedCount)
}
func (g *Game) SparedCountStr() string {
	return "Spared: " + strconv.Itoa(g.State.SparedCount)
}

func (g *Game) InitHomeScreen(screen *ebiten.Image) {
	rendering.DrawCenteredText(screen, sampleText, ScreenWidth/2, ScreenHeight/2)
}

func (g *Game) InitKillFeed(screen *ebiten.Image) {
	rendering.DrawText(screen, g.PurgedCountStr(), ScreenWidth-100, 0)
	rendering.DrawText(screen, g.SparedCountStr(), ScreenWidth-100, 20)
}

func (g *Game) removeHiddenActors() {
	for key, npc := range g.NPCActors {
		if npc.Draw {
			continue
		}
		g.NPCActors = append(g.NPCActors[:key], g.NPCActors[key+1:]...)
	}
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.removeHiddenActors()

	// TODO: set default game mode and update based on user choice
	g.GameMode = 1
	g.PlayMode = gameplay.NewPlayMode(g.GameMode)

	// starts patrolling
	// TODO: maybe not set actor separately but get it from player??
	if g.player == nil {
		g.player = g.PlayMode.InitPlayer()
		g.purgerActor = g.player.Actor
	}

	if g.NPCActors == nil {
		g.NPCActors = g.PlayMode.InitNPCs()
	}

	switch g.State.Status {
	case StatusMap[GameMenu]:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State.Status = StatusMap[GameStarted]
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
		// set initial actors state
		g.PlayMode.InitActors(g.NPCActors)
		if len(g.keys) > 0 {
			g.PlayMode.HandleKeyboardInput(g.State, g.purgerActor, g.NPCActors, g.keys)
			g.purgerActor.SetLimitBounds(ScreenWidthFloat, ScreenHeightFloat)
		}
	case StatusMap[AwaitingUser]:
		g.PlayMode.HandlePlayerInput(g.State, g.NPCActors, g.State.Target)
	}
	return nil
}

func (g *Game) SetupCommonGameComponents(screen *ebiten.Image) {
	g.InitKillFeed(screen)
	g.SpawnActors(screen)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: cache images?
	switch g.State.Status {
	case StatusMap[GameMenu]:
		g.InitHomeScreen(screen)
	case StatusMap[GameStarted]:
		g.SetupCommonGameComponents(screen)
	case StatusMap[GamePaused]:
		g.SetupCommonGameComponents(screen)
		g.PlayMode.PauseGame(g.State, screen, ScreenWidthFloat, ScreenHeightFloat)
	case StatusMap[AwaitingUser]:
		g.SetupCommonGameComponents(screen)
		g.PlayMode.PropmptPlayer(g.State, g.purgerActor, screen)
	case StatusMap[GameEnded]:
		g.PlayMode.EndGame(g.State, screen)
	}

	g.PlayMode.CheckGameOverAndUpdateState(g.State, g.NPCActors)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return ScreenWidth, ScreenHeight
}
