package game

import (
	"image/color"
	"math"
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
	coinSize               = 20
	platformSpacing        = 100
	maxYDeltaTop           = 120
	startingPlatformHeight = 300
	jumpApexHeight         = 140 // todo calculate from jumpForce
	minimumPlatformHeight  = 20
	startingPlatformWidth  = 200
	startingPlatformX      = (screenWidth / 2) - (startingPlatformWidth / 2)
	startingPlatformY      = screenHeight - startingPlatformHeight
	lightGravity           = 0.4
	gravity                = 0.7
	heavyGravity           = 0.8
)

var (
	colorGray            = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	colorSmartyBlue      = color.RGBA{R: 0, G: 102, B: 255, A: 255}
	colorText            = color.White
	bot                  = false
	botFramesLeftJumping = 0
)

type Game struct {
	player           *Player
	platforms        []*Platform
	cameraX          float64
	score            int // Score counter
	gameOver         bool
	font             font.Face
	backgroundLayers []Layer
	coin             *Coin
}

func NewGame() *Game {
	g := &Game{}
	g.initPlatforms()
	g.player = NewPlayer()
	g.font = basicfont.Face7x13 // Use the default basic font from Ebiten
	g.backgroundLayers = NewLayers()
	g.coin = NewCoin()
	return g
}

func (g *Game) Update() error {
	g.handleBackgroundLayers()
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
		prevX := g.player.x // Store previous x position for side collision correction
		g.handlePlayer()
		// g.applyGravity() // todo make bot and player implement interface that applyGravity can use instead of checking for spacebar
		g.handlePlatformCollision(prevX)
		g.handleScreenBounds()
		g.handleCameraMovement()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackgroundLayers(screen)
	//screen.Fill(color.White)
	// If the game is over, display the "Game Over" screen
	if g.gameOver {
		// Text to display
		gameOverText := "That's not a rooftop geocode!"
		restartText := "Press Enter to Restart"

		// Measure the width of the "Game Over" text
		gameOverWidth := font.MeasureString(g.font, gameOverText).Ceil()
		restartWidth := font.MeasureString(g.font, restartText).Ceil()

		// Calculate the x position to center the "Game Over" text
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

		// Calculate the x position to center the "Game Over" text
		botX := (screenWidth - botWidth) / 2
		botX2 := (screenWidth - botWidth2) / 2

		// Draw the "Game Over" text and restart message
		text.Draw(screen, botText, g.font, botX, 60, colorText)
		text.Draw(screen, botText2, g.font, botX2, 80, colorText)
	}

	// Draw the platforms & coins
	for _, p := range g.platforms {
		p.Draw(screen, g.cameraX)
		if !p.visited {
			g.drawCoin(screen, p)
		}
	}
	// Draw player
	g.player.Draw(screen, g.cameraX)

	// Draw score at top left
	text.Draw(screen, "Rooftops Geocoded: "+strconv.Itoa(g.score), g.font, 10, 20, colorText)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) initPlatforms() {
	g.platforms = []*Platform{NewPlatform(startingPlatformX, startingPlatformY, startingPlatformWidth)}
	for i := 1; i < 2; i++ {
		g.platforms = append(g.platforms, GenerateNewRandomPlatform(g.platforms[i-1], g.score))
	}
}

