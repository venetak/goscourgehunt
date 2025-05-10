package gameplay

import (
	_ "image/png"

	"github/actor"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

type GameState struct {
	Status           int
	PromptPlayer     bool
	PromptPlayerText string
	PurgedCount      int
	SparedCount      int
	Target           *actor.Actor
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

func NewPlayMode(gameMode int) PlayMode {
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
