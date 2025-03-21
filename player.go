package main

const (
	left  = -1
	right = 1
)

type Stats struct {
	jumpForce      float64
	playerSpeed    float64
	maxPlayerSpeed float64
}

func (p *Stats) SetJumpForce(jumpForce float64) {
	p.jumpForce = jumpForce
}

func (p *Stats) SetPlayerSpeed(playerSpeed float64) {
	p.playerSpeed = playerSpeed
}

func (p *Stats) SetMaxPlayerSpeed(maxPlayerSpeed float64) {
	p.maxPlayerSpeed = maxPlayerSpeed
}

func (p *Stats) GetJumpForce() float64 { return p.jumpForce }

func (p *Stats) GetPlayerSpeed() float64 { return p.playerSpeed }

func (p *Stats) GetMaxPlayerSpeed() float64 { return p.maxPlayerSpeed }

type Player struct {
	Pos
	Stats
	velocityX float64
	velocityY float64
	isJumping bool
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
