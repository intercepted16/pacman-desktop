package main

import (
	"image/color"
	"math"
	"math/rand"
)

type Point struct {
	x, y float64
}

// Helper to get player's current direction
func (g *Game) getPlayerDirection() Point {
	// This should return the player's current movement direction
	// You'll need to track this in your game state
	return Point{
		x: math.Cos(g.pacman.angle),
		y: math.Sin(g.pacman.angle),
	}
}

// Constants for ghost personality types
const (
	CHASER = 0 // Directly pursues Pacman (Blinky behavior)
	AMBUSH = 1 // Tries to get ahead of Pacman (Pinky behavior)
	PATROL = 2 // Moves between Pacman and scatter point (Inky behavior)
	RANDOM = 3 // Semi-random movement with periodic targeting (Clyde behavior)
)

// Direction constants for movement options
var directions = []Point{
	{x: 1, y: 0},           // Right
	{x: -1, y: 0},          // Left
	{x: 0, y: 1},           // Down
	{x: 0, y: -1},          // Up
	{x: 0.707, y: 0.707},   // Diagonal down-right
	{x: -0.707, y: 0.707},  // Diagonal down-left
	{x: 0.707, y: -0.707},  // Diagonal up-right
	{x: -0.707, y: -0.707}, // Diagonal up-left
}

// Vector helper functions
func normalize(dx, dy float64) (float64, float64) {
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		return 0, 0
	}
	return dx / dist, dy / dist
}

func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

type Pacman struct {
	x, y    float64
	radius  float64
	angle   float64
	variety int
	scatter Point
	speed   float64
	lastDir Point // Store last movement direction to prevent zigzagging
	color   color.Color
}

// Find best available direction towards a target that doesn't hit walls
func (g *Game) findBestDirection(ghost *Pacman, target Point) Point {
	bestDir := Point{x: 0, y: 0}
	bestScore := math.Inf(-1)
	currentDist := distance(ghost.x, ghost.y, target.x, target.y)

	// Check each possible direction
	for _, dir := range directions {
		// Check a few steps ahead for wall collisions
		const lookAhead = 3
		validPath := true
		newX, newY := ghost.x, ghost.y

		for step := 1; step <= lookAhead; step++ {
			checkX := newX + dir.x*ghost.speed*float64(step)
			checkY := newY + dir.y*ghost.speed*float64(step)

			if g.anyCollision(checkX, checkY) {
				validPath = false
				break
			}
		}

		if !validPath {
			continue
		}

		// Calculate how good this direction is
		newX = ghost.x + dir.x*ghost.speed
		newY = ghost.y + dir.y*ghost.speed
		newDist := distance(newX, newY, target.x, target.y)

		// Score based on distance improvement and direction consistency
		score := currentDist - newDist

		// Prefer current direction to prevent zigzagging
		if dir.x*ghost.lastDir.x+dir.y*ghost.lastDir.y > 0 {
			score += 0.5
		}

		if score > bestScore {
			bestScore = score
			bestDir = dir
		}
	}

	// If no valid direction found, try to find any valid direction
	if bestScore == math.Inf(-1) {
		for _, dir := range directions {
			newX := ghost.x + dir.x*ghost.speed
			newY := ghost.y + dir.y*ghost.speed

			if !g.collidesWithWall(newX, newY) {
				return dir
			}
		}
	}

	return bestDir
}

// Main AI function for ghost movement
func (g *Game) ghostAi(pacmen []Pacman) {
	for i := range pacmen {
		p := &pacmen[i]

		// Normal movement logic for ghosts outside cage
		target := g.getGhostTarget(p, Point{x: g.pacman.x, y: g.pacman.y})
		bestDir := g.findBestDirection(p, target)

		newX := p.x + bestDir.x*p.speed
		newY := p.y + bestDir.y*p.speed

		// Collision avoidance with other ghosts
		const minSeparation = 30.0
		for j := range pacmen {
			if i != j {
				other := &pacmen[j]
				dist := distance(newX, newY, other.x, other.y)

				if dist < minSeparation {
					repulsionX, repulsionY := normalize(newX-other.x, newY-other.y)
					strength := (minSeparation - dist) / minSeparation
					proposedX := newX + repulsionX*strength*2
					proposedY := newY + repulsionY*strength*2

					if !g.anyCollision(proposedX, proposedY) {
						newX = proposedX
						newY = proposedY
					}
				}
			}
		}

		if !g.anyCollision(newX, newY) {
			p.x = newX
			p.y = newY
			p.lastDir = bestDir
		}
	}
}

// Calculate target position based on ghost personality
func (g *Game) getGhostTarget(ghost *Pacman, player Point) Point {
	switch ghost.variety {
	case CHASER:
		// Directly chase the player if path available, otherwise target nearby accessible position
		if !g.collidesWithWall(player.x, player.y) {
			return player
		}
		// Find closest accessible point to player
		bestDist := math.Inf(1)
		bestPoint := ghost.scatter

		for _, dir := range directions {
			checkX := player.x + dir.x*50
			checkY := player.y + dir.y*50

			if !g.collidesWithWall(checkX, checkY) {
				dist := distance(ghost.x, ghost.y, checkX, checkY)
				if dist < bestDist {
					bestDist = dist
					bestPoint = Point{x: checkX, y: checkY}
				}
			}
		}
		return bestPoint

	case AMBUSH:
		// Try to get ahead of the player, accounting for walls
		const ambushDistance = 80.0
		playerDir := g.getPlayerDirection()
		// targetX := player.x + playerDir.x*ambushDistance
		// targetY := player.y + playerDir.y*ambushDistance

		// If target is in wall, reduce distance until valid
		for d := ambushDistance; d > 0; d -= 10 {
			checkX := player.x + playerDir.x*d
			checkY := player.y + playerDir.y*d
			if !g.collidesWithWall(checkX, checkY) {
				return Point{x: checkX, y: checkY}
			}
		}
		return ghost.scatter

	case PATROL:
		dist := distance(ghost.x, ghost.y, player.x, player.y)
		if dist < 150 {
			// Find valid intermediate point
			midX := (player.x + ghost.scatter.x) / 2
			midY := (player.y + ghost.scatter.y) / 2

			if !g.collidesWithWall(midX, midY) {
				return Point{x: midX, y: midY}
			}
		}
		return player

	case RANDOM:
		if rand.Float64() < 0.02 {
			// Find random valid position
			for attempts := 0; attempts < 10; attempts++ {
				angle := rand.Float64() * 2 * math.Pi
				targetX := ghost.x + math.Cos(angle)*100
				targetY := ghost.y + math.Sin(angle)*100

				if !g.collidesWithWall(targetX, targetY) {
					return Point{x: targetX, y: targetY}
				}
			}
		}

		dist := distance(ghost.x, ghost.y, player.x, player.y)
		if dist > 200 {
			return player
		}
		return ghost.scatter
	}

	return player
}
