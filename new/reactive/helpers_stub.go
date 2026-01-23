//go:build !js || !wasm

package reactive

import "errors"

// Fetch is a stub for non-wasm builds (IDE support)
func Fetch(url string, callback func(data string, err error)) {
	panic("Fetch only works in wasm")
}

// FetchJSON is a stub for non-wasm builds (IDE support)
func FetchJSON(url string, callback func(data string, err error)) {
	panic("FetchJSON only works in wasm")
}

// FetchSync is a stub for non-wasm builds (IDE support)
func FetchSync(url string) (any, error) {
	return nil, errors.New("FetchSync only works in wasm")
}
