//go:build js && wasm
// +build js,wasm

package game

import (
	"strings"
	"syscall/js"
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
			x := touch.Get("clientX").Int()
			y := touch.Get("clientY").Int()
			fn(x, y)
			//println("touch:", x, y)
			isHeld = true
		}
		return nil
	})

	releaseCallback := js.FuncOf(func(this js.Value, args []js.Value) any {
		//x := args[0].Get("clientX").Int()
		//y := args[0].Get("clientY").Int()
		//println("Clicked at:", x, y)
		isHeld = false
		return nil
	})

	val := canvas.Call("addEventListener", "touchend", releaseCallback)
	val2 := canvas.Call("addEventListener", "touchstart", callback)
	return val, val2 // return them and save them or GC could end them
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
