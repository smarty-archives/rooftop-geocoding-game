package game

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth            = 640
	screenHeight           = 480
	playerSize             = 40
	platformWidth          = 200
	platformSpacing        = 300
	maxYDeltaTop           = 120
	startingPlatformHeight = 300
	jumpApexHeight         = 140 // todo calculate from jumpForce
	startingPlatformX      = (screenWidth / 2) - (platformWidth / 2)
	startingPlatformY      = screenHeight - startingPlatformHeight
	lightGravity           = 0.5
	gravity                = 0.6
	heavyGravity           = 0.7
	framesWithLightGravity = 20
)

var (
	colorGray       = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	colorSmartyBlue = color.RGBA{R: 0, G: 102, B: 255, A: 255}
	colorText       = color.Black
	bot             = false
)

type Game struct {
	player    *Player
	platforms []*Platform
	cameraX   float64
	score     int // Score counter
	gameOver  bool
	font      font.Face
}

func NewGame() *Game {
	g := &Game{}
	g.initPlatforms()
	g.player = NewPlayer()

	// Use the default basic font from Ebiten
	g.font = basicfont.Face7x13

	return g
}

func (g *Game) Update() error {
	g.handlePlatforms()
	g.debug()

	// Player fell too low
	if g.player.y >= screenHeight*2 {
		g.gameOver = true
	}

	// If game over, reset the game when enter key is pressed
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			if g.player.x < g.getFirstPlatform().GetX() {
				bot = true
			}
			g.startOver()
		}
	} else {
		prevX := g.player.x // Store previous X position for side collision correction
		g.handlePlayer()
		g.applyGravity()
		g.handlePlatformCollision(prevX)
		g.handleScreenBounds()
		g.handleCameraMovement()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
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
		text.Draw(screen, gameOverText, g.font, gameOverX, 60, colorText)
		text.Draw(screen, restartText, g.font, restartX, 90, colorText)
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
		text.Draw(screen, botText, g.font, botX, 60, colorText)
		text.Draw(screen, botText2, g.font, botX2, 80, colorText)
	}

	// Draw the platforms
	for _, p := range g.platforms {
		p.Draw(screen, g.cameraX)
	}
	// todo move some stuff to init
	// Draw player
	g.player.Draw(screen, g.cameraX)
	//ebitenutil.DrawRect(screen, g.player.x-g.cameraX, g.player.y, playerSize, playerSize, color.White)

	// Draw score at top left
	text.Draw(screen, "Rooftops Geocoded: "+strconv.Itoa(g.score), g.font, 10, 20, colorText)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) initPlatforms() {
	g.platforms = []*Platform{NewPlatform(startingPlatformX, startingPlatformY)}
	for i := 1; i < 2; i++ {
		g.platforms = append(g.platforms, GenerateNewRandomPlatform(g.platforms[i-1]))
	}
}

func (g *Game) handlePlatforms() {
	// Generate New
	if g.distToLastPlatform() < platformWidth {
		g.platforms = append(g.platforms, GenerateNewRandomPlatform(g.getLastPlatform()))
	}
	// Cleanup
	if g.distToFirstPlatform() > platformSpacing*2 {
		g.platforms = g.platforms[1:]
	}
}

func (g *Game) distToLastPlatform() float64 {
	lastPlatform := g.getLastPlatform()
	return g.player.GetDist(lastPlatform)
}

func (g *Game) distToFirstPlatform() float64 {
	firstPlatform := g.getFirstPlatform()
	return g.player.GetDist(firstPlatform)
}

func (g *Game) getLastPlatform() *Platform {
	return g.platforms[len(g.platforms)-1]
}

func (g *Game) getFirstPlatform() *Platform {
	return g.platforms[0]
}

// WARNING: getFirstPlatform will panic if you don't initialize the platforms after this
func (g *Game) resetGameState() {
	g.player.Reset()
	g.platforms = g.platforms[:0]
	g.cameraX = 0
	g.score = 0
	g.gameOver = false
}

func (g *Game) startOver() {
	g.resetGameState()
	g.player.ResetPlayer()
	g.initPlatforms()
}

