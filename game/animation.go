package game

import (
	"math"

	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

const (
	hesRunning = .8
)

type BaseAnimationState struct {
	animationCounter int
	imageNum         int
}

func (bas *BaseAnimationState) getImageNum() int {
	return bas.imageNum
}

func (bas *BaseAnimationState) shouldAnimate() bool {
	if bas.animationCounter > 0 {
		bas.animationCounter--
		return false
	}
	bas.animationCounter = 10
	return true
}

type PlayerAnimationState interface {
	getNextState(*Player) PlayerAnimationState
	doState(*Player)
	getImageNum() int
	shouldAnimate() bool
}

type IdleState struct {
	BaseAnimationState
}

func (state *IdleState) getNextState(player *Player) PlayerAnimationState {
	if math.Abs(player.velocityX) < hesRunning {
		return state
	}
	return &RunningState{}
}

func (state *IdleState) doState(player *Player) {
}

type RunningState struct {
	BaseAnimationState
}

func (state *RunningState) getNextState(player *Player) PlayerAnimationState {
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
}
