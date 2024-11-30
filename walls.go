package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Wall struct {
	x, y, Width, Height float64
	Color               color.Color
}

var Walls = Map([]Wall{
	// Outer walls
	{x: 0, y: 0, Width: screenWidth, Height: wallThickness},                            // Top
	{x: 0, y: 0, Width: wallThickness, Height: screenHeight},                           // Left
	{x: screenWidth - wallThickness, y: 0, Width: wallThickness, Height: screenHeight}, // Right
	{x: 0, y: screenHeight - wallThickness, Width: screenWidth, Height: wallThickness}, // Bottom

	// Inner walls
	// 'L' shape on the left side
	{x: WallMinimumOffset, y: WallMinimumOffset, Width: WallWidth, Height: wallThickness}, // Horizontal part
	{x: WallMinimumOffset, y: WallMinimumOffset, Width: wallThickness, Height: WallWidth}, // Vertical part

	// 'L' shape on the right side
	{x: screenWidth - WallMinimumOffset - WallWidth, y: WallMinimumOffset, Width: WallWidth, Height: wallThickness},     // Horizontal part
	{x: screenWidth - WallMinimumOffset - wallThickness, y: WallMinimumOffset, Width: wallThickness, Height: WallWidth}, // Vertical part

	// 'L' shape on the bottom left side
	{x: WallMinimumOffset, y: screenHeight - WallMinimumOffset - wallThickness, Width: WallWidth, Height: wallThickness}, // Horizontal part
	{x: WallMinimumOffset, y: screenHeight - WallMinimumOffset - WallWidth, Width: wallThickness, Height: WallWidth},     // Vertical part

	// 'L' shape on the bottom right side
	{x: screenWidth - WallMinimumOffset - WallWidth, y: screenHeight - WallMinimumOffset - wallThickness, Width: WallWidth, Height: wallThickness}, // Horizontal part
	{x: screenWidth - WallMinimumOffset - wallThickness, y: screenHeight - WallMinimumOffset - WallWidth, Width: wallThickness, Height: WallWidth}, // Vertical part

	// Offset 60 from the cage walls horizontally
	{x: cageX - WallMinimumOffset - wallThickness, y: cageY, Width: wallThickness, Height: cageHeight}, // Left
	{x: cageX + cageWidth + WallMinimumOffset, y: cageY, Width: wallThickness, Height: cageHeight},     // Right

	// Offset 60 from the cage walls vertically
	{x: cageX, y: cageY - WallMinimumOffset - wallThickness, Width: cageWidth, Height: wallThickness}, // Top
	{x: cageX, y: cageY + cageHeight + WallMinimumOffset, Width: cageWidth, Height: wallThickness},    // Bottom

}, func(w Wall) Wall {
	colors := []color.Color{yellow, green, blue, lightBlue, red}
	w.Color = colors[rand.Intn(len(colors))]
	return w
})

func (w *Wall) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, w.x, w.y, w.x+w.Width, w.y, w.Color)                   // Top
	ebitenutil.DrawLine(screen, w.x, w.y, w.x, w.y+w.Height, w.Color)                  // Left
	ebitenutil.DrawLine(screen, w.x+w.Width, w.y, w.x+w.Width, w.y+w.Height, w.Color)  // Right
	ebitenutil.DrawLine(screen, w.x, w.y+w.Height, w.x+w.Width, w.y+w.Height, w.Color) // Bottom
}
