package media

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	imagesFilePath          = "assets/images/"
	playerImageFileName     = "guy"
	idleImageFileName       = "idle"
	backgroundImageFileName = "layer"
	coinImageFileName       = "coin"
	imageFileExtension      = ".png"
	NumPlayerImages         = 8
	NumIdleImages           = 2
	NumBackgroundLayers     = 1 // Warning: if this number is more than the number of speeds in NewLayers, then it will panic
)

type Manager struct {
	playerImages     map[string]*ebiten.Image
	idleImages       map[string]*ebiten.Image
	backgroundImages map[string]*ebiten.Image
	coinImage        *ebiten.Image
}

var (
	Instance         = &Manager{}
	errImageNotFound = errors.New("image not found")
)

func init() {
	Instance = NewManager()
}

func NewManager() *Manager {
	result := &Manager{
		playerImages:     make(map[string]*ebiten.Image),
		idleImages:       make(map[string]*ebiten.Image),
		backgroundImages: make(map[string]*ebiten.Image),
	}
	// todo make function for these loops
	// Running Images
	for i := range NumPlayerImages {
		fileName := buildPlayerImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		result.playerImages[fileName] = image
	}
	// Idle Images
	for i := range NumIdleImages {
		fileName := buildIdleImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		result.idleImages[fileName] = image
	}
	// Background Images
	for i := range NumBackgroundLayers {
		fileName := buildBackgroundImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		result.backgroundImages[fileName] = image
	}
	// Coin Image
	fileName := buildCoinImageFileName()
	image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	result.coinImage = image
	return result
}

func (m *Manager) LoadPlayerImage(i int) (*ebiten.Image, error) {
	fileName := buildPlayerImageFileName(i)
	image, ok := m.playerImages[fileName]
	if !ok {
		return nil, errImageNotFound
	}
	return image, nil
}

func (m *Manager) LoadIdleImage(i int) (*ebiten.Image, error) {
	fileName := buildIdleImageFileName(i)
	image, ok := m.idleImages[fileName]
	if !ok {
		return nil, errImageNotFound
	}
	return image, nil
}

func (m *Manager) LoadBackgroundImage(i int) (*ebiten.Image, error) {
	fileName := buildBackgroundImageFileName(i)
	image, ok := m.backgroundImages[fileName]
	if !ok {
		return nil, errImageNotFound
	}
	return image, nil
}

func (m *Manager) LoadCoinImage() (*ebiten.Image, error) {
	return m.coinImage, nil
}

func buildPlayerImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", playerImageFileName, i, imageFileExtension)
}

func buildIdleImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", idleImageFileName, i, imageFileExtension)
}

func buildBackgroundImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", backgroundImageFileName, i, imageFileExtension)
}

func buildCoinImageFileName() string {
	return fmt.Sprintf("%s%s", coinImageFileName, imageFileExtension)
}
