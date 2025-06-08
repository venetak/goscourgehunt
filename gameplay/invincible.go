package gameplay

import (
	_ "image/png"

	"github.com/actor"
	"github.com/player"

	"github.com/utils" // Replace with the correct path to the utils package

	"github.com/rendering" // Replace with the correct path to the rendering package

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ModeInvincible struct {
	BasePlayMode
}

// Constants defining the allowed game states
var (
	GameMenu     GameStatus = "Menu"
	GameStarted  GameStatus = "Started"
	GamePaused   GameStatus = "Paused"
	AwaitingUser GameStatus = "Waiting"
	GameEnded    GameStatus = "Ended"
	GameLost     GameStatus = "Lost"
	GameWon      GameStatus = "Won"
)

var StatusMap = map[GameStatus]int{
	GameMenu:     0,
	GameStarted:  1,
	GamePaused:   2,
	GameEnded:    3,
	AwaitingUser: 5,
	GameLost:     4, // GameLost is treated the same as GameEnded
	GameWon:      6, // GameWon is treated the same as GameEnded
}

// EncounterNPCs is called when the player collides with an NPC.
// It pauses the game and prompts the player to take action.
// The game state is set to AwaitingUser, and the player is prompted to either purge or spare the NPC.
func (playmode *ModeInvincible) EncounterNPCs(gameState *GameState, npc *actor.Actor) {
	// Pause game (stop all movement)
	gameState.Target = npc
	gameState.Status = StatusMap[AwaitingUser]
	gameState.PromptPlayer = true
	gameState.PromptPlayerText = "Press P to Purge or S to Spare"
}

// checkGameOverAndUpdateState checks if there are any remaining actors in the game.
// If there are no actors left, it sets the game state to GameEnded.
func (playmode *ModeInvincible) CheckGameOverAndUpdateState(gameState *GameState, gameActors []*actor.Actor, player *player.Player) {
	if len(gameActors) == 0 {
		gameState.Status = StatusMap[GameEnded]
	}
}

// Purge is called when the player chooses to purge an NPC.
// It increments the purged count and removes the NPC from the game.
// It also checks if the game is over.
func (playmode *ModeInvincible) Purge(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	gameState.PurgedCount += 1
	playmode.RemoveNPC(gameState, gameActors, npcActor)
}

// Spare is called when the player chooses to spare an NPC.
// It increments the spared count and removes the NPC from the game.
// It also checks if the game is over.
func (playmode *ModeInvincible) Spare(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	gameState.SparedCount += 1
	playmode.RemoveNPC(gameState, gameActors, npcActor)
}

func (playmode *ModeInvincible) InitPlayer() *player.Player {
	// Initialize the player actor
	playerTexture := rendering.CreateTexture(utils.LoadFile("./assets/arthas.png"))
	playerActor := actor.NewActor([2]float64{0, 0}, playerTexture, 14, "Purger", true)
	return player.NewPlayer(playerActor)
}

func (playmode *ModeInvincible) InitNPCs() []*actor.Actor {
	// Initialize the NPC actors
	scourgeTexture := rendering.CreateTexture(utils.LoadFile("./assets/pudge.png"))
	scourgeTexture1 := rendering.CreateTexture(utils.LoadFile("./assets/scourge.png"))

	npcActor := actor.NewActor([2]float64{200, 200}, scourgeTexture, 4, "Scourge", true)
	npcActor1 := actor.NewActor([2]float64{400, 200}, scourgeTexture1, 1, "Undead1", true)
	npcActor2 := actor.NewActor([2]float64{500, 300}, scourgeTexture1, 1, "Undead2", true)
	npcActor3 := actor.NewActor([2]float64{250, 50}, scourgeTexture1, 1, "Undead3", true)

	return []*actor.Actor{npcActor, npcActor1, npcActor2, npcActor3}
}

func (playmode *ModeInvincible) HandleKeyboardInput(
	gameState *GameState,
	player *player.Player,
	gameActors []*actor.Actor,
	keys []ebiten.Key) {
	player.HandleInput(keys)
	// What if the NPC goes over the player?
	for _, npcActor := range gameActors {
		if !npcActor.Draw || !npcActor.CollisionEnabled {
			continue
		}
		if player.Actor.CollidesWith(npcActor) {
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
