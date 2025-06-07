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
	Debug        bool
	player       *player.Player
	purgerActor  *actor.Actor
	NPCActors    []*actor.Actor
	State        *gameplay.GameState
	keys         []ebiten.Key
	GameMode     int
	PlayMode     gameplay.PlayMode
	PromptPlayer bool
}

func NewGame(debug bool) *Game {
	return &Game{
		Debug: debug,
		State: &gameplay.GameState{Status: gameplay.StatusMap[gameplay.GameMenu]},
	}
}

func (g *Game) DrawActor(screen *ebiten.Image, actor *actor.Actor) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(actor.Position[0], actor.Position[1])
	rendering.DrawImageWithMatrix(screen, actor.Image, op, g.Debug)
}

func (g *Game) SpawnActors(screen *ebiten.Image, actors []*actor.Actor) {
	for _, actor := range actors {
		if !actor.Draw {
			continue
		}
		g.DrawActor(screen, actor)
	}
}

func (g *Game) DrawPlayer(screen *ebiten.Image) {
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
	NPCActorsCopy := make([]*actor.Actor, 0, len(g.NPCActors)) // Pre-allocate for efficiency
	for _, npc := range g.NPCActors {
		if npc.Draw {
			NPCActorsCopy = append(NPCActorsCopy, npc)
		}
	}
	g.NPCActors = NPCActorsCopy
}

func (g *Game) SetupCommonGameComponents(screen *ebiten.Image) {
	g.InitKillFeed(screen)
	g.DrawPlayer(screen)
	g.SpawnActors(screen, g.NPCActors)
	g.SpawnPlayerAbilities(screen)
}

func (g *Game) SpawnPlayerAbilities(screen *ebiten.Image) {
	for _, ability := range g.player.Abilities {
		g.DrawActor(screen, ability.Actor)
	}
}

// Game lifecycle methods
func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	// TODO: set default game mode and update based on user choice
	g.GameMode = 1
	g.PlayMode = gameplay.NewPlayMode(g.GameMode)

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
			g.PlayMode.HandleKeyboardInput(g.State, g.player, g.NPCActors, g.keys)
			g.purgerActor.SetLimitBounds(ScreenWidthFloat, ScreenHeightFloat)
		}
		g.removeHiddenActors()
		g.player.UpdateAbilitiesDurations()
	case StatusMap[AwaitingUser]:
		g.PlayMode.HandlePlayerInput(g.State, g.NPCActors, g.State.Target)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// TODO: cache images?
	switch g.State.Status {
	case StatusMap[GameMenu]:
		g.InitHomeScreen(screen)
	case StatusMap[GameStarted]:
		g.SetupCommonGameComponents(screen)
		g.PlayMode.CheckGameOverAndUpdateState(g.State, g.NPCActors)
	case StatusMap[GamePaused]:
		g.SetupCommonGameComponents(screen)
		g.PlayMode.PauseGame(g.State, screen, ScreenWidthFloat, ScreenHeightFloat)
	case StatusMap[AwaitingUser]:
		g.SetupCommonGameComponents(screen)
		g.PlayMode.PropmptPlayer(g.State, g.purgerActor, screen)
	case StatusMap[GameEnded]:
		g.PlayMode.EndGame(g.State, screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return ScreenWidth, ScreenHeight
}
