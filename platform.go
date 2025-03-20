package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Platform struct {
	Pos
	width   float64
	visited bool
}

func NewPlatform(x, y float64) *Platform {
	return &Platform{
		Pos:   *NewPos(x, y),
		width: platformWidth,
	}
}

func GenerateNewRandomPlatform(prevPlatform *Platform) *Platform {
	x := prevPlatform.x
	y := prevPlatform.y
	minY := max(y-maxYDeltaTop, playerSize+jumpApexHeight)
	maxY := int(min(screenHeight-20, y+maxYDeltaBottom))
	randY := float64(rand.Intn(maxY-int(minY))) + minY
	return NewPlatform(x+platformSpacing, randY)
}

func (p *Platform) Draw(screen *ebiten.Image, cameraX float64) {
	platformColor := colorGray
	if p.visited {
		platformColor = colorSmartyBlue
	}
	ebitenutil.DrawRect(screen, p.x-cameraX, p.y, p.width, screenHeight-p.y, platformColor)
}
