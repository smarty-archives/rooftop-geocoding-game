package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"slices"
	"sort"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/smarty-archives/rooftop-geocoding-game/clipboard"
	"github.com/smarty-archives/rooftop-geocoding-game/media"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth            = 640
	screenHeight           = 480
	startButtonCenterX     = screenWidth / 2
	startButtonCenterY     = 400
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

	gameLink = "https://www.smarty.com/geocode-jumper"
)

var (
	colorText                  = color.Black
	bot                        = false
	botFramesLeftJumping       = 0
	debugMode                  = false
	filler                 any = nil
	filler2                any = nil
	AudioInitialized           = false
	isHeld                     = false
	copiedSuccessCountdown     = 0
)

type Game struct {
	font             font.Face
	backgroundLayers []Layer
	platforms        []*Platform
	clouds           []*Cloud
	geocodes         []*Geocode
	player           *Player
	startButton      *Button
	shareButton      *Button
	muteButton       *Button
	cameraX          float64
	score            int
	gameStarted      bool
	gameOver         bool
	isMobile         bool
}

func (g *Game) Layout(_, _ int) (int, int) { return screenWidth, screenHeight }

var (
	defaultFont = basicfont.Face7x13 // Use the default basic font from Ebiten
)

func NewGame() *Game {
	g := &Game{}
	g.initClouds()
	g.initPlatforms()
	g.player = NewPlayer()
	g.font = defaultFont
	g.backgroundLayers = NewLayers()
	g.initButtons()
	g.isMobile = IsMobile()
	if g.isMobile {
		filler, filler2 = RegisterClickHandler(func(x, y int) {
			if g.muteButton.Overlaps(x, y) {
				g.muteButton.buttonFn()
				return
			}
			if !g.gameStarted {
				g.gameStarted = true
			} else if g.gameOver {
				if g.shareButton.Overlaps(x, y) && !g.isMobile {
					g.shareButton.buttonFn()
				} else {
					g.startOver()
				}
				return
			}
			g.Jump()
		})
	}
	return g
}

func (g *Game) initClouds() {
	g.clouds = []*Cloud{
		NewCloud(20, g.randomStartingCloudHeight(), .5),
	}
	for range 20 {
		g.clouds = append(g.clouds, g.generateNewRandomStartingCloud())
	}
}

func (g *Game) generateNewRandomStartingCloud() *Cloud {
	prevCloud := g.getLastCloud()
	randX := newRandomCloudX(prevCloud)
	randY := g.randomStartingCloudHeight()
	randSpeed := pickRand(.4, .5, .6, .7)
	return NewCloud(randX, randY, randSpeed)
}

func (g *Game) generateNewRandomCloud() *Cloud {
	prevCloud := g.getLastCloud()
	randX := newRandomCloudX(prevCloud)
	randY := g.randomCloudHeight()
	randSpeed := pickRand(.4, .5, .6, .7)
	return NewCloud(randX, randY, randSpeed)
}

func newRandomCloudX(prevCloud *Cloud) float64 {
	return prevCloud.x + float64(prevCloud.Image.Bounds().Dx()) + giveOrTake(100, 75)
}

func (g *Game) randomStartingCloudHeight() float64 {
	var cloudRange = .1
	return rand.Float64()*cloudRange*screenHeight + 20
}

func (g *Game) randomCloudHeight() float64 {
	var cloudRange float64
	if g.score > 90 {
		cloudRange = 1.0
	} else if g.score > 80 {
		cloudRange = .66
	} else if g.score > 70 {
		cloudRange = .5
	} else {
		cloudRange = .33
	}
	return rand.Float64()*cloudRange*screenHeight + 20
}

func (g *Game) initPlatforms() {
	g.platforms = []*Platform{NewPlatform(startingPlatformX, startingPlatformY, startingPlatformWidth)}
	for i := 1; i < 2; i++ {
		g.platforms = append(g.platforms, GenerateNewRandomPlatform(g.platforms[i-1], g.score))
	}
}

