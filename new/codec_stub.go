//go:build !wasm

package reactive

import "errors"

// Decode stub for non-WASM builds.
func Decode(v any, dst any) error {
	return errors.New("Decode only works in wasm")
}

// Encode stub for non-WASM builds.
func Encode(src any) any {
	return nil
}
