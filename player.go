package main

type Player struct {
	Pos
	xSpeed    float64
	velocityY float64
	isJumping bool
}

func (p *Player) Reset() {
	// todo reset pos
	p.Pos.Reset()
	p.xSpeed = 0
	p.velocityY = 0
	p.isJumping = false
}
