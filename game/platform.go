package game

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
	maxY := float64(screenHeight - 20)
	//randY := pickOne(minY, y+49)
	randY := float64(rand.Intn(int(maxY)-int(minY))) + minY
	return NewPlatform(x+platformSpacing, randY)
}

func pickOne(a, b float64) float64 {
	if rand.Intn(2) == 0 {
		return a
	}
	return b
}

func (p *Platform) Draw(screen *ebiten.Image, cameraX float64) {
	platformColor := colorGray
	if p.visited {
		platformColor = colorSmartyBlue
	}
	ebitenutil.DrawRect(screen, p.x-cameraX, p.y, p.width, screenHeight-p.y, platformColor)
}
