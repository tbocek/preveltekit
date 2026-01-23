//go:build js && wasm

package reactive

// IsBuildTime is always false in WASM - we're running in the browser.
const IsBuildTime = false
