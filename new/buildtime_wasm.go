//go:build wasm

package preveltekit

// IsBuildTime is always false in WASM - we're running in the browser.
const IsBuildTime = false
