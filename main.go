//go:build js && wasm

package main

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/smarty-archives/rooftop-geocoding-game/game"

	"syscall/js"
)

//go:embed assets/music/background-loop-melodic-techno-02-2690.mp3
var bgmData []byte

const sampleRate = 44100 // Standard audio sample rate

var (
	audioContext *audio.Context
	player       *audio.Player
)

func main() {
	js.Global().Get("document").Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if player == nil || !player.IsPlaying() {
			playAudio()
		}
		return nil
	}))
	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

func playAudio() {
	audioContext = audio.NewContext(sampleRate)

	// Decode MP3 from embedded bytes
	stream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(bgmData))
	if err != nil {
		log.Fatal(err)
	}

	// Create an audio player
	player, err = audio.NewPlayer(audioContext, stream)
	if err != nil {
		log.Fatal(err)
	}

	// Autoplay won't work in browsers; JavaScript needs to trigger play
	player.Rewind()
	player.Play()
}
