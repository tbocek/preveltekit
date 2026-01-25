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
func FindComment(marker string) jsValue                         { return jsValue{} }
func SetText(el jsValue, text string)                           {}
func On(el jsValue, event string, handler func())               {}
func OnEvent(el jsValue, event string, handler func(e jsValue)) {}

// Bindable is implemented by types that can be bound to DOM elements.
type Bindable[T any] interface {
	Get() T
	OnChange(func(T))
}

// Settable extends Bindable with Set capability for two-way binding
type Settable[T any] interface {
	Bindable[T]
	Set(T)
}

func Bind[T any](id string, store Bindable[T])                                 {}
func BindText[T any](marker string, store Bindable[T])                         {}
func BindHTML[T any](marker string, store Bindable[T])                         {}
func BindInput(id string, store Settable[string])                              {}
func BindInputInt(id string, store Settable[int])                              {}
func BindCheckbox(id string, store Settable[bool])                             {}
func ToggleClass(el jsValue, class string, add bool)                           {}
func ReplaceContent(anchorMarker string, current jsValue, html string) jsValue { return jsValue{} }
