package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

const (
	left                       = -1
	right                      = 1
	startingJumpForce          = -12
	startingPlayerAcceleration = 0.2
	startingMaxPlayerSpeed     = 4
)

type HitBox struct {
	width, height float64
}

type Player struct {
	Pos
	Stats
	HitBox
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
	p.velocityX = 0
	p.velocityY = 0
	p.isJumping = false
	p.jumpForce = startingJumpForce
	p.playerAcceleration = startingPlayerAcceleration
	p.maxPlayerSpeed = startingMaxPlayerSpeed
	p.width = 20
	p.height = playerSize
	image, err := media.Instance.LoadRunningImage(0)
	if err != nil {
		panic(err)
	}
	p.image = image
	p.animation = &IdleState{}
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
	//if debugMode {
	//	p.DrawHitBox(screen, cameraX)
	//}
	playerCoor := &ebiten.DrawImageOptions{}
	scaleX := playerSize / float64(p.image.Bounds().Dx())
	scaleY := scaleX
	x := p.x - cameraX
	if p.velocityX < 0 {
		scaleX = -scaleX
		x += playerSize
	}
	playerCoor.GeoM.Scale(scaleX, scaleY)
	playerCoor.GeoM.Translate(x, p.y)
	screen.DrawImage(p.image, playerCoor)
}

func (p *Player) DrawHitBox(screen *ebiten.Image, cameraX float64) {
	vector.DrawFilledRect(screen, float32(p.x-cameraX+p.width/2), float32(p.y), float32(p.width), float32(p.height), color.RGBA{R: 255}, false)
}

func (p *Player) cycleImage() {
	p.animation = p.animation.getNextState(p)
	if !p.animation.shouldAnimate() {
		return
	}
	p.animation.doState(p)
}

func (p *Player) LeftX() float64 {
	return p.x + playerSize/2 - p.width/2
}

func (p *Player) RightX() float64 {
	return p.x + playerSize/2 + p.width/2
}

func (p *Player) GetCenterX() float64 {
	return p.x + playerSize/2
}
