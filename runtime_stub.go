//go:build !wasm

package preveltekit

// IsBuildTime is true when running native (pre-rendering).
const IsBuildTime = true

// Stub implementations for non-WASM builds (SSR/pre-rendering)
// These are no-ops since DOM manipulation only happens in the browser

// document is a no-op stub for SSR (uses jsValue from js_stub.go)
var document = &jsValue{}

func getEl(id string) *jsValue           { return &jsValue{} }
func findComment(marker string) *jsValue { return &jsValue{} }

// cleanupBag holds js.Func references for batch release (stub for SSR).
type cleanupBag struct{}

// Add is a no-op for SSR.
func (c *cleanupBag) Add(fn jsFunc) {}

// AddDestroy is a no-op for SSR.
func (c *cleanupBag) AddDestroy(fn func()) {}

// Release is a no-op for SSR.
func (c *cleanupBag) Release() {}

// jsFunc is a stub for js.Func in non-WASM builds
type jsFunc struct{}

// bindable is implemented by types that can be bound to DOM elements.
type bindable[T any] interface {
	Get() T
	OnChange(func(T))
}

// settable extends bindable with Set capability for two-way binding
type settable[T any] interface {
	bindable[T]
	Set(T)
}

func bindInput(id string, store settable[string]) jsFunc  { return jsFunc{} }
func bindInputInt(id string, store settable[int]) jsFunc  { return jsFunc{} }
func bindCheckbox(id string, store settable[bool]) jsFunc { return jsFunc{} }

// Batch binding types and functions (stubs for SSR)
type evt struct {
	ID    string
	Event string
	Fn    func()
}

func bindEvents(c *cleanupBag, events []evt) {}

type inp struct {
	ID    string
	Store settable[string]
}

func bindInputs(c *cleanupBag, bindings []inp) {}

type chk struct {
	ID    string
	Store settable[bool]
}

func bindCheckboxes(c *cleanupBag, bindings []chk) {}
