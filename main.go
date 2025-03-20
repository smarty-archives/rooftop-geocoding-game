package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
	playerImage *ebiten.Image
	player      Player
	platforms   []*Platform
	cameraX     float64
	score       int // Score counter
	gameOver    bool
	font        font.Face
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
	startingPlatformX      = (screenWidth / 2) - (platformWidth / 2)
	startingPlatformY      = screenHeight - startingPlatformHeight
)

var (
	colorGray       = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	colorSmartyBlue = color.RGBA{R: 0, G: 102, B: 255, A: 255}
	bot             = false
)

// todo offload previous platforms and add more platforms instead of initializing all at once

func (g *Game) initPlayer() {
	g.player.x = screenWidth/2 - (playerSize / 2)
	g.player.y = 0 // startingPlatformY - playerSize*2
}

const (
	maxYDeltaTop    = 120.0
	maxYDeltaBottom = 160.0
)

func (g *Game) initPlatforms() {
	g.platforms = []*Platform{NewPlatform(startingPlatformX, startingPlatformY)}
	for i := 1; i < 2; i++ {
		g.platforms = append(g.platforms, GenerateNewRandomPlatform(g.platforms[i-1]))
	}
}

func (g *Game) PlatformHandling() {
	// Generate New
	if g.DistToLastPlatform() < platformWidth {
		g.platforms = append(g.platforms, GenerateNewRandomPlatform(g.GetLastPlatform()))
	}
	// Cleanup
	if g.DistToFirstPlatform() > platformSpacing*2 {
		g.platforms = g.platforms[1:]
	}
}

func (g *Game) DistToLastPlatform() float64 {
	lastPlatform := g.GetLastPlatform()
	return g.player.GetDist(lastPlatform)
}

func (g *Game) DistToFirstPlatform() float64 {
	firstPlatform := g.GetFirstPlatform()
	return g.player.GetDist(firstPlatform)
}

func (g *Game) GetLastPlatform() *Platform {
	return g.platforms[len(g.platforms)-1]
}

func (g *Game) GetFirstPlatform() *Platform {
	return g.platforms[0]
}

func (g *Game) init() {
	g.initPlatforms()
	g.initPlayer()

	// Use the default basic font from Ebiten
	g.font = basicfont.Face7x13
}

func (g *Game) resetGameState() {
	g.player.x = 0
	g.player.y = 0
	g.player.xSpeed = 0
	g.player.velocityY = 0
	g.player.isJumping = false
	// getFirstPlatform will panic if you don't initialize the platforms after this
	g.platforms = g.platforms[:0]
	g.cameraX = 0
	g.score = 0
	g.gameOver = false
}

func (g *Game) startOver() {
	g.resetGameState()
	g.init()
	bot = true
}

func (g *Game) Update() error {
	g.PlatformHandling()
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		fmt.Printf("Total Platforms: %d\n", len(g.platforms))
	}
	// Player fell too low
	if g.player.y >= screenHeight*2 {
		g.gameOver = true
	}

	// If game over, reset the game when a key is pressed
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.startOver()
		}
		return nil
	}

	// Store previous X position for side collision correction
	prevX := g.player.x

	if bot {
		g.movePlayerRight()
		for _, p := range g.platforms {
			if g.botShouldJump(*p) {
				g.jump()
			}

		}
	} else {
		// Move left and right
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			g.movePlayerLeft()
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			g.movePlayerRight()
		} else {
			g.slowPlayer()
		}

		// Jumping logic
		if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.player.isJumping {
			g.jump()
		}
	}

	// Apply gravity
	g.player.velocityY += gravity
	g.player.y += g.player.velocityY

	// Platform collision detection (both vertical and side)
	for _, p := range g.platforms {
		// Check if player is within platform's horizontal range
		playerRight := g.player.x + playerSize
		playerLeft := g.player.x
		platformRight := p.x + p.width
		platformLeft := p.x

		// **Vertical collision (Landing on platform)**
		if playerRight > platformLeft && playerLeft < platformRight && // Player overlaps horizontally
			g.player.y+playerSize > p.y && g.player.y+playerSize-g.player.velocityY <= p.y { // Player is falling onto the platform
			// Land on the platform
			g.player.y = p.y - playerSize
			g.player.velocityY = 0
			g.player.isJumping = false

			if !p.visited {
				p.visited = true
				g.score++
			}
		}

		// **Side collision (Hitting the sides of the platform)**
		if g.player.y+playerSize > p.y { // Player is below platform surface
			if prevX+playerSize <= platformLeft && playerRight > platformLeft { // Hitting left side
				g.player.x = platformLeft - playerSize
				g.player.xSpeed = 0
			} else if prevX >= platformRight && playerLeft < platformRight { // Hitting right side
				g.player.x = platformRight
				g.player.xSpeed = 0
			}
		}
	}

	// todo seems unnecessary now that camera follows the guy
	// Keep player within screen bounds
	if g.player.x < g.cameraX {
		g.player.x = g.cameraX
	} else if g.player.x+playerSize > g.cameraX+screenWidth {
		g.player.x = g.cameraX + screenWidth - playerSize
	}

	// Move camera horizontally as player moves
	g.cameraX = max(g.player.x-screenWidth/2+playerSize/2, g.GetFirstPlatform().x-playerSize)
	if g.cameraX < 0 {
		g.cameraX = 0
	}

	return nil
}

