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
var (
	GameMenu     GameStatus = "Menu"
	GameStarted  GameStatus = "Started"
	GamePaused   GameStatus = "Paused"
	GameEnded    GameStatus = "Ended"
	AwaitingUser GameStatus = "Waiting"
)

var StatusMap = map[GameStatus]int{
	GameMenu:     0,
	GameStarted:  1,
	GamePaused:   2,
	GameEnded:    3,
	AwaitingUser: 4,
}

var currentGameStatus GameStatus = GameMenu

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

type GameState struct {
	Status           int
	PromptPlayer     bool
	PromptPlayerText string
	PurgedCount      int
	SparedCount      int
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
		State:        &GameState{Status: StatusMap[GameMenu]},
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
			g.PlayMode = newPlayMode(g.GameMode)
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
// -----------------------------------------------------------
type PlayMode interface {
	EncounterNPCs(gameState *GameState, npc *actor.Actor)
	InitActors(*actor.Actor, []*actor.Actor, []ebiten.Key)
	PauseGame(gameState *GameState, player *actor.Actor, screen *ebiten.Image)
	Purge(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor)
	Spare(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor)
	HandleKeyboardInput(gameState *GameState, player *actor.Actor, gameActors []*actor.Actor, keys []ebiten.Key)
	HandlePlayerInput(gameState *GameState, npcActors []*actor.Actor, npcActor *actor.Actor)
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
func (playmode *ModeInvincible) EncounterNPCs(gameState *GameState, npc *actor.Actor) {
	// Pause game (stop all movement)
	gameState.Target = npc
	gameState.Status = StatusMap[AwaitingUser]
	gameState.PromptPlayer = true
	gameState.PromptPlayerText = "Press P to Purge or S to Spare"
}

func (playmode *ModeInvincible) UpdateScore() {
	// Update score based on purged and saved NPCs
	// g.purgedCount += 1
	// g.savedCount += 1
}

func (playmode *ModeInvincible) PauseGame(gameState *GameState, player *actor.Actor, screen *ebiten.Image) {
	rendering.DrawPlayerPromptAtActorPos(screen, gameState.PromptPlayerText, player.Position, arcadeFaceSource, 10)
}

func (playmde *ModeInvincible) RemoveActor(gameActors []*actor.Actor, npcActor *actor.Actor) {
	npcActor.Draw = false
}

func (playmode *ModeInvincible) checkGameOverAndUpdateState(gameState *GameState, gameActors []*actor.Actor) {
	if len(gameActors) == 0 {
		gameState.Status = StatusMap[GameEnded]
	} else {
		gameState.Status = StatusMap[GameStarted]
	}
}

func (playmode *ModeInvincible) removeNPC(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	playmode.RemoveActor(gameActors, npcActor)
	gameState.PromptPlayer = false

	// verbose: true...
	playmode.checkGameOverAndUpdateState(gameState, gameActors)
}

func (playmode *ModeInvincible) Purge(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	gameState.PurgedCount += 1
	playmode.removeNPC(gameState, gameActors, npcActor)
}

func (playmode *ModeInvincible) Spare(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	gameState.SparedCount += 1
	playmode.removeNPC(gameState, gameActors, npcActor)
}

// This is a little strage, should it be part of the game package?
func (playmode *ModeInvincible) InitActors(player *actor.Actor, npcActors []*actor.Actor, pressedKeys []ebiten.Key) {
	for npcActor := range npcActors {
		if !npcActors[npcActor].Draw {
			continue
		}
		npcActors[npcActor].Patrol(10)
	}
}

func (playmode *ModeInvincible) HandleKeyboardInput(gameState *GameState, player *actor.Actor, gameActors []*actor.Actor, keys []ebiten.Key) {
	player.HandleInput(keys)
	player.SetLimitBounds(ScreenWidthFloat, ScreenHeightFloat)
	// What if the NPC goes over the player?
	for _, npcActor := range gameActors {
		if !npcActor.Draw {
			continue
		}
		if player.CollidesWith(npcActor) {
			playmode.EncounterNPCs(gameState, npcActor)
		}
	}
}

// TODO: might be better to move this to the main input handler
func (playmode *ModeInvincible) HandlePlayerInput(gameState *GameState, npcActors []*actor.Actor, npcActor *actor.Actor) {

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		playmode.Purge(gameState, npcActors, npcActor)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		playmode.Spare(gameState, npcActors, npcActor)
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
