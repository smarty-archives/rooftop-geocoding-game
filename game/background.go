package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

type Layer struct {
	Image   *ebiten.Image
	Speed   float64 // Speed multiplier for parallax effect
	OffsetX float64
}

func NewLayers() []Layer {
	var layers []Layer
	speeds := []float64{0.2, 0.5} // Warning: if there are fewer speeds than background layers, then this will panic
	for i := range media.NumBackgroundLayers {
		img, err := media.Instance.LoadBackgroundImage(i)
		if err != nil {
			log.Fatal(err)
		}
		layers = append(layers, Layer{
			Image: img,
			Speed: speeds[i],
		})
	}
	return layers
}
