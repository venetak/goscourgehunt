package gameplay

import (
	"bytes"
	_ "image/png"
	"log"
	"strconv"

	"github/actor"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/rendering"
)

type Game interface  {
	scourgeActor *actor.Actor
	purgerActor  *actor.Actor
	NPCActors    []*actor.Actor
	State        *GameState
	keys         []ebiten.Key
	GameMode     int
	PlayMode     PlayMode
	PromptPlayer bool
}