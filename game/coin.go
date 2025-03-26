package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

type Coin struct {
	Image          *ebiten.Image
	AnimationSpeed float64
}

func NewCoin() *Coin {
	img, err := media.Instance.LoadCoinImage()
	if err != nil {
		log.Fatal(err)
	}
	return &Coin{
		Image:          img,
		AnimationSpeed: 5,
	}
}
