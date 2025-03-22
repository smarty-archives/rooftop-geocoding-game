package game

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
