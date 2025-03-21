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
	imagesFilePath      = "assets/images/"
	playerImageFileName = "guy"
	imageFileExtension  = ".png"
	NumPlayerImages     = 8
)

type Manager struct {
	playerImages map[string]*ebiten.Image
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
		playerImages: make(map[string]*ebiten.Image),
	}
	for i := range NumPlayerImages {
		fileName := buildPlayerImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		result.playerImages[fileName] = image
	}
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

func buildPlayerImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", playerImageFileName, i, imageFileExtension)
}
