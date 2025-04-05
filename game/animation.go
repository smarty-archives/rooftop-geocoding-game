package game

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

const (
	hesRunning = .8
)

type BaseAnimationState struct {
	animationCounter int
	imageNum         int
}

func (state *BaseAnimationState) getImageNum() int {
	return state.imageNum
}

func (state *BaseAnimationState) shouldAnimate() bool {
	if state.animationCounter > 0 {
		state.animationCounter--
		return false
	}
	state.animationCounter = 6
	return true
}

type PlayerAnimationState interface {
	getNextState(*Player) PlayerAnimationState
	doState(*Player)
	getImageNum() int
	shouldAnimate() bool
	getImage() *ebiten.Image
}

type IdleState struct {
	BaseAnimationState
}

func (state *IdleState) getNextState(player *Player) PlayerAnimationState {
	if player.isJumping {
		return &JumpingState{}
	}
	if math.Abs(player.velocityX) < hesRunning {
		return state
	}
	return &RunningState{}
}

func (state *IdleState) doState(player *Player) {
	state.imageNum++
	if state.imageNum >= media.NumIdleImages {
		state.imageNum = 0
	}
	player.image = state.getImage()
}

func (state *IdleState) shouldAnimate() bool {
	if state.animationCounter > 0 {
		state.animationCounter--
		return false
	}
	state.animationCounter = 24
	return true
}

func (state *IdleState) getImage() *ebiten.Image {
	image, err := media.Instance.LoadIdleImage(state.getImageNum())
	if err != nil {
		log.Fatal(err)
	}
	return image
}

type RunningState struct {
	BaseAnimationState
}

func (state *RunningState) getImage() *ebiten.Image {
	image, err := media.Instance.LoadRunningImage(state.getImageNum())
	if err != nil {
		log.Fatal(err)
	}
	return image
}

func (state *RunningState) getNextState(player *Player) PlayerAnimationState {
	if player.isJumping {
		return &JumpingState{}
	}
	if math.Abs(player.velocityX) < hesRunning {
		return &IdleState{}
	}
	return state
}

func (state *RunningState) doState(player *Player) {
	state.imageNum++
	if state.imageNum >= media.NumPlayerImages {
		state.imageNum = 0
	}
	player.image = state.getImage()
}

type JumpingState struct {
	BaseAnimationState
}

func (state *JumpingState) getImage() *ebiten.Image {
	image, err := media.Instance.LoadRunningImage(state.getImageNum())
	if err != nil {
		log.Fatal(err)
	}
	return image
}

func (state *JumpingState) getNextState(player *Player) PlayerAnimationState {
	if !player.isJumping {
		return &RunningState{}
	}
	return state
}

func (state *JumpingState) shouldAnimate() bool {
	if state.animationCounter > 0 {
		state.animationCounter--
		return false
	}
	state.animationCounter = 9
	return true
}

func (state *JumpingState) doState(player *Player) {
	state.imageNum++
	if state.imageNum >= media.NumPlayerImages {
		state.imageNum = 0
	}
	player.image = state.getImage()
}
