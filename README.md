# A Game Inspired by WoW
![](https://github.com/venetak/goscourgehunt/blob/main/demo.PNG)

# Project Design

This is a 2D game that uses the [Ebitengine](https://ebitengine.org). The engine handles the game loop - the update and draw callbacks that are responsible for the game state and rendering.
- Update is called by default 60 times per second.
- Draw is synced with the refresh rate of the monitor.

## Components

### Actor
Every object that's part of the game is an actor - player, environment (trees, houses), characters, etc.  
There's no dedicated physics engine, so the Actor handles its locomotion, input processing, collision detection, and state.  
It supports AABB collision detection.  
It supports movement across x, y, and the diagonals.

### Rendering
Responsible for handling the drawing of actors on the scene. It utilizes the drawing API of Ebitengine to provide reusable rendering functionality.

### Game
The Ebitengine Game object. Implements the Update, Draw, and Layout functions. Handles keyboard input. Manages the game state. Provides an abstraction interface that allows painless switching between game modes.

For example:
```go
switch game.State.Status {
    case StatusMap[GamePaused]:
        // in easy gameplay mode will pause
        // in hardcore mode will pause, but any DoT effects will be active... 
        game.PlayMode.PauseGame()
    case StatusMap[GameStarted]:
        // in easy mode will spawn patrolling NPCs
        // in hardcore mode will spawn NPCs, traps, weather AoEs, etc.
        game.PlayMode.InitActors()
    case StatusMap[UserInput]:
        // in easy mode will support movement, menu controls
        // in hardcore mode will support bindings for multiple player abilities
        game.PlayMode.HandleKeyboardInput()
    }
```

### Gameplay
The implementation of the game logic. All instances implement the Gameplay interface. The different game modes determine:
- what happens when the player encounters NPCs.
- the different logic for removing an actor from the game.
- the conditions for game state change - for example, what requirements should be met in order to end the game in easy and hard modes?

# Game Design

## Main Menu
You, Prince Arthas Menethil, arrive at the gates of Stratholme. But you are too late. The batch of infected grain was already delivered and distributed across the city. You must act now, or soon the whole of Lordaeron will have mindless Scourge crawling around, killing innocent people.

### Select Game Mode
- [1] The Boy Who Killed Invincible
    - You'll be given a choice each time you encounter