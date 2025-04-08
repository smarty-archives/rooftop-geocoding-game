//go:build js && wasm
// +build js,wasm

package clipboard

import (
	"syscall/js"
)

func CopyToClipboard(text string) {
	// Get access to the browser's navigator.clipboard API
	navigator := js.Global().Get("navigator")
	if !navigator.Truthy() {
		println("Clipboard API not available")
		return
	}

	clipboard := navigator.Get("clipboard")
	if !clipboard.Truthy() {
		println("Clipboard not available")
		return
	}

	// Use clipboard.writeText(text)
	promise := clipboard.Call("writeText", text)

	// Handle success or failure (optional)
	then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("Copied to clipboard!")
		return nil
	})
	catch := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("Failed to copy to clipboard")
		return nil
	})

	promise.Call("then", then).Call("catch", catch)
}
