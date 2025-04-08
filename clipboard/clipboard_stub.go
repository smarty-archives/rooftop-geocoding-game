//go:build !js || !wasm
// +build !js !wasm

package clipboard

import "fmt"

func CopyToClipboard(text string) {
	fmt.Println("Yeah this ain't gonna work here")
}
