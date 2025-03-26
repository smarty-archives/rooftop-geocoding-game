package game

type Stats struct {
	jumpForce          float64
	playerAcceleration float64
	maxPlayerSpeed     float64
}

func (p *Stats) SetJumpForce(jumpForce float64) {
	p.jumpForce = jumpForce
}

func (p *Stats) SetPlayerAcceleration(playerAcceleration float64) {
	p.playerAcceleration = playerAcceleration
}

func (p *Stats) SetMaxPlayerSpeed(maxPlayerSpeed float64) {
	p.maxPlayerSpeed = maxPlayerSpeed
}

func (p *Stats) GetJumpForce() float64 { return p.jumpForce }

func (p *Stats) GetPlayerAcceleration() float64 { return p.playerAcceleration }

func (p *Stats) GetMaxPlayerSpeed() float64 { return p.maxPlayerSpeed }
