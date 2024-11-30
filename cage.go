package main

import "github.com/hajimehoshi/ebiten/v2"

func (c *Square) Draw(screen *ebiten.Image) {
	c.Top.Draw(screen)
	c.Right.Draw(screen)
	c.Bottom.Draw(screen)
	c.Left.Draw(screen)
}
