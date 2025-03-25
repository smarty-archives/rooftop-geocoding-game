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

func NewPlatform(x, y, width float64) *Platform {
	return &Platform{
		Pos:   *NewPos(x, y),
		width: width,
	}
}

func GenerateNewRandomPlatform(prevPlatform *Platform) *Platform {
	x := prevPlatform.x
	y := prevPlatform.y
	minY := max(y-maxYDeltaTop, playerSize+jumpApexHeight)
	maxY := float64(screenHeight - minimumPlatformHeight)
	randX := x + prevPlatform.width + giveOrTake(platformSpacing, 50)
	randY := float64(rand.Intn(int(maxY)-int(minY))) + minY
	randWidth := giveOrTake(150, 75)
	return NewPlatform(randX, randY, randWidth)
}

func giveOrTake(num, delta float64) float64 {
	return float64(rand.Intn(int(num+delta)-int(num-delta))) + num - delta
}

func (p *Platform) Draw(screen *ebiten.Image, cameraX float64) {
	platformColor := colorGray
	if p.visited {
		platformColor = colorSmartyBlue
	}
	ebitenutil.DrawRect(screen, p.x-cameraX, p.y, p.width, screenHeight-p.y, platformColor)
}
