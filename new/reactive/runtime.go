//go:build js && wasm

package reactive

import (
	"strconv"
	"strings"
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

// SetAttr sets an attribute on an element if it exists
func SetAttr(el js.Value, attr, val string) {
	if !el.IsUndefined() && !el.IsNull() {
		el.Call("setAttribute", attr, val)
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

// Bind binds a string store to an element's textContent
func Bind(id string, store *Store[string]) {
	el := GetEl(id)
	store.OnChange(func(v string) { SetText(el, v) })
	SetText(el, store.Get())
}

// BindInt binds an int store to an element's textContent
func BindInt(id string, store *Store[int]) {
	el := GetEl(id)
	store.OnChange(func(v int) { SetText(el, strconv.Itoa(v)) })
	SetText(el, strconv.Itoa(store.Get()))
}

// BindAttr binds a string store to an element's attribute with template substitution
func BindAttr(selector, attr, tmpl, field string, store *Store[string]) {
	el := Document.Call("querySelector", selector)
	if el.IsUndefined() || el.IsNull() {
		return
	}
	update := func() {
		SetAttr(el, attr, strings.ReplaceAll(tmpl, "{"+field+"}", store.Get()))
	}
	store.OnChange(func(_ string) { update() })
	update()
}

// QuerySelector returns the first element matching a CSS selector
func QuerySelector(selector string) js.Value {
	return Document.Call("querySelector", selector)
}

// CreateElement creates a new DOM element
func CreateElement(tag string) js.Value {
	return Document.Call("createElement", tag)
}
