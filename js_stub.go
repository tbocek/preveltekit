//go:build !wasm

package preveltekit

// ssrPath holds the path currently being pre-rendered; the SSR router reads it
// in place of window.location.pathname.
var ssrPath string

// SetSSRPath sets the path the SSR router renders.
func SetSSRPath(path string) {
	ssrPath = path
}
