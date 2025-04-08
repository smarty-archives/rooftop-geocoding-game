package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	Pos
	width, height float64
}

func NewButton(centerX, centerY, width, height float64) *Button {
	return &Button{
		Pos: Pos{
			x: centerX - width/2,
			y: centerY - height/2,
		},
		width:  width,
		height: height,
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
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

func (b *Button) Pressed() bool {
	x32, y32 := ebiten.CursorPosition()
	x := float64(x32)
	y := float64(y32)
	if debugMode {
		fmt.Printf("x: %f, y: %f\n", x, y)
	}
	cursorOnButton := x < b.x+b.width && x > b.x && y < b.y+b.height && y > b.y
	return cursorOnButton && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}
