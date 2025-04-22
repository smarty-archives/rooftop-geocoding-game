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
	imagesFilePath               = "assets/images/"
	runningImageFileName         = "guy"
	idleImageFileName            = "idle"
	backgroundImageFileName      = "layer"
	buildingImageFileName        = "building"
	visitedBuildingImageFileName = "visited-building"
	shareBtnImageFileName        = "share-btn"
	imageFileExtension           = ".png"
	NumPlayerImages              = 8
	NumIdleImages                = 2
	NumBackgroundLayers          = 1 // Warning: if this number is more than the number of speeds in NewLayers, then it will panic
	NumBuildingImages            = 5
)

type Manager struct {
	playerImages          map[string]*ebiten.Image
	idleImages            map[string]*ebiten.Image
	backgroundImages      map[string]*ebiten.Image
	buildingImages        map[string]*ebiten.Image
	visitedBuildingImages map[string]*ebiten.Image
	cloudImage            *ebiten.Image
	titleImage            *ebiten.Image
	mutedImage            *ebiten.Image
	playingImage          *ebiten.Image
	shareButtonImage      *ebiten.Image
}

var (
	Instance         = &Manager{}
	errImageNotFound = errors.New("image not found")
)

func init() {
	Instance = NewManager()
}

func NewManager() *Manager {
	result := &Manager{}
	result.initializeRunningImages()
	result.initializeIdleImages()
	result.initializeBackgroundImages()
	result.initializeBuildingImages()
	result.initializeVisitedBuildingImages()
	result.initializeCloudImage()
	result.initializeTitleImage()
	result.initializeMutedImage()
	result.initializePlayingImage()
	result.initializeShareButtonImage()
	return result
}

func (m *Manager) LoadRunningImage(i int) (*ebiten.Image, error) {
	fileName := buildRunningImageFileName(i)
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

func (m *Manager) LoadBuildingImage(i int) (*ebiten.Image, error) {
	fileName := buildBuildingImageFileName(i)
	image, ok := m.buildingImages[fileName]
	if !ok {
		return nil, errImageNotFound
	}
	return image, nil
}

func (m *Manager) LoadVisitedBuildingImage(i int) (*ebiten.Image, error) {
	fileName := buildVisitedBuildingImageFileName(i)
	image, ok := m.visitedBuildingImages[fileName]
	if !ok {
		return nil, errImageNotFound
	}
	return image, nil
}

func (m *Manager) GetCloudImage() *ebiten.Image {
	return m.cloudImage
}

func (m *Manager) GetTitleImage() *ebiten.Image {
	return m.titleImage
}

func (m *Manager) GetMutedImage() *ebiten.Image {
	return m.mutedImage
}

func (m *Manager) GetPlayingImage() *ebiten.Image {
	return m.playingImage
}

func (m *Manager) GetShareButtonImage() *ebiten.Image {
	return m.shareButtonImage
}

func (m *Manager) initializeRunningImages() {
	m.playerImages = make(map[string]*ebiten.Image)
	for i := range NumPlayerImages {
		fileName := buildRunningImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		m.playerImages[fileName] = image
	}
}

func (m *Manager) initializeIdleImages() {
	m.idleImages = make(map[string]*ebiten.Image)
	for i := range NumIdleImages {
		fileName := buildIdleImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		m.idleImages[fileName] = image
	}
}

func (m *Manager) initializeBackgroundImages() {
	m.backgroundImages = make(map[string]*ebiten.Image)
	for i := range NumBackgroundLayers {
		fileName := buildBackgroundImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		m.backgroundImages[fileName] = image
	}
}

func (m *Manager) initializeBuildingImages() {
	m.buildingImages = make(map[string]*ebiten.Image)
	for i := range NumBuildingImages {
		fileName := buildBuildingImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		m.buildingImages[fileName] = image
	}
}

func (m *Manager) initializeVisitedBuildingImages() {
	m.visitedBuildingImages = make(map[string]*ebiten.Image)
	for i := range NumBuildingImages {
		fileName := buildVisitedBuildingImageFileName(i)
		image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
		if err != nil {
			log.Fatal(err)
		}
		m.visitedBuildingImages[fileName] = image
	}
}

func (m *Manager) initializeCloudImage() {
	fileName := "cloud.png"
	image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	m.cloudImage = image
}

func (m *Manager) initializeTitleImage() {
	fileName := "title.png"
	image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	m.titleImage = image
}

func (m *Manager) initializeMutedImage() {
	fileName := "muted.png"
	image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	m.mutedImage = image
}

func (m *Manager) initializePlayingImage() {
	fileName := "playing.png"
	image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	m.playingImage = image
}

func (m *Manager) initializeShareButtonImage() {
	fileName := "share-btn.png"
	image, _, err := ebitenutil.NewImageFromFile(filepath.Join(imagesFilePath, fileName))
	if err != nil {
		log.Fatal(err)
	}
	m.shareButtonImage = image
}

func buildRunningImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", runningImageFileName, i, imageFileExtension)
}

func buildIdleImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", idleImageFileName, i, imageFileExtension)
}

func buildBackgroundImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", backgroundImageFileName, i, imageFileExtension)
}

func buildBuildingImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", buildingImageFileName, i, imageFileExtension)
}

func buildVisitedBuildingImageFileName(i int) string {
	return fmt.Sprintf("%s%d%s", visitedBuildingImageFileName, i, imageFileExtension)
}
