package main

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font"
)

const (
	assetsDir = "assets"
	imageDir  = assetsDir + "/img/"
	fontDir   = assetsDir + "/font/"
	audioDir  = assetsDir + "/audio/"
	retroFont = fontDir + "/retro.ttf"
)
const (
	Up = 1 << iota
	Down
	Left
	Right
)

var game *Game

var Dots []Dot

const dev = true

var Cage = Square{
	Top:    Wall{x: cageX, y: cageY, Width: cageWidth, Height: wallThickness, Color: blue},
	Right:  Wall{x: cageX + cageWidth - wallThickness, y: cageY, Width: wallThickness, Height: cageHeight, Color: blue},
	Bottom: Wall{x: cageX, y: cageY + cageHeight - wallThickness, Width: cageWidth, Height: wallThickness, Color: blue},
	Left:   Wall{x: cageX, y: cageY, Width: wallThickness, Height: cageHeight, Color: blue},
}

const (
	sampleRate           = 44100
	screenWidth          = 640
	screenHeight         = 480
	backgroundSampleRate = sampleRate / 2
)

var (
	yellow    = color.RGBA{255, 204, 85, 255}
	green     = color.RGBA{0, 153, 76, 255}
	blue      = color.RGBA{0, 0, 255, 0}
	lightBlue = color.RGBA{100, 170, 230, 255}
	red       = color.RGBA{255, 85, 85, 255}
	orange    = color.RGBA{255, 153, 0, 255}
)

const (
	wallThickness float64 = 12
	cageWidth     float64 = 200
	cageHeight    float64 = 150
	pacmanRadius  float64 = 20
	pacmanOffsetY float64 = WallMinimumOffset * 2
)

var (
	pacmanPos = Point{
		x: cageX + cageWidth/2,
		y: cageY + cageHeight + pacmanOffsetY,
	}
)

var cageX = (screenWidth - cageWidth) / 2
var cageY = (screenHeight - cageHeight) / 2

type Direction int

const (
	None              Direction = iota
	WallMinimumOffset           = 60
	WallWidth                   = 100
	lives                       = 3
)

// Cache for font face
var (
	fontFace     font.Face
	fontFaceOnce sync.Once
	fontMutex    sync.RWMutex
)

// Cache for text measurements
var (
	textCache    = make(map[string]Point)
	textCacheMux sync.RWMutex
)

type Game struct {
	pacman            Pacman
	cage              Square
	mainPlayer        *audio.Player
	mainContext       *audio.Context
	ghost             []Pacman
	direction         Direction
	walls             []Wall
	introMusicPlaying bool
	wallSize          float64
	gameOverState     bool
	livesLeft         int
	backgroundPlayer  *audio.Player
	backgroundContext *audio.Context
	points            int
	audioPlayers      map[string]*audio.Player
	audioMux          sync.RWMutex
	level             int
}
