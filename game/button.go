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
	clickableMargin  float64
	drawMoreStrategy ButtonDrawStrategy
}

func NewButton(centerX, centerY, width, height, clickableMargin float64, btnFunc func()) *Button {
	return &Button{
		Pos: Pos{
			x: centerX - width/2,
			y: centerY - height/2,
		},
		width:            width,
		height:           height,
		buttonFn:         btnFunc,
		clickableMargin:  clickableMargin,
		drawMoreStrategy: BaseStrategy{},
	}
}

func NewTextButton(centerX, centerY, width, height float64, btnFunc func(), text string) *Button {
	b := NewButton(centerX, centerY, width, height, 0, btnFunc)
	b.drawMoreStrategy = NewTextStrategy(text)
	return b
}

func NewImageButton(centerX, centerY, width, height, scale, clickableMargin float64, btnFunc func(), imageFunc func() *ebiten.Image) *Button {
	b := NewButton(centerX, centerY, width, height, clickableMargin, btnFunc)
	b.drawMoreStrategy = NewImageStrategy(imageFunc, scale)
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
	return x < b.x+b.width+b.clickableMargin &&
		x > b.x-b.clickableMargin &&
		y < b.y+b.height+b.clickableMargin &&
		y > b.y-b.clickableMargin
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
	imageFunc func() *ebiten.Image
	scale     float64
}

func NewImageStrategy(imageFunc func() *ebiten.Image, scale float64) ImageStrategy {
	return ImageStrategy{
		imageFunc: imageFunc,
		scale:     scale,
	}
}

func (is ImageStrategy) DrawButton(screen *ebiten.Image, button *Button) {
	image := is.imageFunc()
	imageOptions := &ebiten.DrawImageOptions{}
	imageOptions.GeoM.Scale(is.scale, is.scale)
	imageOptions.GeoM.Translate(button.x, button.y)
	screen.DrawImage(image, imageOptions)
}
