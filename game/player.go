package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

const (
	left                       = -1
	right                      = 1
	startingJumpForce          = -12
	startingPlayerAcceleration = 0.2
	startingMaxPlayerSpeed     = 4
)

type Player struct {
	Pos
	Stats
	velocityX float64
	velocityY float64
	isJumping bool
	animation PlayerAnimationState
	image     *ebiten.Image
}

func NewPlayer() *Player {
	p := &Player{}
	p.ResetPlayer()
	return p
}

func (p *Player) ResetPlayer() {
	p.x = screenWidth/2 - (playerSize / 2)
	p.y = 0
	p.jumpForce = startingJumpForce
	p.playerAcceleration = startingPlayerAcceleration
	p.maxPlayerSpeed = startingMaxPlayerSpeed
	image, err := media.Instance.LoadPlayerImage(0)
	if err != nil {
		panic(err)
	}
	p.image = image
	p.animation = &IdleState{}
}

func (p *Player) Reset() {
	p.Pos.Reset()
	p.velocityX = 0
	p.velocityY = 0
	p.isJumping = false
}

func (p *Player) Jump() {
	p.velocityY = p.GetJumpForce()
	p.isJumping = true
}

func (p *Player) Accelerate(dir float64) {
	if dir*p.velocityX < 0 {
		p.velocityX = 0
	}
	if dir*p.velocityX <= p.GetMaxPlayerSpeed() {
		p.velocityX += dir * p.GetPlayerAcceleration()
	}
	p.x += p.velocityX
}

func (p *Player) AccelerateLeft() { p.Accelerate(left) }

func (p *Player) AccelerateRight() { p.Accelerate(right) }

func (p *Player) SetStats(jumpForce, playerSpeed, maxPlayerSpeed float64) {
	p.SetJumpForce(jumpForce)
	p.SetPlayerAcceleration(playerSpeed)
	p.SetMaxPlayerSpeed(maxPlayerSpeed)
}

func (p *Player) Draw(screen *ebiten.Image, cameraX float64) {
	p.cycleImage()
	playerCoor := &ebiten.DrawImageOptions{}
	scaleX := playerSize / float64(p.image.Bounds().Dx())
	scaleY := playerSize / float64(p.image.Bounds().Dy())
	x := p.x - cameraX
	if p.velocityX < 0 {
		scaleX = -scaleX
		//playerCoor.GeoM.Translate(playerSize, 0)
		x += playerSize
	}
	playerCoor.GeoM.Scale(scaleX, scaleY)
	playerCoor.GeoM.Translate(x, p.y)
	screen.DrawImage(p.image, playerCoor)
}

func (p *Player) cycleImage() {
	if !p.animation.shouldAnimate() {
		return
	}
	p.animation = p.animation.getNextState(p)
	p.animation.doState(p)

	image, err := media.Instance.LoadPlayerImage(p.animation.getImageNum())
	if err != nil {
		log.Fatal(err)
	}
	p.image = image
}
