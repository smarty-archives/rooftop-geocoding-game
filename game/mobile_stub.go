//go:build !js || !wasm
// +build !js !wasm

package game

func RegisterClickHandler(fn func(x, y int)) (any, any) {
	return nil, nil
}

func IsMobile() bool {
	return false
}
