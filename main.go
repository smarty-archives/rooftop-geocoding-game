//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	_ "embed"
	"log"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/smarty-archives/rooftop-geocoding-game/game"
)

//go:embed assets/music/background.mp3
var bgmData []byte

const sampleRate = 44100 // Standard audio sample rate

var (
	audioContext     *audio.Context
	player           *audio.Player
	audioInitialized bool
)

func main() {
	wasmAudio()
	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

func wasmAudio() {
	startAudio := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !audioInitialized {
			audioInitialized = true
			audioContext = audio.NewContext(sampleRate)

			// Decode MP3 from embedded bytes
			stream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(bgmData))
			if err != nil {
				log.Fatal(err)
			}

			loop := audio.NewInfiniteLoop(stream, stream.Length())

			player, err = audio.NewPlayer(audioContext, loop)
			if err != nil {
				log.Fatal(err)
			}

			player.Play()
		}
		return nil
	})

	doc := js.Global().Get("document")
	doc.Call("addEventListener", "click", startAudio)
	doc.Call("addEventListener", "touchstart", startAudio)
}

func playAudio() {
	audioContext = audio.NewContext(sampleRate)

	// Decode MP3 from embedded bytes
	stream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(bgmData))
	if err != nil {
		log.Fatal(err)
	}

	// Wrap the stream in an infinite loop so it repeats
	loop := audio.NewInfiniteLoop(stream, stream.Length())

	// Create an audio player using the looped stream
	player, err = audio.NewPlayer(audioContext, loop)
	if err != nil {
		log.Fatal(err)
	}

	// Play the looped audio
	player.Play()
}
