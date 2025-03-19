package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Platform struct {
	x, y, width, height float64
}

type Game struct {
	playerX   float64
	playerY   float64
	velocityY float64
	isJumping bool
	platforms []Platform
	cameraX   float64
}

const (
	screenWidth     = 640
	screenHeight    = 480
	gravity         = 0.5
	jumpForce       = -12
	playerSpeed     = 5
	playerSize      = 20
	platformWidth   = 200
	platformHeight  = 10
	platformSpacing = 300
)

func (g *Game) initPlatforms() {
	rand.Seed(time.Now().UnixNano())
	maxYDelta := 150.0
	g.platforms = []Platform{{x: 200, y: screenHeight - 300, width: platformWidth, height: platformHeight}}
	for i := 1; i < 40; i++ {
		x := g.platforms[i-1].x
		y := g.platforms[i-1].y
		minY := max(y-maxYDelta, playerSize+20)
		maxY := int(min(screenHeight, y+maxYDelta))
		randY := float64(rand.Intn(maxY-int(minY))) + minY
		g.platforms = append(g.platforms, Platform{
			y:      randY,
			x:      x + platformSpacing,
			width:  platformWidth,
			height: platformHeight,
		})
	}
}

func (g *Game) Update() error {
	// Move left and right
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerX -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerX += playerSpeed
	}

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
		if g.playerY+playerSize >= p.y &&
			g.playerX+playerSize > p.x && g.playerX < p.x+p.width {
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
		ebitenutil.DrawRect(screen, p.x-g.cameraX, p.y, p.width, p.height, color.RGBA{128, 128, 128, 255})
	}
	ebitenutil.DrawRect(screen, g.playerX-g.cameraX, g.playerY, playerSize, playerSize, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{playerX: screenWidth / 2, playerY: screenHeight / 2}
	game.initPlatforms()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