func (g *Game) initButtons() {
	g.startButton = NewImageButton(startButtonCenterX, startButtonCenterY, 187, 60, 1, 0, func() {
		g.gameStarted = true
	}, media.Instance.GetPlayButtonImage)

	g.shareButton = NewImageButton(startButtonCenterX, 400, 360, 60, 1, 0, func() {
		clipboard.CopyToClipboard(fmt.Sprintf("I scored %d on Geocode Jumper!\nTry to beat me\n%s", g.score, gameLink))
		copiedSuccessCountdown = 120
	}, GetShareButtonImage)

	g.muteButton = NewImageButton(screenWidth-30, 30, 24, 24, .5, 20, func() {
		ToggleMute()
	}, GetMuteButtonImage)
}

////////////////////////////////////////////////////////////////////////

func (g *Game) Update() error {
	g.debug()
	g.muteButton.Update()
	g.handleBackgroundLayers()
	g.handleBackgroundClouds()
	if g.gameStarted {
		g.handlePlatforms()
		g.checkGameOver()

		// If game over, reset the game when the enter key is pressed
		if g.gameOver {
			if !g.isMobile {
				g.shareButton.Update()
				if copiedSuccessCountdown > 0 {
					copiedSuccessCountdown--
				}
			}
			if ebiten.IsKeyPressed(ebiten.KeyEnter) {
				if g.player.x < g.getFirstPlatform().GetX() {
					bot = true
				}
				g.startOver()
			}
		} else {
			if bot && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
				g.startOver()
				bot = false
			}
			prevLeft := g.player.LeftX()
			prevRight := g.player.RightX()
			g.handlePlayer()
			g.handlePlatformCollision(prevLeft, prevRight)
			g.handleScreenBounds()
			g.handleCameraMovement()
		}
		g.handleGeocodes()
	} else { // Title Page
		g.startButton.Update()
	}
	return nil
}

func (g *Game) handleBackgroundLayers() {
	for i := range g.backgroundLayers {
		g.backgroundLayers[i].OffsetX = -g.cameraX * g.backgroundLayers[i].Speed
	}
}

func (g *Game) handleBackgroundClouds() {
	if len(g.clouds) == 0 {
		return
	}
	// Generate New
	if g.distToLastCloud() < screenWidth/2 {
		g.clouds = append(g.clouds, g.generateNewRandomCloud())
	}
	// Cleanup
	for i := range g.clouds {
		if g.clouds[i].getOffsetX(g.cameraX) < -screenWidth {
			g.clouds[i] = g.generateNewRandomCloud()
		}
	}
}

func (g *Game) distToLastCloud() float64 {
	return g.player.GetDist(g.getLastCloud())
}

func (g *Game) getLastCloud() *Cloud {
	if len(g.clouds) == 0 {
		return nil
	}
	lastCloud := g.clouds[0]
	for i := range g.clouds {
		if g.clouds[i].getOffsetX(g.cameraX) > lastCloud.getOffsetX(g.cameraX) {
			lastCloud = g.clouds[i]
		}
	}
	return lastCloud
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
	for i := range g.platforms {
		if g.platforms[i].visited {
			g.platforms[i].framesSinceVisited++
		}
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
		debugMode = true
	} else {
		debugMode = false
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
	g.initClouds()
}

// WARNING: getFirstPlatform will panic if you don't initialize the platforms after this
func (g *Game) resetGameState() {
	g.platforms = g.platforms[:0]
	g.clouds = g.clouds[:0]
	g.geocodes = g.geocodes[:0]
	g.cameraX = 0
	g.score = 0
	g.gameOver = false
	copiedSuccessCountdown = 0
}

func (g *Game) handlePlayer() {
	g.player.cycleImage() // this needs to be here so the player image is updated consistently regardless of frame rate
	// todo make bot and player implement interface that applyGravity can use instead of checking for jumping keys
	if g.isMobile {
		if g.score > 0 { // wait until you land on the 1st building
			g.player.AccelerateRight()
		}
		g.applyGravity()
	} else if bot {
		g.botLogic()
		g.applyBotGravity()
	} else {
		g.playerControls()
		g.applyGravity()
	}
}

func (g *Game) botLogic() {
	if g.botShouldAccelerateRight() {
		g.player.AccelerateRight()
	}
	if g.playerCanJump() && g.botShouldJump() {
		g.player.Jump()
	}
}

func (g *Game) botShouldAccelerateRight() bool {
	platform := g.nextUnvisitedPlatform()
	if g.player.x >= platform.x+(platform.width/2)-playerSize {
		return false
	}
	return true
}

func (g *Game) playerCanJump() bool {
	return !g.playerInAir()
}

func (g *Game) Jump() {
	if g.playerCanJump() {
		g.player.Jump()
	}
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
			return g.player.velocityX > 2
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

// heightAfterXFramesOfJumping assumes that the player is moving at top speed
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
		g.Jump()
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
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) || isHeld {
		currentGravity = lightGravity
	}
	g.player.velocityY += currentGravity
	g.player.y += g.player.velocityY
}

