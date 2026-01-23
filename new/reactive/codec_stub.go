//go:build !js || !wasm

package reactive

import "errors"

// Decode is a stub for non-wasm builds
func Decode(v any, dst any) error {
	return errors.New("Decode only works in wasm")
}

// Encode is a stub for non-wasm builds
func Encode(src any) any {
	return nil
}

// Get is a stub for non-wasm builds
func Get[T any](url string) (T, error) {
	var zero T
	return zero, errors.New("Get only works in wasm")
}
