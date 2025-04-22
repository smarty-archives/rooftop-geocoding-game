//go:build js && wasm
// +build js,wasm

package game

import (
	"strings"
	"syscall/js"
)

const (
	gameWidth  = screenWidth  // or whatever your game's pixel width is
	gameHeight = screenHeight // same for height
)

func RegisterClickHandler(fn func(x, y int)) (any, any) {
	canvas := js.Global().Get("document").Call("querySelector", "canvas")
	if canvas.IsUndefined() {
		println("Canvas not found")
		return nil, nil
	}

	callback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		touches := event.Get("touches")
		if touches.Length() > 0 {
			touch := touches.Index(0)
			clientX := touch.Get("clientX").Float()
			clientY := touch.Get("clientY").Float()

			rect := canvas.Call("getBoundingClientRect")
			canvasLeft := rect.Get("left").Float()
			canvasTop := rect.Get("top").Float()
			displayWidth := rect.Get("width").Float()
			displayHeight := rect.Get("height").Float()

			// Calculate scale factor to fit game inside display while preserving aspect ratio
			scaleX := displayWidth / gameWidth
			scaleY := displayHeight / gameHeight
			scale := scaleX
			if scaleY < scaleX {
				scale = scaleY
			}

			// Calculate actual size of the rendered game and its offset
			renderedWidth := gameWidth * scale
			renderedHeight := gameHeight * scale
			offsetX := (displayWidth - renderedWidth) / 2
			offsetY := (displayHeight - renderedHeight) / 2

			// Coordinates relative to canvas
			canvasX := clientX - canvasLeft
			canvasY := clientY - canvasTop

			// Check if inside game rendering area
			//if canvasX < offsetX || canvasX > offsetX+renderedWidth ||
			//	canvasY < offsetY || canvasY > offsetY+renderedHeight {
			//	// Outside the game area — ignore
			//	return nil
			//}

			// Convert to game coordinates
			gameX := (canvasX - offsetX) / scale
			gameY := (canvasY - offsetY) / scale
			fn(int(gameX), int(gameY))
			isHeld = true
		}
		return nil
	})

	releaseCallback := js.FuncOf(func(this js.Value, args []js.Value) any {
		isHeld = false
		return nil
	})

	val := canvas.Call("addEventListener", "touchend", releaseCallback)
	val2 := canvas.Call("addEventListener", "touchstart", callback)
	return val, val2
}

func IsMobile() bool {
	navigator := js.Global().Get("navigator")
	if uaData := navigator.Get("userAgentData"); uaData.Truthy() {
		if uaData.Get("mobile").Bool() {
			return true
		}
	}

	ua := navigator.Get("userAgent").String()
	return strings.Contains(ua, "Mobi") || strings.Contains(ua, "Android") || strings.Contains(ua, "iPhone")
}
