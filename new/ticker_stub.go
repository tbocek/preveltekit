//go:build !wasm

package reactive

// SetInterval stub for non-WASM builds - does nothing during pre-render.
func SetInterval(ms int, callback func()) func() {
	return func() {}
}

// SetTimeout stub for non-WASM builds - does nothing during pre-render.
func SetTimeout(ms int, callback func()) func() {
	return func() {}
}

// Debounce stub for non-WASM builds - executes callback immediately.
func Debounce(ms int, callback func()) func() {
	return callback
}

// Throttle stub for non-WASM builds - executes callback immediately.
func Throttle(ms int, callback func()) func() {
	return callback
}
