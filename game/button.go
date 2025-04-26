package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

type Button struct {
	Pos
	width, height    float64
	buttonFn         func()
	isPressed        bool
	drawMoreStrategy ButtonDrawStrategy
}

func NewButton(centerX, centerY, width, height float64, btnFunc func()) *Button {
	return &Button{
		Pos: Pos{
			x: centerX - width/2,
			y: centerY - height/2,
		},
		width:            width,
		height:           height,
		buttonFn:         btnFunc,
		drawMoreStrategy: BaseStrategy{},
	}
}

func NewTextButton(centerX, centerY, width, height float64, btnFunc func(), text string) *Button {
	b := NewButton(centerX, centerY, width, height, btnFunc)
	b.drawMoreStrategy = NewTextStrategy(text)
	return b
}

func NewImageButton(centerX, centerY, width, height float64, btnFunc func(), image *ebiten.Image) *Button {
	b := NewButton(centerX, centerY, width, height, btnFunc)
	b.drawMoreStrategy = NewImageStrategy(image)
	return b
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.drawMoreStrategy.DrawButton(screen, b)
}

func (b *Button) getJustPressed() bool {
	return b.cursorOnButton() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (b *Button) Overlaps(x32, y32 int) bool {
	x := float64(x32)
	y := float64(y32)
	return x < b.x+b.width && x > b.x && y < b.y+b.height && y > b.y
}

func (b *Button) getIsPressed() bool {
	return b.cursorOnButton() && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

func (b *Button) cursorOnButton() bool {
	x32, y32 := ebiten.CursorPosition()
	return b.Overlaps(x32, y32)
}

func (b *Button) Update() {
	if b.getJustPressed() {
		b.buttonFn()
		b.isPressed = true
	} else {
		b.isPressed = false
	}
}

type ButtonDrawStrategy interface {
	DrawButton(screen *ebiten.Image, button *Button)
}

type BaseStrategy struct {
}

func (bs BaseStrategy) DrawButton(screen *ebiten.Image, b *Button) {
	vector.DrawFilledRect(
		screen,
		float32(b.x),
		float32(b.y),
		float32(b.width),
		float32(b.height),
		color.White,
		true,
	)
}

type TextStrategy struct {
	bg         ButtonDrawStrategy
	font       font.Face
	text       string
	textWidth  int
	textHeight int
}

func NewTextStrategy(text string) ButtonDrawStrategy {
	currFont := defaultFont
	return TextStrategy{
		bg:         BaseStrategy{},
		font:       currFont,
		text:       text,
		textWidth:  font.MeasureString(currFont, text).Ceil(),
		textHeight: currFont.Metrics().Ascent.Ceil(),
	}
}

func (ts TextStrategy) DrawButton(screen *ebiten.Image, button *Button) {
	ts.bg.DrawButton(screen, button)
	text.Draw(screen, ts.text, ts.font, int(button.width/2+button.x)-(ts.textWidth/2), int(button.height/2+button.y)+(ts.textHeight/2), colorText)
}

type ImageStrategy struct {
	image *ebiten.Image
}

func NewImageStrategy(image *ebiten.Image) ImageStrategy {
	return ImageStrategy{
		image: image,
	}
}

func (is ImageStrategy) DrawButton(screen *ebiten.Image, button *Button) {
	imageOptions := &ebiten.DrawImageOptions{}
	scaleX := button.width / float64(is.image.Bounds().Dx())
	scaleY := scaleX
	if button.getIsPressed() {
		scaleX *= .9
		scaleY *= .9
	}
	imageOptions.GeoM.Scale(scaleX, scaleY)
	imageOptions.GeoM.Translate(button.x, button.y)
	screen.DrawImage(is.image, imageOptions)
}
