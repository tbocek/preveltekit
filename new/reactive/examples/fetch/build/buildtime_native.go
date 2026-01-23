//go:build !js || !wasm

package main

// IsBuildTime is true when running native (pre-rendering).
const IsBuildTime = true
