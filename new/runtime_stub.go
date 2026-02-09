//go:build !wasm

package preveltekit

// IsBuildTime is true when running native (pre-rendering).
const IsBuildTime = true

// Stub implementations for non-WASM builds (SSR/pre-rendering)
// These are no-ops since DOM manipulation only happens in the browser

// Document is a no-op stub for SSR (uses jsValue from js_stub.go)
var Document = &jsValue{}

func GetEl(id string) *jsValue           { return &jsValue{} }
func FindComment(marker string) *jsValue { return &jsValue{} }

// Cleanup holds js.Func references for batch release (stub for SSR).
type Cleanup struct{}

// Add is a no-op for SSR.
func (c *Cleanup) Add(fn jsFunc) {}

// AddDestroy is a no-op for SSR.
func (c *Cleanup) AddDestroy(fn func()) {}

// Release is a no-op for SSR.
func (c *Cleanup) Release() {}

// jsFunc is a stub for js.Func in non-WASM builds
type jsFunc struct{}

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

func BindInput(id string, store Settable[string]) jsFunc  { return jsFunc{} }
func BindInputInt(id string, store Settable[int]) jsFunc  { return jsFunc{} }
func BindCheckbox(id string, store Settable[bool]) jsFunc { return jsFunc{} }

// Batch binding types and functions (stubs for SSR)
type Evt struct {
	ID    string
	Event string
	Fn    func()
}

func BindEvents(c *Cleanup, events []Evt) {}

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
