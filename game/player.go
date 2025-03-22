package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

const (
	left                   = -1
	right                  = 1
	startingJumpForce      = -12.0
	startingPlayerSpeed    = 0.2
	startingMaxPlayerSpeed = 5.0
)

type Player struct {
	Pos
	Stats
	velocityX        float64
	velocityY        float64
	isJumping        bool
	image            *ebiten.Image
	imageNum         int
	animationCounter int
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
	p.playerSpeed = startingPlayerSpeed
	p.maxPlayerSpeed = startingMaxPlayerSpeed
	image, err := media.Instance.LoadPlayerImage(0)
	if err != nil {
		panic(err)
	}
	p.image = image
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

func (p *Player) Move(dir float64) {
	if dir*p.velocityX < 0 {
		p.velocityX = 0
	}
	if dir*p.velocityX <= p.GetMaxPlayerSpeed() {
		p.velocityX += dir * p.GetPlayerSpeed()
	}
	p.x += p.velocityX
}

func (p *Player) MoveLeft() { p.Move(left) }

func (p *Player) MoveRight() { p.Move(right) }

func (p *Player) SetStats(jumpForce, playerSpeed, maxPlayerSpeed float64) {
	p.SetJumpForce(jumpForce)
	p.SetPlayerSpeed(playerSpeed)
	p.SetMaxPlayerSpeed(maxPlayerSpeed)
}

func (p *Player) Draw(screen *ebiten.Image, cameraX float64) {
	p.cycleImage()
	playerCoor := &ebiten.DrawImageOptions{}
	scaleX := playerSize / float64(p.image.Bounds().Dx())
	scaleY := playerSize / float64(p.image.Bounds().Dy())
	playerCoor.GeoM.Scale(scaleX, scaleY)
	playerCoor.GeoM.Translate(p.x-cameraX, p.y)
	screen.DrawImage(p.image, playerCoor)
}

func (p *Player) cycleImage() {
	if p.animationCounter > 0 {
		p.animationCounter--
		return
	}
	p.animationCounter = 10
	p.imageNum++
	if p.imageNum >= media.NumPlayerImages {
		p.imageNum = 0
	}
	image, err := media.Instance.LoadPlayerImage(p.imageNum)
	if err != nil {
		log.Fatal(err)
	}
	p.image = image
}
