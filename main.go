//go:build js && wasm
// +build js,wasm

package main

import (
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smarty-archives/rooftop-geocoding-game/game"
)

//go:embed assets/music/background.mp3
var bgmData []byte

func main() {
	// Initialize audio in the main file with the embedded bgmData
	game.InitializeAudio(bgmData)

	// Initialize the game
	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
