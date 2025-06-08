package player

import (
	"time"

	"github.com/actor"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/rendering"
	"github.com/utils" // Import the utils package
)

type Abilities struct {
	DeathAndDecay *actor.Actor // The Area of Effect (AoE) actor for the player
	DeathCoil     *actor.Actor // Projectile Death Coil damage ability
	BurstOfLight  *actor.Actor // Healing single-target ability
}

type AbilityType string

const (
	DeathAndDecayType AbilityType = "AOEDuration"
	DeathCoilType     AbilityType = "ProjectileDamage"
	BurstOfLightType  AbilityType = "SingleTargetHeal"
)

type Ability struct {
	Actor     *actor.Actor
	Duration  float64     // Duration in seconds for the Death and Decay ability
	Type      AbilityType // Type of the ability, e.g., "AoE", "Damage", "Heal"
	StartTime time.Time   // Time when the ability was activated
}

// Player represents the player character in the game.
type Player struct {
	Actor     *actor.Actor // The actor representing the player
	Abilities []*Ability   // The Area of Effect (AoE) actor for the player
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

// HandleInput processes the input from the user and updates the actor's movement direction accordingly.
// It takes a slice of ebiten.Key values representing the keys currently being pressed.
// The function resets the actor's movement direction, checks which keys are pressed, and updates
// the movement direction (MoveDirectionX and MoveDirectionY) based on the arrow keys.
// Finally, it calculates the new position and moves the actor in the specified direction.
func (player *Player) HandleInput(inputDeviceActionButtonNames []ebiten.Key) {
	actor := player.Actor
	actor.ResetMoveDirection()

	for _, key := range inputDeviceActionButtonNames {
		if ebiten.IsKeyPressed(key) {
			switch key {
			case ebiten.KeyArrowLeft:
				actor.MoveDirectionX = -1
			case ebiten.KeyArrowRight:
				actor.MoveDirectionX = 1
			case ebiten.KeyArrowUp:
				actor.MoveDirectionY = -1
			case ebiten.KeyArrowDown:
				actor.MoveDirectionY = 1
			}
		}
	}

	newPosition := [2]float64{actor.MoveDirectionX, actor.MoveDirectionY}
	actor.MoveIn(newPosition)
}

func (p *Player) DeathAndDecay() {
	playerBonds := p.Actor.GetBoundingRect()
	playerCenterX := playerBonds.PositionX + playerBonds.Width/2
	playerCenterY := playerBonds.PositionY + playerBonds.Height/2

	aoeTexture := rendering.CreateTexture(utils.LoadFile("./assets/circle1.png"))

	aoeActor := actor.NewActor([2]float64{0, 0}, aoeTexture, 4, "Death and Decay", false)
	aoeBonds := aoeActor.GetBoundingRect()
	aoeActor.Position = [2]float64{
		playerCenterX - aoeBonds.Width/2,
		playerCenterY - aoeBonds.Height/2,
	}

	DeathAndDecay := &Ability{
		Actor:     aoeActor,
		Duration:  3,                 // Duration in seconds for the Death and Decay ability
		Type:      DeathAndDecayType, // Type of the ability
		StartTime: time.Now(),        // Time when the ability was activated
	}

	p.Abilities = append(p.Abilities, DeathAndDecay)
}

// UpdateAbilitiesDurations updates the durations of the player's abilities.
// It removes any abilities that have expired based on their start time and duration.
func (p *Player) UpdateAbilitiesDurations() {
	abilitiesCopy := make([]*Ability, 0, len(p.Abilities)) // Pre-allocate for efficiency

	for _, ability := range p.Abilities {
		if ability.StartTime.IsZero() {
			continue
		}
		timeElapsed := time.Since(ability.StartTime).Seconds()
		if timeElapsed >= ability.Duration {
			continue // Skip abilities that have expired
		}
		abilitiesCopy = append(abilitiesCopy, ability)
	}
	p.Abilities = abilitiesCopy
}

// LevelUp increases the player's level and resets health and mana.
func (p *Player) LevelUp() {
	p.Level++
	p.Health = 100 // Reset health to default value
	p.Mana = 50    // Reset mana to default value
}
