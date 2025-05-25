package player

import (
	"github.com/actor"
)

type Abilities struct {
	DeathAndDecay *actor.Actor // The Area of Effect (AoE) actor for the player
	DeathCoil     *actor.Actor // Projectile Death Coil damage ability
	BurstOfLight  *actor.Actor // Healing single-target ability
}

// Player represents the player character in the game.
type Player struct {
	Actor     *actor.Actor // The actor representing the player
	Abilities *Abilities   // The Area of Effect (AoE) actor for the player
	Health    int          // Player's health
	Mana      int          // Player's mana
	Level     int          // Player's level
	Target    *actor.Actor // The current target of the player
}

// NewPlayer creates a new Player instance with the given actor and AoE actor.
func NewPlayer(actor *actor.Actor) *Player {
	return &Player{
		Actor:  actor,
		Health: 100, // Default health
		Mana:   50,  // Default mana
		Level:  1,   // Starting level
	}
}

func (p *Player) DeathAndDecay() {

}
