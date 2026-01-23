//go:build js && wasm

package main

// IsBuildTime is always false in WASM - we're running in the browser.
const IsBuildTime = false
