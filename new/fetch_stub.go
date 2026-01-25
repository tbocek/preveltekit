//go:build !wasm

package preveltekit

import "errors"

var errWasmOnly = errors.New("fetch only works in wasm")

// Get stub for non-WASM builds.
func Get[T any](url string) (T, error) {
	var zero T
	return zero, errWasmOnly
}

// Post stub for non-WASM builds.
func Post[T any](url string, body any) (T, error) {
	var zero T
	return zero, errWasmOnly
}

// Put stub for non-WASM builds.
func Put[T any](url string, body any) (T, error) {
	var zero T
	return zero, errWasmOnly
}

// Patch stub for non-WASM builds.
func Patch[T any](url string, body any) (T, error) {
	var zero T
	return zero, errWasmOnly
}

// Delete stub for non-WASM builds.
func Delete[T any](url string) (T, error) {
	var zero T
	return zero, errWasmOnly
}
