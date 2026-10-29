//go:build wasm

package preveltekit

import (
	"syscall/js"
)

// IsBuildTime is always false in WASM - we're running in the browser.
const IsBuildTime = false

// document is a cached reference to the DOM document
var document = js.Global().Get("document")

// nodeFilterShowComment is cached for TreeWalker (NodeFilter.SHOW_COMMENT = 128)
var nodeFilterShowComment = js.ValueOf(128)

// getEl returns an element by ID
func getEl(id string) js.Value {
	return document.Call("getElementById", id)
}

// exists returns true if el is a valid element
func exists(el js.Value) bool {
	return !el.IsNull() && !el.IsUndefined()
}

// cleanupBag holds js.Func references and destroy callbacks for batch release.
// Use this to prevent memory leaks when components unmount or re-render.
type cleanupBag struct {
	funcs     []js.Func
	onDestroy []func()
}

// Add registers a js.Func for later cleanup.
// Safe to call with zero-value js.Func.
func (c *cleanupBag) Add(fn js.Func) {
	if fn.Value.IsUndefined() {
		return
	}
	c.funcs = append(c.funcs, fn)
}

// AddDestroy registers a destroy callback to run on Release.
func (c *cleanupBag) AddDestroy(fn func()) {
	c.onDestroy = append(c.onDestroy, fn)
}

// Release runs all destroy callbacks and frees all registered js.Func references.
// Safe to call multiple times.
func (c *cleanupBag) Release() {
	for _, fn := range c.onDestroy {
		fn()
	}
	c.onDestroy = nil
	for _, fn := range c.funcs {
		fn.Release()
	}
	c.funcs = nil
}

// bindable is implemented by types that can be bound to DOM elements.
type bindable[T any] interface {
	Get() T
	OnChange(func(T))
}

// findComment finds a comment node with the given marker text using TreeWalker
func findComment(marker string) js.Value {
	walker := document.Call("createTreeWalker",
		document.Get("body"),
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

// settable extends bindable with Set capability for two-way binding
type settable[T any] interface {
	bindable[T]
	Set(T)
}

// bindInput binds a text input to a string store (two-way).
// Returns the js.Func for cleanup.
func bindInput(id string, store settable[string]) js.Func {
	el := getEl(id)
	if !exists(el) {
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

// bindInputInt binds a text input to an int store (two-way).
// Returns the js.Func for cleanup.
func bindInputInt(id string, store settable[int]) js.Func {
	el := getEl(id)
	if !exists(el) {
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

// bindCheckbox binds a checkbox to a bool store (two-way).
// Returns the js.Func for cleanup.
func bindCheckbox(id string, store settable[bool]) js.Func {
	el := getEl(id)
	if !exists(el) {
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

// bindEvent attaches a handler to the element with the given ID, applying
// the modifiers registered via PreventDefault/StopPropagation.
func bindEvent(c *cleanupBag, id, event string, handler func()) {
	el := getEl(id)
	if !exists(el) {
		return
	}
	mods := GetHandlerModifiers(id)
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
		handler()
		return nil
	})
	el.Call("addEventListener", event, fn)
	c.Add(fn)
}
