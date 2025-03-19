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
	jumpApexHeight  = 140 // todo calculate from jumpForce
	playerSpeed     = 5
	playerSize      = 20
	platformWidth   = 200
	platformHeight  = 10
	platformSpacing = 300
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
