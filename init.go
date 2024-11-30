package main

import "github.com/hajimehoshi/ebiten/v2/audio"

func init() {
	// Initialize dots
	// If dev is true, generate test dots
	// Otherwise, create dots
	if dev {
		Dots = generateTestDots()
	} else {
		Dots = generateDots()
	}
	fontFace = generateGameFont()
	audioContext := audio.NewContext(sampleRate)
	game = &Game{
		pacman: Pacman{
			x:      pacmanPos.x,
			y:      pacmanPos.y,
			radius: pacmanRadius,
			angle:  0,
			color:  yellow,
		},
		cage:              Cage,
		mainContext:       audioContext,
		ghost:             append([]Pacman{}, Ghost[:]...),
		direction:         None,
		walls:             Walls[:],
		introMusicPlaying: true,
		wallSize:          wallThickness,
		livesLeft:         lives,
		backgroundContext: audioContext,
		audioPlayers:      make(map[string]*audio.Player),
		points:            0,
		level:             1,
	}

	// Initialize intro music
	// a
	player, err := game.getAudioPlayer(audioDir + "intro.wav")
	if err == nil {
		game.mainPlayer = player
		player.Play()
	}
}
