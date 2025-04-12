package game

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	minLat = 30.0
	maxLat = 50.0
	minLon = -120.0
	maxLon = -70.0
)

type Geocode struct {
	str     string
	opacity int
	Pos
}

func NewGeocode(pos Pos) *Geocode {
	lat, lon := randomGeocode()
	return &Geocode{
		str:     fmt.Sprintf("%f, %f", lat, lon),
		opacity: 300,
		Pos:     pos,
	}
}

func (g *Geocode) Draw(screen *ebiten.Image, fontObj font.Face, cameraX float64) {
	textWidth := font.MeasureString(fontObj, g.str).Ceil()
	textHeight := fontObj.Metrics().Ascent.Ceil()
	drawX := int(g.x) - textWidth/2
	drawY := int(g.y) - textHeight/2
	opacity := min(g.opacity, 255)
	text.Draw(screen, g.str, fontObj, drawX-int(cameraX), drawY, color.RGBA{A: uint8(opacity)})
}

func randomGeocode() (float32, float32) {
	return randomFloat32(minLat, maxLat), randomFloat32(minLon, maxLon)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