func (g *Game) handlePlatformCollision(prevLeft, prevRight float64) {
	for _, p := range g.platforms {
		// **Vertical collision (Landing on the platform)**
		if g.playerOnPlatform(*p) {
			// Land on the platform
			g.player.y = p.y - playerSize
			g.player.velocityY = 0
			g.player.isJumping = false

			if !p.visited {
				p.Visit()
				g.score++
				g.addGeocode()
			}
		}
		platformLeft := p.x
		platformRight := p.x + p.width
		// **Side collision (Hitting the sides of the platform)**
		if g.player.y+playerSize > p.y { // Player is below the platform surface
			if prevRight <= platformLeft && g.player.RightX() > platformLeft { // Hitting left side
				// Move player to edge of the platform
				g.player.x = platformLeft - g.player.width/2 - playerSize/2
				g.player.velocityX = 0
			} else if prevLeft >= platformRight && g.player.LeftX() < platformRight { // Hitting the right side
				// Move player to the edge of the platform
				g.player.x = platformRight + g.player.width/2 - playerSize/2
				g.player.velocityX = 0
			}
		}
	}
}

func (g *Game) playerOnPlatform(p Platform) bool {
	// Check if player is within platform's horizontal range
	platformRight := p.x + p.width
	platformLeft := p.x

	// **Vertical collision (Landing on the platform)**
	return g.player.RightX() > platformLeft && g.player.LeftX() < platformRight && // Player overlaps horizontally
		g.player.y+playerSize >= p.y && g.player.y+playerSize-g.player.velocityY <= p.y // Player is falling onto the platform
}

func (g *Game) addGeocode() {
	g.geocodes = append(g.geocodes, NewGeocode(Pos{
		x: g.player.GetCenterX(),
		y: g.player.GetY() - 20,
	}))
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

func (g *Game) handleGeocodes() {
	fadeRate := 5
	keep := g.geocodes[:0] // reuse the underlying array
	for _, geocode := range g.geocodes {
		geocode.opacity -= fadeRate
		if geocode.opacity > 0 {
			keep = append(keep, geocode)
		}
	}
	g.geocodes = keep
}

////////////////////////////////////////////////////////////////////////

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackgroundLayers(screen)

	if !g.gameStarted { // Title Page
		g.drawBackgroundClouds(screen)
		g.drawTitle(screen)
		g.startButton.Draw(screen)
	} else { // Game Started
		g.drawPlatforms(screen)
		g.player.Draw(screen, g.cameraX)
		g.drawBackgroundClouds(screen)
		if g.gameOver {
			if !bot {
				if !g.isMobile {
					g.shareButton.Draw(screen)
				}
				g.drawGameOverScreen(screen)
			}
		}
		if bot {
			g.drawBotScreen(screen)
		}
	}
	g.muteButton.Draw(screen)
	g.DrawAllText(screen)
}

