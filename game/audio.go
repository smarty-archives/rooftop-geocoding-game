package game

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

// Global variables for audio state
var (
	audioContext *audio.Context
	player       *audio.Player
	isMuted      bool
)

const sampleRate = 44100 // Standard audio sample rate

// InitializeAudio initializes the audio context and player with the provided MP3 data
func InitializeAudio(bgmData []byte) {
	audioContext = audio.NewContext(sampleRate)

	// Decode MP3 from the provided byte slice (bgmData)
	stream, err := mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(bgmData))
	if err != nil {
		log.Fatal(err)
	}

	loop := audio.NewInfiniteLoop(stream, stream.Length())
	player, err = audioContext.NewPlayer(loop)
	if err != nil {
		log.Fatal(err)
	}

	// Play the audio
	player.Play()
}

// ToggleMute toggles the mute state of the audio player
func ToggleMute() {
	if isMuted {
		// Unmute audio
		player.SetVolume(1)
		isMuted = false
	} else {
		// Mute audio
		player.SetVolume(0)
		isMuted = true
	}
}
