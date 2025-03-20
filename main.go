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
	screenWidth            = 640
	screenHeight           = 480
	gravity                = 0.5
	jumpForce              = -12
	jumpApexHeight         = 140 // todo calculate from jumpForce
	playerSpeed            = .2
	playerSize             = 40
	platformWidth          = 200
	platformSpacing        = 300
	startingPlatformHeight = 300.0
	startingPlatformX      = 200.0
	startingPlatformY      = screenHeight - startingPlatformHeight
)

var (
	colorGray = color.RGBA{R: 128, G: 128, B: 128, A: 255}
)

// todo offload previous platforms and add more platforms instead of initializing all at once

func (g *Game) initPlayer() {
	g.playerX = screenWidth / 2
	g.playerY = startingPlatformY - playerSize*2
}

func (g *Game) initPlatforms() {
	rand.Seed(time.Now().UnixNano())
	maxYDeltaTop := 120.0
	maxYDeltaBottom := 160.0
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

func (g *Game) init() {
	g.initPlatforms()
	g.initPlayer()
}

func (g *Game) resetGameState() {
	g.playerX = 0
	g.playerY = 0
	g.playerXSpeed = 0
	g.velocityY = 0
	g.isJumping = false
	g.platforms = g.platforms[:0]
	g.cameraX = 0
}

func (g *Game) Update() error {
	// Player fell too low
	if g.playerY >= screenHeight*2 {
		// Game Over
		g.resetGameState()
		g.init()
	}
	// Store previous X position for side collision correction
	prevX := g.playerX

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

	// Platform collision detection (both vertical and side)
	for _, p := range g.platforms {
		// Check if player is within platform's horizontal range
		playerRight := g.playerX + playerSize
		playerLeft := g.playerX
		platformRight := p.x + p.width
		platformLeft := p.x

		// **Vertical collision (Landing on platform)**
		if playerRight > platformLeft && playerLeft < platformRight && // Player overlaps horizontally
			g.playerY+playerSize > p.y && g.playerY+playerSize-g.velocityY <= p.y { // Player is falling onto the platform
			// Land on the platform
			g.playerY = p.y - playerSize
			g.velocityY = 0
			g.isJumping = false
		}

		// **Side collision (Hitting the sides of the platform)**
		if g.playerY+playerSize > p.y { // Player is below platform surface
			if prevX+playerSize <= platformLeft && playerRight > platformLeft { // Hitting left side
				g.playerX = platformLeft - playerSize
				g.playerXSpeed = 0
			} else if prevX >= platformRight && playerLeft < platformRight { // Hitting right side
				g.playerX = platformRight
				g.playerXSpeed = 0
			}
		}
	}

	// todo seems unnecessary now that camera follows the guy
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
	// todo move some stuff to init
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
	game := &Game{playerImage: image}
	game.init()
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