func (g *Game) drawBackgroundLayers(screen *ebiten.Image) {
	for _, layer := range g.backgroundLayers {
		imgWidth := layer.Image.Bounds().Dx()
		imgHeight := layer.Image.Bounds().Dy()

		// Scale the image to match the screen height
		scaleY := float64(screenHeight) / float64(imgHeight)
		scaleX := scaleY // Maintain the aspect ratio horizontally

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

func (g *Game) drawBackgroundClouds(screen *ebiten.Image) {
	speedMap := map[float64][]*Cloud{}
	var speeds []float64
	for i := range g.clouds {
		s := g.clouds[i].Speed
		if !slices.Contains(speeds, s) {
			speeds = append(speeds, s)
		}
		speedMap[s] = append(speedMap[s], g.clouds[i])
	}
	sort.Float64s(speeds)
	for _, speed := range speeds {
		for _, cloud := range speedMap[speed] {
			cloud.Draw(screen, g.cameraX)
		}
	}
}

func (g *Game) drawTitle(screen *ebiten.Image) {
	drawImage(screen, media.Instance.GetTitleImage(), screenWidth/2, screenHeight/2-50)
}

func (g *Game) drawGameOverScreen(screen *ebiten.Image) {
	var restartImage *ebiten.Image
	if g.isMobile {
		restartImage = media.Instance.GetMobileRestartButtonImage()
	} else {
		restartImage = media.Instance.GetRestartButtonImage()
	}
	drawImage(screen, restartImage, screenWidth/2, screenHeight/2)
}

func (g *Game) drawBotScreen(screen *ebiten.Image) {
	g.drawTextCenteredOn(screen, "Smarty will take it from here.", screenWidth/2, 60)
	g.drawTextCenteredOn(screen, "Press enter if you want to go back to the hard way.", screenWidth/2, 80)
}

func (g *Game) drawTextCenteredOn(screen *ebiten.Image, content string, x, y int) {
	textWidth := font.MeasureString(g.font, content).Ceil()
	textHeight := g.font.Metrics().Ascent.Ceil()
	drawX := x - textWidth/2
	drawY := y - textHeight/2
	text.Draw(screen, content, g.font, drawX, drawY, colorText)
}

func (g *Game) drawPlatforms(screen *ebiten.Image) {
	for _, p := range g.platforms {
		p.Draw(screen, g.cameraX)
	}
}

func (g *Game) drawGeocodes(screen *ebiten.Image) {
	for _, geocode := range g.geocodes {
		geocode.Draw(screen, g.font, g.cameraX)
	}
}

func (g *Game) drawScore(screen *ebiten.Image) {
	text.Draw(screen, "Rooftops Geocoded: "+strconv.Itoa(g.score), g.font, 10, 20, colorText)
}

func (g *Game) DrawAllText(screen *ebiten.Image) {
	if g.gameStarted {
		g.drawGeocodes(screen)
		g.drawScore(screen)
	}
}

////////////////////////////////////////////////////////////////////////

func drawImage(screen *ebiten.Image, image *ebiten.Image, centerX, centerY float64) {
	titleOptions := &ebiten.DrawImageOptions{}
	scaleX, scaleY := 1.0, 1.0
	x := centerX - float64(image.Bounds().Dx()/2)
	y := centerY - float64(image.Bounds().Dy()/2)
	titleOptions.GeoM.Scale(scaleX, scaleY)
	titleOptions.GeoM.Translate(x, y)
	screen.DrawImage(image, titleOptions)
}

func GetMuteButtonImage() *ebiten.Image {
	if isMuted {
		return media.Instance.GetMutedImage()
	} else {
		return media.Instance.GetPlayingImage()
	}
}

func GetShareButtonImage() *ebiten.Image {
	if copiedSuccessCountdown > 0 {
		return media.Instance.GetCopyScoreSuccessButtonImage()
	} else {
		return media.Instance.GetCopyScorePromptButtonImage()
	}
}
