//go:build wasm

package preveltekit

import (
	"syscall/js"
)

// IsBuildTime is always false in WASM - we're running in the browser.
const IsBuildTime = false

// Document is a cached reference to the DOM document
var Document = js.Global().Get("document")

// nodeFilterShowComment is cached for TreeWalker (NodeFilter.SHOW_COMMENT = 128)
var nodeFilterShowComment = js.ValueOf(128)

// GetEl returns an element by ID
func GetEl(id string) js.Value {
	return Document.Call("getElementById", id)
}

// ok returns true if el is a valid element
func ok(el js.Value) bool {
	return !el.IsNull() && !el.IsUndefined()
}

// Cleanup holds js.Func references for batch release.
// Use this to prevent memory leaks when components unmount or re-render.
type Cleanup struct {
	funcs []js.Func
}

// Add registers a js.Func for later cleanup.
// Safe to call with zero-value js.Func.
func (c *Cleanup) Add(fn js.Func) {
	if fn.Value.IsUndefined() {
		return
	}
	c.funcs = append(c.funcs, fn)
}

// Release frees all registered js.Func references.
// Safe to call multiple times.
func (c *Cleanup) Release() {
	for _, fn := range c.funcs {
		fn.Release()
	}
	c.funcs = nil
}

// Bindable is implemented by types that can be bound to DOM elements.
type Bindable[T any] interface {
	Get() T
	OnChange(func(T))
}

// FindComment finds a comment node with the given marker text using TreeWalker
func FindComment(marker string) js.Value {
	walker := Document.Call("createTreeWalker",
		Document.Get("body"),
		nodeFilterShowComment,
		js.Null(),
	)
	for {
		node := walker.Call("nextNode")
		if node.IsNull() {
			return js.Null()
		}
		if node.Get("nodeValue").String() == marker {
			return node
		}
	}
}

// Settable extends Bindable with Set capability for two-way binding
type Settable[T any] interface {
	Bindable[T]
	Set(T)
}

// BindInput binds a text input to a string store (two-way).
// Returns the js.Func for cleanup.
func BindInput(id string, store Settable[string]) js.Func {
	el := GetEl(id)
	if !ok(el) {
		return js.Func{}
	}
	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		store.Set(this.Get("value").String())
		return nil
	})
	el.Call("addEventListener", "input", fn)
	store.OnChange(func(v string) { el.Set("value", v) })
	return fn
}

// BindInputInt binds a text input to an int store (two-way).
// Returns the js.Func for cleanup.
func BindInputInt(id string, store Settable[int]) js.Func {
	el := GetEl(id)
	if !ok(el) {
		return js.Func{}
	}
	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		store.Set(atoiSafe(this.Get("value").String()))
		return nil
	})
	el.Call("addEventListener", "input", fn)
	store.OnChange(func(v int) { el.Set("value", itoa(v)) })
	return fn
}

// BindCheckbox binds a checkbox to a bool store (two-way).
// Returns the js.Func for cleanup.
func BindCheckbox(id string, store Settable[bool]) js.Func {
	el := GetEl(id)
	if !ok(el) {
		return js.Func{}
	}
	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		store.Set(this.Get("checked").Bool())
		return nil
	})
	el.Call("addEventListener", "change", fn)
	store.OnChange(func(v bool) { el.Set("checked", v) })
	el.Set("checked", store.Get())
	return fn
}

// === Batch Binding Types (for smaller WASM) ===

// Evt represents an event binding for batch processing
type Evt struct {
	ID    string
	Event string
	Fn    func()
}

// BindEvents binds multiple events in a loop (smaller WASM than separate calls).
// Pass a Cleanup to collect js.Func references for later release.
func BindEvents(c *Cleanup, events []Evt) {
	for _, e := range events {
		e := e // Capture loop variable for closure
		el := GetEl(e.ID)
		if !ok(el) {
			continue
		}
		mods := GetHandlerModifiers(e.ID)
		fn := js.FuncOf(func(this js.Value, args []js.Value) any {
			if len(args) > 0 {
				ev := args[0]
				for _, mod := range mods {
					switch mod {
					case "preventDefault":
						ev.Call("preventDefault")
					case "stopPropagation":
						ev.Call("stopPropagation")
					}
				}
			}
			e.Fn()
			return nil
		})
		el.Call("addEventListener", e.Event, fn)
		c.Add(fn)
	}
}

// Inp represents an input binding for batch processing
type Inp struct {
	ID    string
	Store Settable[string]
}

// BindInputs binds multiple inputs in a loop.
// Pass a Cleanup to collect js.Func references for later release.
func BindInputs(c *Cleanup, bindings []Inp) {
	for _, b := range bindings {
		c.Add(BindInput(b.ID, b.Store))
	}
}

// Chk represents a checkbox binding for batch processing
type Chk struct {
	ID    string
	Store Settable[bool]
}

// BindCheckboxes binds multiple checkboxes in a loop.
// Pass a Cleanup to collect js.Func references for later release.
func BindCheckboxes(c *Cleanup, bindings []Chk) {
	for _, b := range bindings {
		c.Add(BindCheckbox(b.ID, b.Store))
	}
}
