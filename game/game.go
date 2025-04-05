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
	platformSpacing        = 100
	maxYDeltaTop           = 120
	minimumPlatformHeight  = 20
	maxPlatformHeight      = 285 // this is a little bit less than the height of the building assets
	maxPlatformWidth       = 175
	startingPlatformHeight = maxPlatformHeight
	startingPlatformWidth  = maxPlatformWidth
	startingPlatformX      = (screenWidth / 2) - (startingPlatformWidth / 2)
	startingPlatformY      = screenHeight - startingPlatformHeight
	lightGravity           = 0.4
	gravity                = 0.7
	heavyGravity           = 0.8
)

var (
	colorGray            = color.RGBA{R: 128, G: 128, B: 128, A: 255}
	colorSmartyBlue      = color.RGBA{R: 0, G: 102, B: 255, A: 255}
	colorText            = color.Black
	bot                  = false
	botFramesLeftJumping = 0
)

type Game struct {
	font             font.Face
	backgroundLayers []Layer
	platforms        []*Platform
	player           *Player
	cameraX          float64
	score            int
	gameStarted      bool
	gameOver         bool
}

func NewGame() *Game {
	g := &Game{}
	g.initPlatforms()
	g.player = NewPlayer()
	g.font = basicfont.Face7x13 // Use the default basic font from Ebiten
	g.backgroundLayers = NewLayers()
	return g
}

func (g *Game) initPlatforms() {
	g.platforms = []*Platform{NewPlatform(startingPlatformX, startingPlatformY, startingPlatformWidth)}
	for i := 1; i < 2; i++ {
		g.platforms = append(g.platforms, GenerateNewRandomPlatform(g.platforms[i-1], g.score))
	}
}

////////////////////////////////////////////////////////////////////////

func (g *Game) Update() error {
	g.handleBackgroundLayers()
	if g.gameStarted {
		g.handlePlatforms()
		g.debug()
		g.checkGameOver()

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
			g.handlePlatformCollision(prevX)
			g.handleScreenBounds()
			g.handleCameraMovement()
		}
	}
	return nil
}

func (g *Game) handleBackgroundLayers() {
	for i := range g.backgroundLayers {
		g.backgroundLayers[i].OffsetX = -g.cameraX * g.backgroundLayers[i].Speed
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

func (g *Game) getLastPlatform() *Platform {
	return g.platforms[len(g.platforms)-1]
}

func (g *Game) distToFirstPlatform() float64 {
	firstPlatform := g.getFirstPlatform()
	return g.player.GetDist(firstPlatform)
}

func (g *Game) debug() {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		//fmt.Printf("Total Platforms: %d\n", len(g.platforms))
	}
}

// checkGameOver will set gameOver to true if Player fell too low
func (g *Game) checkGameOver() {
	if g.player.y >= screenHeight*2 {
		g.gameOver = true
	}
}

func (g *Game) getFirstPlatform() *Platform {
	return g.platforms[0]
}

func (g *Game) startOver() {
	g.resetGameState()
	g.player.ResetPlayer()
	g.initPlatforms()
}

// WARNING: getFirstPlatform will panic if you don't initialize the platforms after this
func (g *Game) resetGameState() {
	g.player.Reset()
	g.platforms = g.platforms[:0]
	g.cameraX = 0
	g.score = 0
	g.gameOver = false
}

func (g *Game) handlePlayer() {
	// todo make bot and player implement interface that applyGravity can use instead of checking for jumping keys
	if bot {
		g.botLogic()
		g.applyBotGravity()
	} else {
		g.playerControls()
		g.applyGravity()
	}
}

func (g *Game) botLogic() {
	g.player.AccelerateRight()
	if g.playerCanJump() && g.botShouldJump() {
		g.player.Jump()
	}
}

func (g *Game) playerCanJump() bool {
	return !g.playerInAir()
}

func (g *Game) playerInAir() bool {
	for _, p := range g.platforms {
		if g.playerOnPlatform(*p) {
			return false
		}
	}
	return true
}

func (g *Game) botShouldJump() bool {
	platformPos := g.nextUnvisitedPlatform()
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

func (g *Game) nextUnvisitedPlatform() *Platform {
	for _, p := range g.platforms {
		if !p.visited {
			return p
		}
	}
	return nil
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
	if (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)) && !g.player.isJumping {
		if g.playerCanJump() {
			g.player.Jump()
		}
	}
}

func (g *Game) slowPlayer() {
	g.player.velocityX *= .8
}

func (g *Game) applyGravity() {
	// todo make this use a function that can be called by applyBotGravity and heightAfterXFramesOfJumping
	currentGravity := gravity
	if g.player.velocityY > 0 {
		currentGravity = heavyGravity
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		currentGravity = lightGravity
	}
	g.player.velocityY += currentGravity
	g.player.y += g.player.velocityY
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

////////////////////////////////////////////////////////////////////////

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackgroundLayers(screen)

	// Title
	if !g.gameStarted {
		g.drawTitleText(screen, "Click anywhere to start", 200)
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.gameStarted = true
		}
	} else {
		// Game Started
		if g.gameOver {
			g.drawGameOverScreen(screen)
		}

		if bot {
			g.drawBotScreen(screen)
		}
		g.drawScore(screen)
		g.drawPlatforms(screen)
		g.player.Draw(screen, g.cameraX)
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

func (g *Game) drawGameOverScreen(screen *ebiten.Image) {
	g.drawText(screen, "That's not a rooftop geocode!", 60)
	g.drawText(screen, "Press Enter to Restart", 90)
}

func (g *Game) drawBotScreen(screen *ebiten.Image) {
	g.drawText(screen, "You can't be trusted to do it yourself.", 60)
	g.drawText(screen, "Now you have to use Smarty.", 80)
}

func (g *Game) drawText(screen *ebiten.Image, content string, y int) {
	textWidth := font.MeasureString(g.font, content).Ceil()
	x := (screenWidth - textWidth) / 2
	text.Draw(screen, content, g.font, x, y, colorText)
}

func (g *Game) drawTitleText(screen *ebiten.Image, content string, y int) {
	// Measure text size
	bounds := text.BoundString(g.font, content)
	width := bounds.Dx()
	height := bounds.Dy()
	// Create a temporary image to draw the original text
	textImage := ebiten.NewImage(width, height)
	text.Draw(textImage, content, g.font, 0, g.font.Metrics().Ascent.Ceil(), colorText)

	// Scale the image 2x
	scale := 2.0
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)

	// Center it horizontally on the screen
	x := float64((screenWidth - int(float64(width)*scale)) / 2)
	op.GeoM.Translate(x, float64(y))

	screen.DrawImage(textImage, op)
}

func (g *Game) drawPlatforms(screen *ebiten.Image) {
	for _, p := range g.platforms {
		p.Draw(screen, g.cameraX)
	}
}

func (g *Game) drawScore(screen *ebiten.Image) {
	text.Draw(screen, "Rooftops Geocoded: "+strconv.Itoa(g.score), g.font, 10, 20, colorText)
}

////////////////////////////////////////////////////////////////////////

func (g *Game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }
