//go:build !js || !wasm

package preveltekit

// Stub implementations for non-WASM builds (SSR/pre-rendering)
// These are no-ops since DOM manipulation only happens in the browser

type jsValue struct{}

func (jsValue) IsUndefined() bool { return true }
func (jsValue) IsNull() bool      { return true }

// Document is a no-op stub for SSR
var Document = jsValue{}

func InjectStyle(name, css string)                              {}
func GetEl(id string) jsValue                                   { return jsValue{} }
func SetText(el jsValue, text string)                           {}
func On(el jsValue, event string, handler func())               {}
func OnEvent(el jsValue, event string, handler func(e jsValue)) {}

// Bindable is implemented by types that can be bound to DOM elements.
type Bindable[T any] interface {
	Get() T
	OnChange(func(T))
}

func Bind[T any](id string, store Bindable[T]) {}
