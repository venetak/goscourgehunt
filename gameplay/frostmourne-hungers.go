package gameplay

import (
	"github/actor"

	"github.com/hajimehoshi/ebiten/v2"
)

type ModeFrostmourneHungers struct {
}

func (playmode *ModeFrostmourneHungers) EncounterNPCs(gameState *GameState, npc *actor.Actor) {
	// If the player is witin range/in collision with the NPC
	// Enable an action PURGE! button that the player can press

	// If the player's AoE ability is active, the player can purge all NPCs in range
}

func (playmode *ModeFrostmourneHungers) InitActors(*actor.Actor, []*actor.Actor, []ebiten.Key) {
	// Initialize actors for the Frostmourne Hungers mode
	// Include people, undead and abominations, animals
	// More NPCs, more undead and higher level abbominations
}

func (playmode *ModeFrostmourneHungers) PropmptPlayer(gameState *GameState, player *actor.Actor, screen *ebiten.Image) {
	// Draw the player prompt at the actor's position on the screen
}

func (playmode *ModeFrostmourneHungers) PauseGame(gameState *GameState, screen *ebiten.Image, ScreenWidth, ScreenHeight float64) {
	// Same as the other, but the text might be different
}

func (playmode *ModeFrostmourneHungers) Purge(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	// Implement mass purge as well.
	// Some of the NPCs could come back and undead
}

func (playmode *ModeFrostmourneHungers) Spare(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	// Painless death :D
}
func (playmode *ModeFrostmourneHungers) HandleKeyboardInput(gameState *GameState, player *actor.Actor, gameActors []*actor.Actor, keys []ebiten.Key) {
	// Will have handlers for mass purge
	// and maybe some additionl abilities, like Death Coil
}
func (playmode *ModeFrostmourneHungers) HandlePlayerInput(gameState *GameState, npcActors []*actor.Actor, npcActor *actor.Actor) {
	// Might not be needed as this mode does not have "Waiting" as game state
}

func (playmode *ModeFrostmourneHungers) UpdateScore() {
	// Update score based on purged NPCs
	// g.purgedCount += 1
}
