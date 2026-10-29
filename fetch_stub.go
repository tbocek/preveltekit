//go:build !wasm

package preveltekit

import "errors"

var errWasmOnly = errors.New("fetch only works in wasm")

// FetchOptions configures a fetch request (stub for SSR)
type FetchOptions struct {
	Method  string
	Body    any
	Headers map[string]string
	Signal  any // placeholder for js.Value
}

// jsValue stub for non-WASM builds
type jsValueFetch struct{}

// NewAbortController stub for non-WASM builds.
func NewAbortController() (jsValueFetch, func()) {
	return jsValueFetch{}, func() {}
}

// Fetch stub for non-WASM builds.
func Fetch[T any](url string, opts *FetchOptions) (T, error) {
	var zero T
	return zero, errWasmOnly
}

func request[T any](_, _ string) (T, error) {
	var zero T
	return zero, errWasmOnly
}

func requestBody[T any](_, _ string, _ any) (T, error) {
	var zero T
	return zero, errWasmOnly
}
