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

type ModeFrostmourneHungers struct {
	BasePlayMode
}

// These functions need to be clalled each game tick
func (playmode *ModeFrostmourneHungers) Tick(gameState *GameState, gameActors []*actor.Actor, player *player.Player) {
	playmode.InitActors(gameActors)
	playmode.PurgeIfInAoE(gameState, gameActors, player)
	playmode.HandlePlayerInput(gameState, gameActors, player.Actor)
	playmode.CheckGameOverAndUpdateState(gameState, gameActors, player)
}

func (playmode *ModeFrostmourneHungers) EncounterNPCs(gameState *GameState, npc *actor.Actor) {
	// If the player is witin range/in collision with the NPC
	// Enable an action PURGE! button that the player can press

	// If the player's AoE ability is active, the player can purge all NPCs in range
}

func (playmode *ModeFrostmourneHungers) Purge(gameState *GameState, gameActors []*actor.Actor, npcActor *actor.Actor) {
	gameState.PurgedCount += 1
	playmode.RemoveNPC(gameState, gameActors, npcActor)
}

func (playmode *ModeFrostmourneHungers) PurgeIfInAoE(gameState *GameState, gameActors []*actor.Actor, player *player.Player) {
	for _, npcActor := range gameActors {
		if !npcActor.Draw || !npcActor.CollisionEnabled || len((player.Abilities)) == 0 {
			continue
		}
		if !npcActor.CollidesWithAbility(player.Abilities[0].Actor) {
			continue
		}
		// If the NPC is in the player's AoE ability, purge it
		playmode.Purge(gameState, gameActors, npcActor)
		if gameState.PurgedCount%4 == 0 {
			player.LevelUp()
		}
	}
}

func (playmode *ModeFrostmourneHungers) HandleKeyboardInput(
	gameState *GameState,
	player *player.Player,
	gameActors []*actor.Actor,
	keys []ebiten.Key) {
	player.HandleInput(keys)
	// What if the NPC goes over the player?

	playmode.PurgeIfInAoE(gameState, gameActors, player)

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		player.DeathAndDecay()
	}
}
func (playmode *ModeFrostmourneHungers) HandlePlayerInput(gameState *GameState, npcActors []*actor.Actor, npcActor *actor.Actor) {
	// Might not be needed as this mode does not have "Waiting" as game state
}

func (playmode *ModeFrostmourneHungers) CheckGameOverAndUpdateState(gameState *GameState, gameActors []*actor.Actor, player *player.Player) {
	if len(gameActors) == 0 {
		gameState.Status = StatusMap[GameWon]
	} else if player.Health <= 0 {
		gameState.Status = StatusMap[GameLost]
	}
}

func (playmode *ModeFrostmourneHungers) InitPlayer() *player.Player {
	// Initialize the player actor
	playerTexture := rendering.CreateTexture(utils.LoadFile("./assets/dk.png"))
	playerActor := actor.NewActor([2]float64{0, 0}, playerTexture, 14, "Purger", true)
	return player.NewPlayer(playerActor)
}

func (playmode *ModeFrostmourneHungers) InitNPCs() []*actor.Actor {
	// Initialize the NPC actors
	scourgeTexture := rendering.CreateTexture(utils.LoadFile("./assets/scv.png"))

	npcActor := actor.NewActor([2]float64{200, 200}, scourgeTexture, 4, "Scourge", true)
	npcActor1 := actor.NewActor([2]float64{400, 200}, scourgeTexture, 1, "Undead1", true)
	npcActor2 := actor.NewActor([2]float64{500, 300}, scourgeTexture, 1, "Undead2", true)
	npcActor3 := actor.NewActor([2]float64{250, 50}, scourgeTexture, 1, "Undead3", true)
	npcActor4 := actor.NewActor([2]float64{350, 50}, scourgeTexture, 1, "Undead3", true)
	npcActor5 := actor.NewActor([2]float64{450, 70}, scourgeTexture, 1, "Undead3", true)
	npcActor6 := actor.NewActor([2]float64{200, 600}, scourgeTexture, 1, "Undead3", true)

	return []*actor.Actor{npcActor, npcActor1, npcActor2, npcActor3, npcActor4, npcActor5, npcActor6}
}