func (g *Game) handlePlatforms() {
	// Generate New
	if g.distToLastPlatform() < screenWidth/2 {
		g.platforms = append(g.platforms, GenerateNewRandomPlatform(g.getLastPlatform(), g.score))
	}
	// Cleanup
	if g.distToFirstPlatform() > screenWidth {
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
	if g.playerCanJump() && g.botShouldJump() {
		g.player.Jump()
	}
}

// heightAfterXFramesOfJumping assumes that player is moving at top speed
func (g *Game) heightAfterXFramesOfJumping(jumpFrames, totalFrames int) (y float64, velocityY float64) {
	finalY := g.player.y
	velocity := g.player.GetJumpForce()
	for i := range totalFrames {
		if velocity > 0 {
			velocity += heavyGravity
		} else if i < jumpFrames {
			velocity += lightGravity
		} else {
			velocity += gravity
		}
		finalY += velocity
	}
	return finalY, velocity
}

func (g *Game) NextUnvisitedPlatform() *Platform {
	for _, p := range g.platforms {
		if !p.visited {
			return p
		}
	}
	return nil
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

func (g *Game) botShouldJump() bool {
	platformPos := g.NextUnvisitedPlatform()
	numFrames := (platformPos.x - g.player.x) / g.player.maxPlayerSpeed
	for i := range 30 {
		newY, newVelocityY := g.heightAfterXFramesOfJumping(i, int(numFrames)+1)
		if newY < platformPos.y && newVelocityY > 0 {
			if i > 1 && g.playerHasMoreRunway() {
				return false
			}
			botFramesLeftJumping = i
			return true
		}
	}
	return false
}

func (g *Game) playerHasMoreRunway() bool {
	for i, platform := range g.platforms {
		if g.playerOnPlatform(*platform) {
			if g.platforms[i+1].y+80 <= platform.y { // this assumes there is always a platform after the one the player is on
				return false
			}
			return platform.x+platform.width > g.player.x+50
		}
	}
	return false
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
	// todo make this use a function that can be called by applyBotGravity and heightAfterXFramesOfJumping
	currentGravity := gravity
	if g.player.velocityY > 0 {
		currentGravity = heavyGravity
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		currentGravity = lightGravity
	}
	g.player.velocityY += currentGravity
	g.player.y += g.player.velocityY
}

func (g *Game) applyBotGravity() {
	currentGravity := gravity
	if g.player.velocityY > 0 {
		currentGravity = heavyGravity
	} else if botFramesLeftJumping > 0 {
		currentGravity = lightGravity
		botFramesLeftJumping--
	}
	g.player.velocityY += currentGravity
	g.player.y += g.player.velocityY
}

func (g *Game) handlePlayer() {
	if bot {
		g.botLogic()
		g.applyBotGravity()
	} else {
		g.playerControls()
		g.applyGravity()
	}
}

func (g *Game) debug() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		//fmt.Printf("Total Platforms: %d\n", len(g.platforms))
	}
}

func (g *Game) drawBackgroundLayers(screen *ebiten.Image) {
	for _, layer := range g.backgroundLayers {
		imgWidth := layer.Image.Bounds().Dx()
		imgHeight := layer.Image.Bounds().Dy()

		// Scale the image to match the screen height
		scaleY := float64(screenHeight) / float64(imgHeight)
		scaleX := scaleY // Maintain aspect ratio horizontally

		// Calculate the total width of a single scaled image
		scaledWidth := float64(imgWidth) * scaleX

		// Ensure the offset wraps around seamlessly
		xOffset := layer.OffsetX
		xOffset = math.Mod(xOffset, scaledWidth)
		if xOffset > 0 {
			xOffset -= scaledWidth
		}

		// Draw enough images to cover the entire screen width
		for x := xOffset; x < screenWidth; x += scaledWidth {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scaleX, scaleY)
			op.GeoM.Translate(x, 0)
			screen.DrawImage(layer.Image, op)
		}
	}
}

func (g *Game) handleBackgroundLayers() {
	for i := range g.backgroundLayers {
		g.backgroundLayers[i].OffsetX = -g.cameraX * g.backgroundLayers[i].Speed
	}
}

func (g *Game) drawCoin(screen *ebiten.Image, p *Platform) {
	op := &ebiten.DrawImageOptions{}
	scaleX := coinSize / float64(g.coin.Image.Bounds().Dx())
	scaleY := coinSize / float64(g.coin.Image.Bounds().Dy())
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(p.x+p.width/2-g.cameraX-coinSize/2, p.y-coinSize*1.5)
	screen.DrawImage(g.coin.Image, op)
}
