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

// textNodeRefs stores the current text node reference for each marker (updated on rebind)
var textNodeRefs = make(map[string]js.Value)

// boundMarkers tracks which markers have had OnChange registered
var boundMarkers = make(map[string]bool)

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

// ok returns true if el is a valid element
func ok(el js.Value) bool {
	return !el.IsNull() && !el.IsUndefined()
}

// SetText sets textContent on an element if it exists
func SetText(el js.Value, text string) {
	if ok(el) {
		el.Set("textContent", text)
	}
}

// On adds an event listener to an element
func On(el js.Value, event string, handler func()) {
	if ok(el) {
		el.Call("addEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) any {
			handler()
			return nil
		}))
	}
}

// OnEvent adds an event listener with access to the event object
func OnEvent(el js.Value, event string, handler func(js.Value)) {
	if ok(el) {
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

// Bind binds any store to an element's textContent (legacy, uses getElementById).
func Bind[T any](id string, store Bindable[T]) {
	el := GetEl(id)
	store.OnChange(func(v T) { SetText(el, toString(v)) })
	SetText(el, toString(store.Get()))
}

// FindComment finds a comment node with the given marker text using TreeWalker
func FindComment(marker string) js.Value {
	walker := Document.Call("createTreeWalker",
		Document.Get("body"),
		js.ValueOf(128), // NodeFilter.SHOW_COMMENT
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

// bindMarker is the unified implementation for BindText and BindHTML.
// isHTML=false: binds to text node (nodeType 3), uses nodeValue
// isHTML=true: binds to element (nodeType 1), uses innerHTML
func bindMarker[T any](marker string, store Bindable[T], isHTML bool) {
	comment := FindComment(marker)
	if comment.IsNull() {
		return
	}
	var node js.Value
	prevSibling := comment.Get("previousSibling")
	nodeType := 3
	if isHTML {
		nodeType = 1
	}
	if !prevSibling.IsNull() && prevSibling.Get("nodeType").Int() == nodeType {
		node = prevSibling
	} else {
		if isHTML {
			node = Document.Call("createElement", "span")
			node.Set("innerHTML", toString(store.Get()))
		} else {
			node = Document.Call("createTextNode", toString(store.Get()))
		}
		comment.Get("parentNode").Call("insertBefore", node, comment)
	}
	// Remove comment marker after hydration (no longer needed)
	comment.Call("remove")
	textNodeRefs[marker] = node
	prop := "nodeValue"
	if isHTML {
		prop = "innerHTML"
	}
	node.Set(prop, toString(store.Get()))
	if !boundMarkers[marker] {
		boundMarkers[marker] = true
		store.OnChange(func(v T) {
			if n, ok := textNodeRefs[marker]; ok {
				n.Set(prop, toString(v))
			}
		})
	}
}

// BindText binds a store to a text node, using a comment marker for hydration.
func BindText[T any](marker string, store Bindable[T]) {
	bindMarker(marker, store, false)
}

// BindHTML binds a store to innerHTML, using a comment marker for hydration.
func BindHTML[T any](marker string, store Bindable[T]) {
	bindMarker(marker, store, true)
}

// Settable extends Bindable with Set capability for two-way binding
type Settable[T any] interface {
	Bindable[T]
	Set(T)
}

// BindInput binds a text input to a string store (two-way)
func BindInput(id string, store Settable[string]) {
	el := GetEl(id)
	if !ok(el) {
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
	if !ok(el) {
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
	if !ok(el) {
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
	if ok(el) {
		el.Get("classList").Call("toggle", class, add)
	}
}

// ReplaceContent replaces if-block content: removes old, inserts new HTML before anchor comment
func ReplaceContent(anchorMarker string, current js.Value, html string) js.Value {
	anchor := FindComment(anchorMarker)
	newEl := Document.Call("createElement", "span")
	newEl.Set("innerHTML", html)
	if ok(current) && current.Truthy() {
		current.Call("remove")
	}
	if ok(anchor) {
		anchor.Get("parentNode").Call("insertBefore", newEl, anchor)
	}
	return newEl
}
