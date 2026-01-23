//go:build !js || !wasm

package reactive

// Ticker holds a JavaScript interval that fires at regular intervals.
// This is a stub for non-WASM builds.
type Ticker struct {
	C chan struct{}
}

// NewTicker creates a new Ticker (stub for non-WASM builds).
func NewTicker(ms int) *Ticker {
	return &Ticker{C: make(chan struct{})}
}

// Stop stops the ticker (stub for non-WASM builds).
func (t *Ticker) Stop() {}

// SetInterval creates a JavaScript interval (stub for non-WASM builds).
func SetInterval(ms int, callback func()) func() {
	return func() {}
}

// SetTimeout calls the callback after ms milliseconds (stub for non-WASM builds).
func SetTimeout(ms int, callback func()) func() {
	return func() {}
}
