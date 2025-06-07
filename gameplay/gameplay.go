package gameplay

import (
	_ "image/png"

	"github.com/actor"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/player"
)

type Hud struct {
	PlayerControls  []*ebiten.Image // help overlay for player controls
	PlayerAbilities []*ebiten.Image // buttons for player abilities (with bindings)
	PlayerStats     []*ebiten.Image // player stats (health, mana, etc.)
	PlayerNameplate *ebiten.Image   // player nameplate
	GameMenu        []*ebiten.Image // game menu images (start, pause, etc.)
	NPCNameplate    *ebiten.Image   // NPC nameplate
}

type GameStatus string

type GameState struct {
	Status           int
	PromptPlayer     bool
	PromptPlayerText string
	PurgedCount      int
	SparedCount      int
	Target           *actor.Actor
	TimeElapsed      float64
}

type PlayMode interface {
	EncounterNPCs(gameState *GameState, npc *actor.Actor)
	InitActors(npcActors []*actor.Actor)
	PauseGame(gameState *GameState, screen *ebiten.Image, ScreenWidth, ScreenHeight float64)
	PropmptPlayer(gameState *GameState, player *actor.Actor, screen *ebiten.Image)
	Purge(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor)
	Spare(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor)
	HandleKeyboardInput(gameState *GameState, player *player.Player, gameActors []*actor.Actor, keys []ebiten.Key)
	HandlePlayerInput(gameState *GameState, npcActors []*actor.Actor, npcActor *actor.Actor)
	EndGame(gameState *GameState, screen *ebiten.Image)
	CheckGameOverAndUpdateState(gameState *GameState, gameActors []*actor.Actor)
	InitPlayer() *player.Player
	InitNPCs() []*actor.Actor
}

// Game mode factory
// This function creates a new PlayMode instance based on the provided gameMode parameter.
func NewPlayMode(gameMode int) PlayMode {
	switch gameMode {
	case 1:
		return &ModeInvincible{}
	case 2:
		return &ModeFrostmourneHungers{}
	default:
		return nil
	}
}
