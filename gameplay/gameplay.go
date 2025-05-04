package gameplay

import (
	"github/actor"

	"github.com/hajimehoshi/ebiten/v2"
)

type Gameplay interface {
	EncounterNPCs()
	SpawnActors(*actor.Actor, []*actor.Actor, []ebiten.Key)
}

type ModeInvincible struct {
	purgedCount int
	savedCount  int
}

type ModeFrostmourneHungers struct {
	purgedCount int
}

func GetGameplay(gameMode int) Gameplay {
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
func (g *ModeInvincible) EncounterNPCs() {
	// Pause game (stop all movement)
	// freeze drawing and ask for user input
	// show popup asking if the player wants to spare or purge the NPC
}

func (g *ModeInvincible) UpdateScore() {
	// Update score based on purged and saved NPCs
	// g.purgedCount += 1
	// g.savedCount += 1
}

// This is a little strage, should it be part of the game package?
func (g *ModeInvincible) SpawnActors(player *actor.Actor, npcActors []*actor.Actor, pressedKeys []ebiten.Key) {
	for npcActor := range npcActors {
		npcActors[npcActor].Patrol(10)
	}
	if len(pressedKeys) > 0 {
		player.HandleInput(pressedKeys)
		if player.CollidesWith(npcActors[0]) {
			player.RollbackPosition()
		}
	}
}

// -----------------------------------------------------------
func (g *ModeFrostmourneHungers) EncounterNPCs() {
	// If the player is witin range/in collision with the NPC
	// Enable an action PURGE! button that the player can press
}

func (g *ModeFrostmourneHungers) UpdateScore() {
	// Update score based on purged NPCs
	// g.purgedCount += 1
}

func (g *ModeFrostmourneHungers) SpawnActors() {}

// manage state:
// - if the actor is moving - check for collision
// - if the actor is not moving - check for collision?
// - if the actor is is collision:
// - pause the game
// stop all movement; freeze drawing and ask for user input
// - purge the npc - yes - delete the npc
// - remove it from game objects
// increase purged count??
// - spare it?
