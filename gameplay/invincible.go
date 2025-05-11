package gameplay

import (
	_ "image/png"

	"github/actor"

	"github.com/rendering" // Replace with the correct path to the rendering package

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ModeInvincible struct {
}

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

func (playmode *ModeInvincible) UpdateScore() {
	// Update score based on purged and saved NPCs
	// g.purgedCount += 1
	// g.savedCount += 1
}

// It draws the player prompt at the actor's position on the screen.
func (playmode *ModeInvincible) PauseGame(gameState *GameState, screen *ebiten.Image, ScreenWidth, ScreenHeight float64) {
	choiceText := "Game Paused"
	rendering.DrawBox(screen, float32(ScreenWidth/2-100), float32(ScreenHeight/2-50), 200, 100)
	rendering.DrawCenteredText(screen, choiceText, ScreenWidth/2, ScreenHeight/2)
}

// It draws the player prompt at the actor's position on the screen.
func (playmode *ModeInvincible) PropmptPlayer(gameState *GameState, player *actor.Actor, screen *ebiten.Image) {
	rendering.DrawPlayerPromptAtActorPos(screen, gameState.PromptPlayerText, player.Position)
}

// RemoveActor is called to remove an actor from the game.
// It sets the actor's Draw property to false, effectively removing it from the game.
// TODO: optimize by either removing the actor from the slice or using a pool of actors
func (playmde *ModeInvincible) RemoveActor(gameActors []*actor.Actor, npcActor *actor.Actor) {
	npcActor.Draw = false
}

// checkGameOverAndUpdateState checks if there are any remaining actors in the game.
// If there are no actors left, it sets the game state to GameEnded.
func (playmode *ModeInvincible) CheckGameOverAndUpdateState(gameState *GameState, gameActors []*actor.Actor) {
	if len(gameActors) == 0 {
		gameState.Status = StatusMap[GameEnded]
	}
}

// removeNPC is called to remove an NPC from the game.
// It removes the actor from the game, sets the PromptPlayer to false,
// and checks if the game is over.
func (playmode *ModeInvincible) removeNPC(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	playmode.RemoveActor(gameActors, npcActor)
	gameState.PromptPlayer = false
	gameState.Status = StatusMap[GameStarted]
}

// Purge is called when the player chooses to purge an NPC.
// It increments the purged count and removes the NPC from the game.
// It also checks if the game is over.
func (playmode *ModeInvincible) Purge(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	gameState.PurgedCount += 1
	playmode.removeNPC(gameState, gameActors, npcActor)
}

// Spare is called when the player chooses to spare an NPC.
// It increments the spared count and removes the NPC from the game.
// It also checks if the game is over.
func (playmode *ModeInvincible) Spare(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	gameState.SparedCount += 1
	playmode.removeNPC(gameState, gameActors, npcActor)
}

// InitActors initializes the actors in the game.
// It sets the patrol speed for each NPC actor and starts their patrol behavior.
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

// EndGame is called when the game is over.
// It displays a message indicating that the game has ended and the player has completed the game.
func (playmode *ModeInvincible) EndGame(gameState *GameState, screen *ebiten.Image) {
	choiceText := "Congratulations! You have completed the game!"
	rendering.DrawBox(screen, float32(screen.Bounds().Dx()/2-200), float32(screen.Bounds().Dy()/2-50), 400, 100)
	rendering.DrawCenteredText(screen, choiceText, float64(screen.Bounds().Dx()/2), float64(screen.Bounds().Dy()/2))
	gameState.Status = StatusMap[GameEnded]
}
