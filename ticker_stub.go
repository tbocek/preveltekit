//go:build !wasm

package preveltekit

// SetInterval stub for non-WASM builds - does nothing during pre-render.
func SetInterval(ms int, callback func()) func() {
	return func() {}
}

// SetTimeout stub for non-WASM builds - does nothing during pre-render.
func SetTimeout(ms int, callback func()) func() {
	return func() {}
}

// Debounce stub for non-WASM builds - executes callback immediately.
// Returns the debounced function and a no-op cleanup function.
func Debounce(ms int, callback func()) (func(), func()) {
	return callback, func() {}
}

// Throttle stub for non-WASM builds - executes callback immediately.
func Throttle(ms int, callback func()) func() {
	return callback
}
