package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
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
	backgroundSampleRate = 44100 / 2
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

func (g *Game) getAudioPlayer(filename string) (*audio.Player, error) {
	g.audioMux.RLock()
	player, exists := g.audioPlayers[filename]
	g.audioMux.RUnlock()

	if exists {
		return player, nil
	}

	g.audioMux.Lock()
	defer g.audioMux.Unlock()

	if player, exists = g.audioPlayers[filename]; exists {
		return player, nil
	}

	newPlayer, err := initializeAudioPlayerWithContext(filename, g.mainContext)
	if err != nil {
		return nil, err
	}

	g.audioPlayers[filename] = newPlayer
	return newPlayer, nil
}

func repositionGhost(p *Pacman, g *Game) {
	for {
		p.x = rand.Float64() * screenWidth
		p.y = rand.Float64() * screenHeight
		if !g.anyCollision(p.x, p.y) {
			break
		}
	}
}

func measureText(textToDisplay string) (x, y int) {
	generateGameFont()

	fontMutex.RLock()
	bounds := text.BoundString(fontFace, textToDisplay)
	fontMutex.RUnlock()

	return bounds.Max.X - bounds.Min.X, bounds.Max.Y - bounds.Min.Y
}

func drawText(screen *ebiten.Image, point Point, textToDisplay string, textColor color.Color) {
	generateGameFont()

	fontMutex.RLock()
	text.Draw(screen, textToDisplay, fontFace, int(point.x), int(point.y), textColor)
	fontMutex.RUnlock()
}

func (g *Game) respawnPacman() {
	g.pacman.x = pacmanPos.x
	g.pacman.y = pacmanPos.y
	g.direction = None
	g.introMusicPlaying = true

	player, err := g.getAudioPlayer(audioDir + "intro.wav")
	if err == nil {
		g.mainPlayer = player
		g.mainPlayer.Play()
	}

	for i := range g.ghost {
		// reset every ghost
		repositionGhost(&g.ghost[i], g)
	}

}

func (g *Game) Update() error {
	if g.gameOverState {
		return nil
	}
	if g.introMusicPlaying {
		if !g.mainPlayer.IsPlaying() {
			g.introMusicPlaying = false
			for i := range g.ghost {
				repositionGhost(&g.ghost[i], g)
			}
		}
		return nil
	}
	speed := 2.0

	if g.backgroundPlayer == nil {
		player, err := g.getAudioPlayer(audioDir + "siren.wav")
		if err == nil {
			g.backgroundPlayer = player
		}
	}

	if g.backgroundPlayer != nil && !g.backgroundPlayer.IsPlaying() {
		g.backgroundPlayer.Seek(0)
		g.backgroundPlayer.Play()
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.direction = Up
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.direction = Down
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.direction = Left
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.direction = Right
	}

	switch g.direction {
	case Up:
		if g.pacman.y-speed >= g.pacman.radius && !g.anyCollision(g.pacman.x, g.pacman.y-speed) {
			g.pacman.y -= speed
		}
	case Down:
		if g.pacman.y+speed <= screenHeight-g.pacman.radius && !g.anyCollision(g.pacman.x, g.pacman.y+speed) {
			g.pacman.y += speed
		}
	case Left:
		if g.pacman.x-speed >= g.pacman.radius && !g.anyCollision(g.pacman.x-speed, g.pacman.y) {
			g.pacman.x -= speed
		}
	case Right:
		if g.pacman.x+speed <= screenWidth-g.pacman.radius && !g.anyCollision(g.pacman.x+speed, g.pacman.y) {
			g.pacman.x += speed
		}
	}

	for i := len(Dots) - 1; i >= 0; i-- {
		dot := Dots[i]
		if distance(g.pacman.x, g.pacman.y, dot.x, dot.y) < g.pacman.radius {
			Dots = append(Dots[:i], Dots[i+1:]...)
			g.points += 10

			player, err := g.getAudioPlayer(audioDir + "dot.wav")
			if err == nil && !player.IsPlaying() {
				player.Seek(0)
				player.Play()
			}
		}
	}

	if len(Dots) == 0 {
		println("You win!")
		g.level++
		println("Respwaning for next level")
		player, err := g.getAudioPlayer(audioDir + "intermission.wav")
		if err == nil {
			player.Seek(0)
			player.Play()
			for player.IsPlaying() {
			}
		}
		// Make ghosts faster based on level
		switch {
		case g.level == 5 || g.level == 10 || g.level == 15:
			// Increase speed for levels 5, 10, and 15
			for i := range g.ghost {
				g.ghost[i].speed += 1
			}
		case g.level >= 20:
			// Increase speed by 0.25 for levels 20 and above
			for i := range g.ghost {
				g.ghost[i].speed += 0.25
			}
		}
		g.respawnPacman()
		Dots = generateTestDots()

	}

	for _, ghost := range g.ghost {
		if distance(g.pacman.x, g.pacman.y, ghost.x, ghost.y) < g.pacman.radius {
			g.livesLeft--
			if g.livesLeft == 0 {
				g.gameOverState = true
				player, err := g.getAudioPlayer(audioDir + "gameover.wav")
				if err == nil {
					player.Seek(0)
					player.Play()
				}
			} else {
				println("Lost a life")
				player, err := g.getAudioPlayer(audioDir + "gameover.wav")
				if err == nil {
					player.Seek(0)
					player.Play()
					for player.IsPlaying() {
					}
				}
				g.respawnPacman()
			}
		}
	}

	g.ghostAi(g.ghost)
	return nil
}

