//go:build js && wasm
// +build js,wasm

package clipboard

import (
	"syscall/js"
)

func fallbackCopyTextToClipboard(text string) {
	document := js.Global().Get("document")
	textarea := document.Call("createElement", "textarea")
	textarea.Set("value", text)

	document.Get("body").Call("appendChild", textarea)
	textarea.Call("focus")
	textarea.Call("select")

	success := js.Global().Get("document").Call("execCommand", "copy").Bool()
	if success {
		println("Copied using fallback!")
	} else {
		println("Fallback failed.")
	}
	textarea.Call("remove")
}

func CopyToClipboard(text string) {
	navigator := js.Global().Get("navigator")
	if navigator.Truthy() {
		clipboard := navigator.Get("clipboard")
		if clipboard.Truthy() {
			promise := clipboard.Call("writeText", text)

			then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				println("Copied to clipboard!")
				return nil
			})
			catch := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				println("Failed to copy via Clipboard API, using fallback.")
				fallbackCopyTextToClipboard(text)
				return nil
			})

			promise.Call("then", then).Call("catch", catch)
			return
		}
	}
	// Clipboard API not available, use fallback
	fallbackCopyTextToClipboard(text)
}
