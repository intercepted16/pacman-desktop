package main

import (
	"io/ioutil"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Initialize font face once
func generateGameFont() font.Face {
	var face font.Face
	fontFaceOnce.Do(func() {
		tt, err := opentype.Parse(getFontBytes(retroFont))
		if err != nil {
			log.Fatal(err)
		}
		face, err = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    24,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			log.Fatal(err)
		}
	})
	return face
}

// Function to open TTF file and get bytes
func getFontBytes(filePath string) []byte {
	if fontBytes, ok := fontBytesCache.Get(filePath); ok {
		// println("Cache hit")
		return fontBytes
	}
	println("Cache miss")
	fontBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	fontBytesCache.Add(filePath, fontBytes)
	return fontBytes
}