func (g *Game) botShouldJump(p Platform) bool {
	pRight := p.x + p.width
	return pRight-50 < g.player.x && g.player.x < pRight && !g.player.isJumping
}

func (g *Game) Draw(screen *ebiten.Image) {
	// If the game is over, display the "Game Over" screen
	if g.gameOver {
		// Text to display
		gameOverText := "That's not a rooftop geocode!"
		restartText := "Press Enter to Restart"

		// Measure the width of the "Game Over" text
		gameOverWidth := font.MeasureString(g.font, gameOverText).Ceil()
		restartWidth := font.MeasureString(g.font, restartText).Ceil()

		// Calculate the X position to center the "Game Over" text
		gameOverX := (screenWidth - gameOverWidth) / 2
		restartX := (screenWidth - restartWidth) / 2

		// Draw the "Game Over" text and restart message
		text.Draw(screen, gameOverText, g.font, gameOverX, 60, color.White)
		text.Draw(screen, restartText, g.font, restartX, 90, color.White)
	}
	if bot {
		botText := "You can't be trusted to do it yourself."
		botText2 := "Now you have to use Smarty."

		// Measure the width of the "Game Over" text
		botWidth := font.MeasureString(g.font, botText).Ceil()
		botWidth2 := font.MeasureString(g.font, botText2).Ceil()

		// Calculate the X position to center the "Game Over" text
		botX := (screenWidth - botWidth) / 2
		botX2 := (screenWidth - botWidth2) / 2

		// Draw the "Game Over" text and restart message
		text.Draw(screen, botText, g.font, botX, 60, color.White)
		text.Draw(screen, botText2, g.font, botX2, 80, color.White)
	}

	// Draw the platforms
	for _, p := range g.platforms {
		p.Draw(screen, g.cameraX)
	}
	// todo move some stuff to init
	// Draw player
	playerCoor := &ebiten.DrawImageOptions{}
	scaleX := playerSize / float64(g.playerImage.Bounds().Dx())
	scaleY := playerSize / float64(g.playerImage.Bounds().Dy())
	playerCoor.GeoM.Scale(scaleX, scaleY)
	playerCoor.GeoM.Translate(g.player.x-g.cameraX, g.player.y)
	screen.DrawImage(g.playerImage, playerCoor)
	//ebitenutil.DrawRect(screen, g.player.x-g.cameraX, g.player.y, playerSize, playerSize, color.White)

	// Draw score at top left
	text.Draw(screen, "Rooftops Geocoded: "+strconv.Itoa(g.score), g.font, 10, 20, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) movePlayerLeft() {
	if g.player.xSpeed > 0 {
		g.player.xSpeed = 0
	}
	if g.player.xSpeed >= -5 {
		g.player.xSpeed -= playerSpeed
	}
	g.player.x += g.player.xSpeed
}

func (g *Game) movePlayerRight() {
	if g.player.xSpeed < 0 {
		g.player.xSpeed = 0
	}
	if g.player.xSpeed <= 5 {
		g.player.xSpeed += playerSpeed
	}
	g.player.x += g.player.xSpeed
}

func (g *Game) slowPlayer() {
	g.player.xSpeed *= .8
}

func (g *Game) jump() {
	g.player.velocityY = jumpForce
	g.player.isJumping = true
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
