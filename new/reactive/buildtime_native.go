//go:build !js || !wasm

package reactive

// IsBuildTime is true when running native (pre-rendering).
const IsBuildTime = true
