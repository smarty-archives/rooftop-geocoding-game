//go:build !js || !wasm
// +build !js !wasm

package game

func RegisterClickHandler(_ func(x, y int)) (any, any) {
	return nil, nil
}

func IsMobile() bool {
	return false
}
