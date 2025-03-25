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
	wasmAudio()
	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

func wasmAudio() {
	js.Global().Get("document").Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if player == nil || !player.IsPlaying() {
			playAudio()
		}
		return nil
	}))
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
