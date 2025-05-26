package actor

import (
	"log"
	"math"

	"github.com/google/uuid"

	"github.com/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type Actor struct {
	Id               string
	Name             string
	Position         [2]float64
	initialPosition  [2]float64
	targetPosition   [2]float64
	Image            *ebiten.Image
	Speed            float64
	MoveDirectionX   float64
	MoveDirectionY   float64
	moveRange        float64
	Draw             bool
	CollisionEnabled bool // Indicates if the actor can collide with other actors
}

type BoundingRect struct {
	PositionX float64
	PositionY float64
	Width     float64
	Height    float64
}

func NewActor(position [2]float64, image *ebiten.Image, speed float64, name string, collision bool) *Actor {
	return &Actor{
		Id:               uuid.New().String(),
		Name:             name,
		Position:         position,
		initialPosition:  position,
		targetPosition:   position,
		Image:            image,
		Speed:            speed, // Default speed
		MoveDirectionX:   0.0,   // Default direction
		MoveDirectionY:   0.0,
		moveRange:        100.0, // Default move range
		Draw:             true,
		CollisionEnabled: collision, // Default collision behavior
	}
}

func (actor *Actor) ResetMoveDirection() {
	actor.MoveDirectionX = 0
	actor.MoveDirectionY = 0
}

// RollbackPosition moves the actor back to its previous position.
// This is useful when the actor collides with another actor or an obstacle.
func (actor *Actor) RollbackPosition() {
	actor.MoveIn([2]float64{-actor.MoveDirectionX, -actor.MoveDirectionY})
}

// MoveIn moves the actor in the specified direction.
// The direction is represented as a 2D vector (dx, dy).
// The function calculates the distance to move based on the speed of the actor.
// It normalizes the direction vector to ensure smooth movement, even when moving diagonally.
func (actor *Actor) MoveIn(direction [2]float64) {
	dx := direction[0]
	dy := direction[1]

	// Pythagorean theorem for the hypotenuse od a right triangle
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	if distance > 0 {
		// normalize the direction vector to get smooth movement on the diagonals
		dx = (dx / distance) * actor.Speed
		dy = (dy / distance) * actor.Speed

		actor.Position[0] += dx
		actor.Position[1] += dy
	}
}

// GetBoundingRect calculates the bounding rectangle of the actor.
// The bounding rectangle is used for collision detection.
// It takes into account the actor's position and the dimensions of the image.
func (actor *Actor) GetBoundingRect() *BoundingRect {
	playerRect := actor.Image.Bounds()

	return &BoundingRect{
		PositionX: actor.Position[0],
		PositionY: actor.Position[1],
		Width:     float64(playerRect.Dx()),
		Height:    float64(playerRect.Dy()),
	}
}

// AABB (Axis-Aligned Bounding Box) Collision detection
// Checks if two actors are colliding based on their bounding rectangles.
// - npc (non-player character) The second actor
// returns true if they are colliding, and false otherwise.
func (actor *Actor) CollidesWith(npc *Actor) bool {
	playerRect := actor.GetBoundingRect()
	npcRect := npc.GetBoundingRect()

	return rectCollition(playerRect, npcRect)
}

// TODO: maybe move to phisics package?
func rectCollition(rect1, rect2 *BoundingRect) bool {
	if rect1.PositionX < rect2.PositionX+rect2.Width &&
		rect1.PositionX+rect1.Width > rect2.PositionX &&
		rect1.PositionY < rect2.PositionY+rect2.Height &&
		rect1.PositionY+rect1.Height > rect2.PositionY {
		log.Print("Collision!!!!!--------")
		return true
	}

	return false
}

func (actor *Actor) CollidesWithAbility(ability *Actor) bool {
	abilityBounds := ability.GetBoundingRect()
	abilityRadius := abilityBounds.Width / 2 // Assuming the ability is circular, use half of its width as radius
	abilityCenterX := ability.Position[0] + abilityRadius
	abilityCenterY := ability.Position[1] + abilityRadius
	return rectWithCircleCollision(actor.GetBoundingRect(), abilityCenterX, abilityCenterY, abilityRadius)
}

// detects if a rectangle collides with a circle.
// It calculates the closest point on the rectangle to the circle's center
// and checks if the distance from that point to the circle's center is less than or equal to the circle's radius.
func rectWithCircleCollision(rect *BoundingRect, circleX, circleY, circleRadius float64) bool {
	// Find the closest point on the rectangle to the circle
	closestX := math.Max(rect.PositionX, math.Min(circleX, rect.PositionX+rect.Width))
	closestY := math.Max(rect.PositionY, math.Min(circleY, rect.PositionY+rect.Height))

	// Calculate the distance from the closest point to the circle's center
	dx := closestX - circleX
	dy := closestY - circleY

	// If the distance is less than or equal to the circle's radius, there is a collision
	return (dx*dx + dy*dy) <= (circleRadius * circleRadius)
}

// Moves the actor towards the target position.
// Calculates the angle between the positive x-axis and the line connecting the actor's current position to the target position
// todetirmine the direction of movement.
func (actor *Actor) MoveTo(targetPosition [2]float64) {
	x := actor.Position[0]
	y := actor.Position[1]

	dx := targetPosition[0] - x
	dy := targetPosition[1] - y

	teta := math.Atan2(dy, dx)

	nextPositionX := actor.Position[0] + math.Cos(teta)*actor.Speed
	nextPositionY := actor.Position[1] + math.Sin(teta)*actor.Speed

	actor.Position[0] = nextPositionX
	actor.Position[1] = nextPositionY

	// if x or y oversteps the target position, set it to the target position (clamping)
	if dx > 0 && nextPositionX > actor.targetPosition[0] {
		actor.Position[0] = actor.targetPosition[0]
	}
	if dx < 0 && nextPositionX < actor.targetPosition[0] {
		actor.Position[0] = actor.targetPosition[0]
	}
	if dy > 0 && nextPositionY > actor.targetPosition[1] {
		actor.Position[1] = actor.targetPosition[1]
	}
	if dy < 0 && nextPositionY < actor.targetPosition[1] {
		actor.Position[1] = actor.targetPosition[1]
	}
}

// Initiates a patrol movement for the actor.
// It patrols between the initial position and a random target position within a specified move range.
// If the actor reaches the target position, it generates a new random target position within the move range.
func (actor *Actor) Patrol(moveRange float64) {
	if actor.Position[0] == actor.targetPosition[0] &&
		actor.Position[1] == actor.targetPosition[1] {
		actor.targetPosition[0] = utils.GetRandomNumInRange(actor.initialPosition[0]-actor.moveRange, actor.initialPosition[0]+actor.moveRange)
		actor.targetPosition[1] = utils.GetRandomNumInRange(actor.initialPosition[1]-actor.moveRange, actor.initialPosition[1]+actor.moveRange)
	}

	actor.MoveTo(actor.targetPosition)
}

// SetLimitBounds sets the limits for the actor's movement.
// It ensures that the actor does not move outside the specified bounds (limiX, limitY).
func (actor *Actor) SetLimitBounds(limiX, limitY float64) {
	playerRect := actor.GetBoundingRect()
	playerWidth := playerRect.Width
	playerHeight := playerRect.Height

	if playerRect.PositionX < 0 {
		actor.Position[0] = 0
	}
	if playerRect.PositionY < 0 {
		actor.Position[1] = 0
	}
	if playerRect.PositionX+playerWidth > limiX {
		actor.Position[0] = limiX - playerWidth
	}
	if playerRect.PositionY+playerHeight > limitY {
		actor.Position[1] = limitY - playerHeight
	}
}
