package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

type Cloud struct {
	Image        *ebiten.Image
	Speed        float64 // Speed multiplier for parallax effect
	Transparency float32
	Pos
}

func NewCloud(x, y float64, speed float64) *Cloud {
	return &Cloud{
		Image:        media.Instance.GetCloudImage(),
		Speed:        speed,
		Pos:          Pos{x, y},
		Transparency: 0, // 100% transparency (although it's not THAT transparent)
	}
}

func (c *Cloud) Draw(screen *ebiten.Image, cameraX float64) {
	cloudCoor := &ebiten.DrawImageOptions{}
	scaleX, scaleY := 1.0, 1.0
	x := c.getOffsetX(cameraX)

	cloudCoor.GeoM.Scale(scaleX, scaleY)
	cloudCoor.GeoM.Translate(x, c.GetY())

	// Set transparency (alpha value between 0.0 and 1.0)
	cloudCoor.ColorScale.Scale(1, 1, 1, c.Transparency)

	screen.DrawImage(c.Image, cloudCoor)
}

func (c *Cloud) getOffsetX(cameraX float64) float64 {
	return c.GetX() - cameraX*c.Speed
}

func pickRand[T any](items ...T) (winner T) {
	if len(items) == 0 {
		return winner
	}
	return items[rand.Intn(len(items))]
}
