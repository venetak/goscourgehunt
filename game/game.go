package game

import (
	"bytes"
	_ "image/png"
	"log"
	"strconv"

	"github/actor"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rendering"
)

var (
	arcadeFaceSource *text.GoTextFaceSource
)

type GameStatus string

// Constants defining the allowed game states
const (
	GameMenu     GameStatus = "Menu"
	GameStarted  GameStatus = "Started"
	GamePaused   GameStatus = "Paused"
	GameEnded    GameStatus = "Ended"
	AwaitingUser GameStatus = "Waiting"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 400
	LineSpacing  = 0
	FontSize     = 10
)

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

type GameState struct {
	status           GameStatus
	PromptPlayer     bool
	PromptPlayerText string
	PurgedCount      int
	SavedCount       int
	Target           *actor.Actor
}

type Game struct {
	scourgeActor *actor.Actor
	purgerActor  *actor.Actor
	NPCActors    []*actor.Actor
	State        *GameState
	keys         []ebiten.Key
	GameMode     int
	PlayMode     PlayMode
	PromptPlayer bool
}

func NewGame(playerActor *actor.Actor, npcActors []*actor.Actor) *Game {
	return &Game{
		State:        &GameState{status: GameMenu},
		scourgeActor: playerActor,
		purgerActor:  playerActor,
		NPCActors:    npcActors,
	}
}

func (g *Game) SpawnActors(screen *ebiten.Image) {
	for _, npc := range g.NPCActors {
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

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	switch g.State.status {
	case GameMenu:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.State.status = GameStarted
			g.GameMode = 1
			// TODO: rename g to game
			// AND Gameplay to PlayMode
			// and you'll have game.PlayMode
			g.PlayMode = newPlayMode(g.GameMode)
		}
	case GamePaused:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State.status = GameStarted
			return nil
		}
	case GameStarted:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State.status = GamePaused
			return nil
		}

		// starts patrolling
		g.PlayMode.InitActors(g.purgerActor, g.NPCActors, g.keys)
		if len(g.keys) > 0 {
			g.purgerActor.HandleInput(g.keys)
			if g.purgerActor.CollidesWith(g.NPCActors[0]) {
				g.PlayMode.EncounterNPCs(g, g.NPCActors[0])
			}
		}
	case AwaitingUser:
		if ebiten.IsKeyPressed(ebiten.KeyP) {
			g.PlayMode.Purge(g, g.State.Target)
		}
		// if ebiten.IsKeyPressed(ebiten.KeyS) {
		// 	g.PlayMode.Purge(g, g.State.Target)
		// 	g.State.status = GameStarted
		// 	g.State.PromptPlayer = false
		// }
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Image at: (0, 0)"))
	// TODO: cache images?
	switch g.State.status {
	case GameMenu:
		rendering.DrawCenteredText(screen, sampleText, ScreenWidth/2, ScreenHeight/2, arcadeFaceSource, FontSize)
	case GameStarted:
		g.SpawnActors(screen)
		rendering.DrawText(screen, g.PurgedCountStr(), ScreenWidth-100, 0, arcadeFaceSource, FontSize)
	case GamePaused:
		g.SpawnActors(screen)
		choiceText := "Game Paused"
		rendering.DrawBox(screen, ScreenWidth/2-100, ScreenHeight/2-50, 200, 100)
		rendering.DrawCenteredText(screen, choiceText, ScreenWidth/2, ScreenHeight/2, arcadeFaceSource, FontSize)
	case AwaitingUser:
		g.SpawnActors(screen)
		g.PlayMode.PauseGame(g, screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (ScreenWidth, screenHeight int) {
	return 800, 400
}

// PlayMode interface
// -----------------------------------------------------------
type PlayMode interface {
	EncounterNPCs(game *Game, npc *actor.Actor)
	InitActors(*actor.Actor, []*actor.Actor, []ebiten.Key)
	PauseGame(g *Game, screen *ebiten.Image)
	Purge(g *Game, npcActor *actor.Actor)
}

type ModeInvincible struct {
}

type ModeFrostmourneHungers struct {
}

func newPlayMode(gameMode int) PlayMode {
	switch gameMode {
	case 1:
		return &ModeInvincible{}
	case 2:
		// return &ModeFrostmourneHungers{}
	default:
		return nil
	}

	return &ModeInvincible{}
}

// InvincibleGameplay implements the Gameplay interface
// -----------------------------------------------------------
func (playmode *ModeInvincible) EncounterNPCs(game *Game, npc *actor.Actor) {
	// Pause game (stop all movement)
	state := game.State
	game.State.Target = npc
	state.status = AwaitingUser
	state.PromptPlayer = true
	state.PromptPlayerText = "Press P to Purge or S to Spare"
}

func (playmode *ModeInvincible) UpdateScore() {
	// Update score based on purged and saved NPCs
	// g.purgedCount += 1
	// g.savedCount += 1
}

func (playmode *ModeInvincible) PauseGame(g *Game, screen *ebiten.Image) {
	rendering.DrawPlayerPromptAtActorPos(screen, g.State.PromptPlayerText, g.purgerActor.Position, arcadeFaceSource, 10)
}

func (playmode *ModeInvincible) Purge(g *Game, npcActor *actor.Actor) {
	for i, npc := range g.NPCActors {
		if npc.Id == npcActor.Id {
			g.NPCActors = append(g.NPCActors[:i], g.NPCActors[i+1:]...)
			break
		}
	}
	g.State.PurgedCount += 1
	g.State.status = GameStarted
	g.State.PromptPlayer = false
}

// This is a little strage, should it be part of the game package?
func (playmode *ModeInvincible) InitActors(player *actor.Actor, npcActors []*actor.Actor, pressedKeys []ebiten.Key) {
	for npcActor := range npcActors {
		npcActors[npcActor].Patrol(10)
	}
}

// -----------------------------------------------------------
func (playmode *ModeFrostmourneHungers) EncounterNPCs() {
	// If the player is witin range/in collision with the NPC
	// Enable an action PURGE! button that the player can press
}

func (playmode *ModeFrostmourneHungers) UpdateScore() {
	// Update score based on purged NPCs
	// g.purgedCount += 1
}

func (playmode *ModeFrostmourneHungers) InitActors() {}