func (g *Game) collidesWithWall(x, y float64) bool {
	for _, wall := range g.walls {
		if x+g.pacman.radius > wall.x && x-g.pacman.radius < wall.x+wall.Width &&
			y+g.pacman.radius > wall.y && y-g.pacman.radius < wall.y+wall.Height {
			return true
		}
	}
	return false
}

func (g *Game) collidesWithCage(x, y float64) bool {
	if x+g.pacman.radius > g.cage.Left.x && x-g.pacman.radius < g.cage.Right.x+g.cage.Right.Width &&
		y+g.pacman.radius > g.cage.Top.y && y-g.pacman.radius < g.cage.Bottom.y+g.cage.Bottom.Height {
		return true
	}
	return false
}

func (g *Game) anyCollision(x, y float64) bool {
	return g.collidesWithWall(x, y) || g.collidesWithCage(x, y)
}

func drawCenteredText(screen *ebiten.Image, textToDisplay string, textColor color.Color) {
	textCacheMux.RLock()
	point, exists := textCache[textToDisplay]
	textCacheMux.RUnlock()

	if !exists {
		textWidth, textHeight := measureText(textToDisplay)
		x := (screenWidth - textWidth) / 2
		y := (screenHeight - textHeight) / 2

		textCacheMux.Lock()
		textCache[textToDisplay] = Point{x: float64(x), y: float64(y)}
		textCacheMux.Unlock()

		point = Point{x: float64(x), y: float64(y)}
	}

	drawText(screen, point, textToDisplay, textColor)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.pacman.Draw(screen)
	g.cage.Draw(screen)
	for _, dot := range Dots {
		dot.Draw(screen)
	}
	for _, p := range g.ghost {
		p.Draw(screen)
	}
	for _, wall := range g.walls {
		wall.Draw(screen)
	}
	if g.introMusicPlaying {
		drawCenteredText(screen, "READY!", color.White)
	} else {
		drawCenteredText(screen, "Points: "+strconv.Itoa(g.points), color.White)
	}

	if g.gameOverState {
		drawCenteredText(screen, "GAME OVER", color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	ebiten.SetTPS(30)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("PacMan Desktop")

	pacmanIcon, _, err := ebitenutil.NewImageFromFile(imageDir + "pacman.png")
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowIcon([]image.Image{pacmanIcon})

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
