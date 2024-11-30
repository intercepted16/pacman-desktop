package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (p *Pacman) Draw(screen *ebiten.Image) {
	// Draw the body of Pacman
	ebitenutil.DrawCircle(screen, p.x, p.y, p.radius, p.color)

	// Draw the mouth of Pacman
	mouthAngle := math.Pi / 4
	for i := -mouthAngle; i <= mouthAngle; i += 0.01 {
		x := p.x + p.radius*math.Cos(p.angle+i)
		y := p.y + p.radius*math.Sin(p.angle+i)
		screen.Set(int(x), int(y), color.White)
	}
}