func (g *Game) playerOnPlatform(p Platform) bool {
	// Check if player is within platform's horizontal range
	playerRight := g.player.x + playerSize
	playerLeft := g.player.x
	platformRight := p.x + p.width
	platformLeft := p.x

	// **Vertical collision (Landing on platform)**
	return playerRight > platformLeft && playerLeft < platformRight && // Player overlaps horizontally
		g.player.y+playerSize >= p.y && g.player.y+playerSize-g.player.velocityY <= p.y // Player is falling onto the platform
}

func (g *Game) playerInAir() bool {
	for _, p := range g.platforms {
		if g.playerOnPlatform(*p) {
			return false
		}
	}
	return true
}

func (g *Game) botLogic() {
	g.player.AccelerateRight()
	for i, p := range g.platforms {
		if i < len(g.platforms)-1 {
			if g.botShouldJump(*p, *g.platforms[i+1]) {
				g.player.Jump()
			}
		}

	}
}

func (g *Game) playerControls() {
	// Accelerate left and right
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.AccelerateLeft()
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.AccelerateRight()
	} else {
		g.slowPlayer()
	}

	// Jumping logic
	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.player.isJumping {
		if g.playerCanJump() {
			g.player.Jump()
		}
	}
}

func (g *Game) botShouldJump(p, nextP Platform) bool {
	if p.y+50 < nextP.y {
		return false
	}
	pRight := p.x + p.width
	return pRight-50 < g.player.x && g.player.x < pRight && !g.player.isJumping
}

func (g *Game) slowPlayer() {
	g.player.velocityX *= .8
}

func (g *Game) playerCanJump() bool {
	return !g.playerInAir()
}

func (g *Game) handlePlatformCollision(prevX float64) {
	for _, p := range g.platforms {
		// **Vertical collision (Landing on platform)**
		if g.playerOnPlatform(*p) {
			// Land on the platform
			g.player.y = p.y - playerSize
			g.player.velocityY = 0
			g.player.isJumping = false

			if !p.visited {
				p.visited = true
				g.score++
			}
		}

		playerRight := g.player.x + playerSize
		playerLeft := g.player.x
		platformRight := p.x + p.width
		platformLeft := p.x
		// **Side collision (Hitting the sides of the platform)**
		if g.player.y+playerSize > p.y { // Player is below platform surface
			if prevX+playerSize <= platformLeft && playerRight > platformLeft { // Hitting left side
				g.player.x = platformLeft - playerSize
				g.player.velocityX = 0
			} else if prevX >= platformRight && playerLeft < platformRight { // Hitting right side
				g.player.x = platformRight
				g.player.velocityX = 0
			}
		}
	}
}

func (g *Game) handleScreenBounds() {
	if g.player.x < g.cameraX {
		g.player.x = g.cameraX
	} else if g.player.x+playerSize > g.cameraX+screenWidth {
		g.player.x = g.cameraX + screenWidth - playerSize
	}
}

func (g *Game) handleCameraMovement() {
	g.cameraX = max(g.player.x-screenWidth/2+playerSize/2, g.getFirstPlatform().x-playerSize)
	if g.cameraX < 0 {
		g.cameraX = 0
	}
}

func (g *Game) applyGravity() {
	if g.player.isJumping {
		g.player.framesSinceJump++
	} else {
		g.player.framesSinceJump = 0
	}
	currentGravity := gravity
	if g.player.velocityY > 0 {
		currentGravity = heavyGravity
	}
	if g.player.framesSinceJump < framesWithLightGravity {
		currentGravity = lightGravity
	}
	g.player.velocityY += currentGravity
	g.player.y += g.player.velocityY
}

func (g *Game) handlePlayer() {
	if bot {
		g.botLogic()
	} else {
		g.playerControls()
	}
}

func (g *Game) debug() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		//fmt.Printf("Total Platforms: %d\n", len(g.platforms))
		fmt.Printf("Frames Since Jumping: %d\n", g.player.framesSinceJump)
	}
}
