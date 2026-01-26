//go:build !js || !wasm

package preveltekit

// Stub implementations for non-WASM builds (SSR/pre-rendering)
// These are no-ops since DOM manipulation only happens in the browser

type jsValue struct{}

func (jsValue) IsUndefined() bool { return true }
func (jsValue) IsNull() bool      { return true }

// Document is a no-op stub for SSR
var Document = jsValue{}

func InjectStyle(name, css string)                {}
func GetEl(id string) jsValue                     { return jsValue{} }
func FindComment(marker string) jsValue           { return jsValue{} }
func On(el jsValue, event string, handler func()) {}

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

func BindText[T any](marker string, store Bindable[T])                         {}
func BindHTML[T any](marker string, store Bindable[T])                         {}
func BindInput(id string, store Settable[string])                              {}
func BindInputInt(id string, store Settable[int])                              {}
func BindCheckbox(id string, store Settable[bool])                             {}
func ToggleClass(el jsValue, class string, add bool)                           {}
func ReplaceContent(anchorMarker string, current jsValue, html string) jsValue { return jsValue{} }
func FindExistingIfContent(anchorMarker string) jsValue                        { return jsValue{} }

// Batch binding types and functions (stubs for SSR)
type Evt struct {
	ID    string
	Event string
	Fn    func()
}

func BindEvents(events []Evt) {}

type Txt[T any] struct {
	Marker string
	Store  Bindable[T]
}

func BindTexts[T any](bindings []Txt[T]) {}

type Inp struct {
	ID    string
	Store Settable[string]
}

func BindInputs(bindings []Inp) {}

type Chk struct {
	ID    string
	Store Settable[bool]
}

func BindCheckboxes(bindings []Chk) {}
