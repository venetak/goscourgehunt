package gameplay

import (
	_ "image/png"

	"github.com/actor"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/player"
	"github.com/rendering"
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
	Won              bool
	Lost             bool
}

type PlayMode interface {
	EncounterNPCs(gameState *GameState, npc *actor.Actor)
	InitActors(npcActors []*actor.Actor)
	PauseGame(gameState *GameState, screen *ebiten.Image, ScreenWidth, ScreenHeight float64)
	PropmptPlayer(gameState *GameState, player *actor.Actor, screen *ebiten.Image)
	Purge(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor)
	HandleKeyboardInput(gameState *GameState, player *player.Player, gameActors []*actor.Actor, keys []ebiten.Key)
	HandlePlayerInput(gameState *GameState, npcActors []*actor.Actor, npcActor *actor.Actor)
	EndGame(gameState *GameState, screen *ebiten.Image)
	CheckGameOverAndUpdateState(gameState *GameState, gameActors []*actor.Actor, player *player.Player)
	InitPlayer() *player.Player
	InitNPCs() []*actor.Actor
	RemoveNPC(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor)
}

type BasePlayMode struct{}

// EndGame is called when the game is over.
// It displays a message indicating that the game has ended and the player has completed the game.
func (b *BasePlayMode) EndGame(gameState *GameState, screen *ebiten.Image) {
	choiceText := "Congratulations! You have completed the game!"
	rendering.DrawBox(screen, float32(screen.Bounds().Dx()/2-200), float32(screen.Bounds().Dy()/2-50), 400, 100)
	rendering.DrawCenteredText(screen, choiceText, float64(screen.Bounds().Dx()/2), float64(screen.Bounds().Dy()/2))
	gameState.Status = StatusMap[GameEnded]
}

// It draws the player prompt at the actor's position on the screen.
func (playmode *BasePlayMode) PauseGame(gameState *GameState, screen *ebiten.Image, ScreenWidth, ScreenHeight float64) {
	choiceText := "Game Paused"
	rendering.DrawBox(screen, float32(ScreenWidth/2-100), float32(ScreenHeight/2-50), 200, 100)
	rendering.DrawCenteredText(screen, choiceText, ScreenWidth/2, ScreenHeight/2)
}

func (playmode *BasePlayMode) InitActors(npcActors []*actor.Actor) {
	for npcActor := range npcActors {
		if !npcActors[npcActor].Draw {
			continue
		}
		npcActors[npcActor].Patrol(10)
	}
}

// It draws the player prompt at the actor's position on the screen.
func (playmode *BasePlayMode) PropmptPlayer(gameState *GameState, player *actor.Actor, screen *ebiten.Image) {
	rendering.DrawPlayerPromptAtActorPos(screen, gameState.PromptPlayerText, player.Position)
}

// RemoveActor is called to remove an actor from the game.
// It sets the actor's Draw property to false, effectively removing it from the game.
// TODO: optimize by either removing the actor from the slice or using a pool of actors
func (playmde *BasePlayMode) RemoveActor(gameActors []*actor.Actor, npcActor *actor.Actor) {
	npcActor.Draw = false
}

// removeNPC is called to remove an NPC from the game.
// It removes the actor from the game, sets the PromptPlayer to false,
// and checks if the game is over.
func (playmode *BasePlayMode) RemoveNPC(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	playmode.RemoveActor(gameActors, npcActor)
	gameState.PromptPlayer = false
	gameState.Status = StatusMap[GameStarted]
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
