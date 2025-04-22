package game

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
)

type Platform struct {
	Pos
	image              *ebiten.Image
	visitedImage       *ebiten.Image
	width              float64
	framesSinceVisited int
	visited            bool
}

func NewPlatform(x, y, width float64) *Platform {
	img, err := media.Instance.LoadBuildingImage(getPlatformIndex(width))
	if err != nil {
		log.Fatal(err)
	}
	return &Platform{
		Pos:   *NewPos(x, y),
		width: width,
		image: img,
	}
}

func getPlatformIndex(width float64) int {
	return (int(width) - 75) / 25 // 5 different widths from 75 to 175
}

func GenerateNewRandomPlatform(prevPlatform *Platform, score int) *Platform {
	x := prevPlatform.x
	minX := x + prevPlatform.width + platformSpacing - 50
	maxX := x + prevPlatform.width + platformSpacing + 50
	randX := float64(rand.Intn(int(maxX)-int(minX))) + minX

	y := prevPlatform.y
	minY := max(y-maxYDeltaTop, maxPlatformHeight)
	maxY := float64(screenHeight - minimumPlatformHeight)
	randY := float64(rand.Intn(int(maxY)-int(minY))) + minY

	randWidth := pickWidth(score, 175, 150, 125, 100, 75) // these numbers correspond to the widths of the building assets
	return NewPlatform(randX, randY, randWidth)
}

func pickWidth(counter int, numbers ...float64) float64 {
	var weights []float64
	// As counter increases, shift probability towards later numbers
	if counter > 60 {
		weights = []float64{0.01, 0.02, 0.3, 0.3, 0.37}
	} else if counter > 50 {
		weights = []float64{0.05, 0.05, 0.3, 0.3, 0.3}
	} else if counter > 40 {
		weights = []float64{0.15, 0.15, 0.25, 0.25, 0.2}
	} else if counter > 30 {
		weights = []float64{0.3, 0.2, 0.2, 0.15, 0.15}
	} else if counter > 20 {
		weights = []float64{0.5, 0.3, 0.1, 0.05, 0.05}
	} else if counter > 10 {
		weights = []float64{0.37, 0.3, 0.3, 0.02, 0.01}
	} else {
		weights = []float64{0.5, 0.4, 0.1, 0, 0}
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
	return numbers[len(numbers)-1] // Fallback should never happen
}

func giveOrTake(num, delta float64) float64 {
	return float64(rand.Intn(int(num+delta)-int(num-delta))) + num - delta
}

func (p *Platform) Draw(screen *ebiten.Image, cameraX float64) {
	x := p.x - cameraX
	scaleX, scaleY := 1.0, 1.0

	imgWidth := p.image.Bounds().Dx()
	imgHeight := p.image.Bounds().Dy()

	const revealFrames = 10
	const minAlpha = 0.9
	const maxAlpha = 1.0
	const buildingFeature = true

	if p.framesSinceVisited < revealFrames && buildingFeature {
		// unvisited image
		unvisitedCoor := &ebiten.DrawImageOptions{}
		unvisitedCoor.GeoM.Scale(scaleX, scaleY)
		unvisitedCoor.GeoM.Translate(x, p.y)

		screen.DrawImage(p.image, unvisitedCoor)
		if !p.visited {
			return
		}
		// progressively reveal image
		progress := float64(p.framesSinceVisited) / float64(revealFrames)

		for y := 0; y < imgHeight; y++ {
			sliceProgress := float64(y) / float64(imgHeight)
			alphaProgress := progress - sliceProgress

			if alphaProgress <= 0 {
				continue // not revealed yet
			}
			if alphaProgress > 1 {
				alphaProgress = 1
			}

			alpha := minAlpha + alphaProgress*(maxAlpha-minAlpha)
			slice := p.visitedImage.SubImage(image.Rect(0, y, imgWidth, y+1)).(*ebiten.Image)

			revealCoor := &ebiten.DrawImageOptions{}
			revealCoor.GeoM.Translate(x, p.y+float64(y))
			revealCoor.ColorScale.Scale(1, 1, 1, float32(alpha))

			screen.DrawImage(slice, revealCoor)
		}
	} else {
		visitedCoor := &ebiten.DrawImageOptions{}
		visitedCoor.GeoM.Scale(scaleX, scaleY)
		visitedCoor.GeoM.Translate(x, p.y)
		screen.DrawImage(p.visitedImage, visitedCoor)
	}
}

func (p *Platform) drawHitBox(screen *ebiten.Image, cameraX float64) {
	vector.DrawFilledRect(screen, float32(p.x-cameraX), float32(p.y), float32(p.width), float32(maxPlatformHeight), color.RGBA{R: 255}, false)
}

func (p *Platform) Visit() {
	p.visited = true
	newImage, err := media.Instance.LoadVisitedBuildingImage(getPlatformIndex(p.width))
	if err != nil {
		log.Fatal(err)
	}
	p.visitedImage = newImage
}
