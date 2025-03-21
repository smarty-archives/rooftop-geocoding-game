package main

import "math"

type Poser interface {
	GetPos() Pos
	GetX() float64
	GetY() float64
	SetPos(p Pos)
	SetX(x float64)
	SetY(y float64)
}

type Pos struct {
	x, y float64
}

func NewPos(x, y float64) *Pos {
	return &Pos{x, y}
}

func (p *Pos) GetPos() Pos {
	return *p
}

func (p *Pos) GetX() float64 {
	return p.x
}

func (p *Pos) GetY() float64 {
	return p.y
}

func (p *Pos) SetPos(pos Pos) {
	*p = pos
}

func (p *Pos) SetX(x float64) {
	p.x = x
}

func (p *Pos) SetY(y float64) {
	p.y = y
}

func (p *Pos) GetDist(other Poser) float64 {
	return math.Abs(p.x - other.GetX())
}

func (p *Pos) Reset() {
	p.x = 0
	p.y = 0
}
