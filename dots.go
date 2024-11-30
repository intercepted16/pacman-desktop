package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (d *Dot) Draw(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, d.x, d.y, d.radius, d.color)
}

type Dot struct {
	x, y   float64
	radius float64
	color  color.Color
}

func generateDots() []Dot {
	dots := []Dot{}

	spacing := 30.0 // Space between dots
	radius := 3.0   // Size of dots
	buffer := 10.0  // Minimum distance away from walls

	for x := spacing; x <= screenWidth-spacing; x += spacing {
		for y := spacing; y <= screenHeight-spacing; y += spacing {
			// Check if point is in cage
			inCage := x >= cageX && x <= cageX+cageWidth &&
				y >= cageY && y <= cageY+cageHeight

			// Check if point is near any wall, including cage walls with buffer distance

			allWalls := append(Walls, Cage.Bottom, Cage.Left, Cage.Right, Cage.Top)

			nearWall := false
			for _, wall := range allWalls {
				if x >= wall.x-buffer && x <= wall.x+wall.Width+buffer &&
					y >= wall.y-buffer && y <= wall.y+wall.Height+buffer {
					nearWall = true
					break
				}
			}

			// If not in cage and not near wall, add dot
			if !inCage && !nearWall {
				dots = append(dots, Dot{
					x:      x,
					y:      y,
					radius: radius,
					color:  color.White,
				})
			}
		}
	}
	return dots
}

func generateTestDots() []Dot {
	dots := []Dot{
		{30, 30, 3, color.White},
	}
	return dots
}
