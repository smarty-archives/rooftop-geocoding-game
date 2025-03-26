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

func GenerateNewRandomPlatform(prevPlatform *Platform, score int) *Platform {
	x := prevPlatform.x
	y := prevPlatform.y
	minY := max(y-maxYDeltaTop, playerSize+jumpApexHeight)
	maxY := float64(screenHeight - minimumPlatformHeight)
	randX := x + prevPlatform.width + giveOrTake(platformSpacing, 50)
	randY := float64(rand.Intn(int(maxY)-int(minY))) + minY
	randWidth := pickWidth(score, 175, 150, 125, 100, 75)
	//randWidth := giveOrTake(150, 75)
	return NewPlatform(randX, randY, randWidth)
}

func pickWidth(counter int, numbers ...float64) float64 {
	// Define weight distribution based on counter value
	weights := []float64{0.5, 0.4, 0.1, 0, 0} // Initial probabilities

	// As counter increases, shift probability towards later numbers
	if counter > 50 {
		weights = []float64{0.37, 0.3, 0.3, 0.02, 0.01}
	}
	if counter > 100 {
		weights = []float64{0.5, 0.3, 0.1, 0.05, 0.05}
	}
	if counter > 150 {
		weights = []float64{0.3, 0.2, 0.2, 0.15, 0.15}
	}
	if counter > 200 {
		weights = []float64{0.15, 0.15, 0.25, 0.25, 0.2}
	}
	if counter > 250 {
		weights = []float64{0.05, 0.05, 0.3, 0.3, 0.3}
	}
	if counter > 300 {
		weights = []float64{0.01, 0.02, 0.3, 0.3, 0.37}
	}

	// Pick a number based on weighted probabilities
	r := rand.Float64()
	sum := 0.0
	for i, w := range weights {
		sum += w
		if r < sum {
			return numbers[i]
		}
	}
	return numbers[len(numbers)-1] // Fallback, should never happen
}

//func pickRand[T any](items ...T) (winner T) {
//	if len(items) == 0 {
//		return winner
//	}
//	return items[rand.Intn(len(items))]
//}

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
