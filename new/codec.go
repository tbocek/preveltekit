//go:build wasm

package preveltekit

import (
	"syscall/js"
)

// JSDecoder is implemented by types that can decode from JS
type JSDecoder interface {
	FromJS(js.Value)
}

// JSEncoder is implemented by types that can encode to JS
type JSEncoder interface {
	ToJS() js.Value
}

// Decode converts a js.Value to a Go value.
// If dst implements JSDecoder, uses that.
func Decode(v js.Value, dst any) error {
	if v.IsUndefined() || v.IsNull() {
		return nil
	}
	if dec, ok := dst.(JSDecoder); ok {
		dec.FromJS(v)
	}
	return nil
}

// Encode converts a Go value to a js.Value.
// If src implements JSEncoder, uses that.
func Encode(src any) js.Value {
	if src == nil {
		return js.Null()
	}
	if enc, ok := src.(JSEncoder); ok {
		return enc.ToJS()
	}
	return js.ValueOf(src)
}
