package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	lru "github.com/hashicorp/golang-lru/v2"
)

// AudioData represents cached audio data
type AudioData struct {
	data []byte
	mu   sync.Mutex
}

// NewReader creates a new io.ReadSeeker from the cached data
func (ad *AudioData) NewReader() io.ReadSeeker {
	ad.mu.Lock()
	defer ad.mu.Unlock()
	return bytes.NewReader(ad.data)
}

var audioCache *lru.Cache[string, *AudioData]

func init() {
	var err error
	audioCache, err = lru.New[string, *AudioData](100) // Set a capacity of 100 items
	if err != nil {
		log.Fatalf("Failed to initialize audio cache: %v", err)
	}
}

func initializeAudioPlayerWithContext(path string, context *audio.Context) (*audio.Player, error) {
	// Check cache
	if cached, ok := audioCache.Get(path); ok {
		log.Println("Audio player cache hit")
		reader := cached.NewReader()
		decodedWav, err := wav.Decode(context, reader)
		if err != nil {
			return nil, fmt.Errorf("failed to decode cached WAV data: %w", err)
		}

		audioPlayer, err := context.NewPlayer(decodedWav)
		if err != nil {
			return nil, fmt.Errorf("failed to create audio player from cache: %w", err)
		}
		return audioPlayer, nil
	}

	// Load and decode the audio file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	// Read the entire file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio file: %w", err)
	}

	// Cache the raw audio data
	audioCache.Add(path, &AudioData{data: data})

	// Create a new reader from the data for immediate use
	reader := bytes.NewReader(data)
	decodedWav, err := wav.Decode(context, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode WAV file: %w", err)
	}

	// Create a new audio player
	audioPlayer, err := context.NewPlayer(decodedWav)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio player: %w", err)
	}

	return audioPlayer, nil
}
