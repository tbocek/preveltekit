//go:build !wasm

package preveltekit

import "errors"

var errWasmOnly = errors.New("fetch only works in wasm")

// FetchError provides detailed HTTP error information
type FetchError struct {
	Status     int
	StatusText string
	URL        string
}

func (e *FetchError) Error() string {
	return errWasmOnly.Error()
}

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
