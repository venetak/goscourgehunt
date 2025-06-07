package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Hud struct {
	GameMenu        []*ebiten.Image // game menu images (start, pause, etc.)
	PlayerControls  []*ebiten.Image // help overlay for player controls
	PlayerAbilities []*ebiten.Image // buttons for player abilities (with bindings)
	PlayerStats     []*ebiten.Image // player stats (health, mana, etc.)
	PlayerNameplate *ebiten.Image   // player nameplate
	NPCNameplate    *ebiten.Image   // NPC nameplate
}
