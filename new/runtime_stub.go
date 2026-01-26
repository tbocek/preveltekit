//go:build !js || !wasm

package preveltekit

// Stub implementations for non-WASM builds (SSR/pre-rendering)
// These are no-ops since DOM manipulation only happens in the browser

type jsValue struct{}

func (jsValue) IsUndefined() bool { return true }
func (jsValue) IsNull() bool      { return true }

// Document is a no-op stub for SSR
var Document = jsValue{}

func InjectStyle(name, css string)      {}
func GetEl(id string) jsValue           { return jsValue{} }
func FindComment(marker string) jsValue { return jsValue{} }

// Cleanup holds js.Func references for batch release (stub for SSR).
type Cleanup struct{}

// Add is a no-op for SSR.
func (c *Cleanup) Add(fn jsFunc) {}

// Release is a no-op for SSR.
func (c *Cleanup) Release() {}

// jsFunc is a stub for js.Func in non-WASM builds
type jsFunc struct{}

// On is a no-op for SSR, returns zero-value jsFunc.
func On(el jsValue, event string, handler func()) jsFunc { return jsFunc{} }

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
func BindInput(id string, store Settable[string]) jsFunc                       { return jsFunc{} }
func BindInputInt(id string, store Settable[int]) jsFunc                       { return jsFunc{} }
func BindCheckbox(id string, store Settable[bool]) jsFunc                      { return jsFunc{} }
func ToggleClass(el jsValue, class string, add bool)                           {}
func ReplaceContent(anchorMarker string, current jsValue, html string) jsValue { return jsValue{} }
func FindExistingIfContent(anchorMarker string) jsValue                        { return jsValue{} }

// Batch binding types and functions (stubs for SSR)
type Evt struct {
	ID    string
	Event string
	Fn    func()
}

func BindEvents(c *Cleanup, events []Evt) {}

type Txt[T any] struct {
	Marker string
	Store  Bindable[T]
}

func BindTexts[T any](bindings []Txt[T]) {}

type Inp struct {
	ID    string
	Store Settable[string]
}

func BindInputs(c *Cleanup, bindings []Inp) {}

type Chk struct {
	ID    string
	Store Settable[bool]
}

func BindCheckboxes(c *Cleanup, bindings []Chk) {}
