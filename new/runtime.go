//go:build js && wasm

package preveltekit

import (
	"strconv"
	"syscall/js"
)

// Document is a cached reference to the DOM document
var Document = js.Global().Get("document")

// injectedStyles tracks which component styles have been injected
var injectedStyles = make(map[string]bool)

// InjectStyle injects a component's CSS once (deduplicated by name)
func InjectStyle(name, css string) {
	if injectedStyles[name] || css == "" {
		return
	}
	injectedStyles[name] = true
	style := Document.Call("createElement", "style")
	style.Set("textContent", css)
	Document.Get("head").Call("appendChild", style)
}

// GetEl returns an element by ID
func GetEl(id string) js.Value {
	return Document.Call("getElementById", id)
}

// SetText sets textContent on an element if it exists
func SetText(el js.Value, text string) {
	if !el.IsUndefined() && !el.IsNull() {
		el.Set("textContent", text)
	}
}

// On adds an event listener to an element
func On(el js.Value, event string, handler func()) {
	if !el.IsUndefined() && !el.IsNull() {
		el.Call("addEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) any {
			handler()
			return nil
		}))
	}
}

// OnEvent adds an event listener with access to the event object
func OnEvent(el js.Value, event string, handler func(e js.Value)) {
	if !el.IsUndefined() && !el.IsNull() {
		el.Call("addEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) any {
			handler(args[0])
			return nil
		}))
	}
}

// Bindable is implemented by types that can be bound to DOM elements.
type Bindable[T any] interface {
	Get() T
	OnChange(func(T))
}

// toString converts a value to string for display
func toString[T any](v T) string {
	switch val := any(v).(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

// Bind binds any store to an element's textContent.
func Bind[T any](id string, store Bindable[T]) {
	el := GetEl(id)
	store.OnChange(func(v T) { SetText(el, toString(v)) })
	SetText(el, toString(store.Get()))
}

// Settable extends Bindable with Set capability for two-way binding
type Settable[T any] interface {
	Bindable[T]
	Set(T)
}

// BindInput binds a text input to a string store (two-way)
func BindInput(id string, store Settable[string]) {
	el := GetEl(id)
	if el.IsNull() || el.IsUndefined() {
		return
	}
	el.Call("addEventListener", "input", js.FuncOf(func(this js.Value, args []js.Value) any {
		store.Set(this.Get("value").String())
		return nil
	}))
	store.OnChange(func(v string) { el.Set("value", v) })
}

// BindInputInt binds a text input to an int store (two-way)
func BindInputInt(id string, store Settable[int]) {
	el := GetEl(id)
	if el.IsNull() || el.IsUndefined() {
		return
	}
	el.Call("addEventListener", "input", js.FuncOf(func(this js.Value, args []js.Value) any {
		if v, err := strconv.Atoi(this.Get("value").String()); err == nil {
			store.Set(v)
		}
		return nil
	}))
	store.OnChange(func(v int) { el.Set("value", strconv.Itoa(v)) })
}

// BindCheckbox binds a checkbox to a bool store (two-way)
func BindCheckbox(id string, store Settable[bool]) {
	el := GetEl(id)
	if el.IsNull() || el.IsUndefined() {
		return
	}
	el.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) any {
		store.Set(this.Get("checked").Bool())
		return nil
	}))
	store.OnChange(func(v bool) { el.Set("checked", v) })
	el.Set("checked", store.Get())
}

// ToggleClass adds or removes a class based on a condition
func ToggleClass(el js.Value, class string, add bool) {
	if el.IsNull() || el.IsUndefined() {
		return
	}
	if add {
		el.Get("classList").Call("add", class)
	} else {
		el.Get("classList").Call("remove", class)
	}
}
