package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Platform struct {
	x, y, width, height float64
}

func (p *Platform) Draw(screen *ebiten.Image, cameraX float64) {
	ebitenutil.DrawRect(screen, p.x-cameraX, p.y, p.width, p.height, colorGray)
}

type Game struct {
	playerImage  *ebiten.Image
	playerX      float64
	playerY      float64
	playerXSpeed float64
	velocityY    float64
	isJumping    bool
	platforms    []Platform
	cameraX      float64
}

const (
	screenWidth     = 640
	screenHeight    = 480
	gravity         = 0.5
	jumpForce       = -12
	jumpApexHeight  = 140 // todo calculate from jumpForce
	playerSpeed     = .2
	playerSize      = 40
	platformWidth   = 200
	platformHeight  = 10
	platformSpacing = 300
)

var (
	colorGray = color.RGBA{R: 128, G: 128, B: 128, A: 255}
)

// todo offload previous platforms and add more platforms instead of initializing all at once

func (g *Game) initPlatforms() {
	rand.Seed(time.Now().UnixNano())
	maxYDeltaTop := 120.0
	maxYDeltaBottom := 160.0
	startingPlatformHeight := 300.0
	startingPlatformX := 200.0
	startingPlatformY := screenHeight - startingPlatformHeight
	g.platforms = []Platform{{x: startingPlatformX, y: startingPlatformY, width: platformWidth, height: startingPlatformHeight}}
	for i := 1; i < 4000; i++ {
		x := g.platforms[i-1].x
		y := g.platforms[i-1].y
		minY := max(y-maxYDeltaTop, playerSize+jumpApexHeight)
		maxY := int(min(screenHeight-20, y+maxYDeltaBottom))
		randY := float64(rand.Intn(maxY-int(minY))) + minY
		g.platforms = append(g.platforms, Platform{
			y:      randY,
			x:      x + platformSpacing,
			width:  platformWidth,
			height: screenHeight - randY,
		})
	}
}

func (g *Game) Update() error {
	// Move left and right
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if g.playerXSpeed > 0 {
			g.playerXSpeed = 0
		}
		if g.playerXSpeed >= -5 {
			g.playerXSpeed -= playerSpeed
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if g.playerXSpeed < 0 {
			g.playerXSpeed = 0
		}
		if g.playerXSpeed <= 5 {
			g.playerXSpeed += playerSpeed
		}
	} else {
		g.playerXSpeed *= .8
	}
	g.playerX += g.playerXSpeed

	// Jumping logic
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.isJumping {
		g.velocityY = jumpForce
		g.isJumping = true
	}

	// Apply gravity
	g.velocityY += gravity
	g.playerY += g.velocityY

	// Collision with platforms
	for i := range g.platforms {
		p := &g.platforms[i]
		playerOnOrLowerThanPlatform := g.playerY+playerSize >= p.y
		playerInXRangeOfPlatform := g.playerX+playerSize > p.x && g.playerX < p.x+p.width
		if playerOnOrLowerThanPlatform && playerInXRangeOfPlatform {
			g.playerY = p.y - playerSize
			g.velocityY = 0
			g.isJumping = false
		}
	}

	// Keep player within screen bounds
	if g.playerX < g.cameraX {
		g.playerX = g.cameraX
	} else if g.playerX+playerSize > g.cameraX+screenWidth {
		g.playerX = g.cameraX + screenWidth - playerSize
	}

	// Move camera horizontally as player moves
	g.cameraX = g.playerX - screenWidth/2 + playerSize/2
	if g.cameraX < 0 {
		g.cameraX = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.platforms {
		p.Draw(screen, g.cameraX)
	}
	playerCoor := &ebiten.DrawImageOptions{}
	scaleX := playerSize / float64(g.playerImage.Bounds().Dx())
	scaleY := playerSize / float64(g.playerImage.Bounds().Dy())
	playerCoor.GeoM.Scale(scaleX, scaleY)
	playerCoor.GeoM.Translate(g.playerX-g.cameraX, g.playerY)
	screen.DrawImage(g.playerImage, playerCoor)
	//ebitenutil.DrawRect(screen, g.playerX-g.cameraX, g.playerY, playerSize, playerSize, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	image, _, err := ebitenutil.NewImageFromFile("assets/guy.png")
	if err != nil {
		log.Fatal(err)
	}
	game := &Game{playerX: screenWidth / 2, playerY: screenHeight / 2, playerImage: image}
	game.initPlatforms()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
