//go:build js && wasm

package preveltekit

import (
	"syscall/js"
)

// Document is a cached reference to the DOM document
var Document = js.Global().Get("document")

// textNodeRefs stores the current text node reference for each marker (updated on rebind)
var textNodeRefs = make(map[string]js.Value)

// boundMarkers tracks which markers have had OnChange registered
var boundMarkers = make(map[string]bool)

// ClearBoundMarker removes a marker from the boundMarkers map.
// This allows re-binding when if-block content is replaced.
func ClearBoundMarker(marker string) {
	delete(boundMarkers, marker)
	delete(textNodeRefs, marker)
}

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

// toString converts a value to string for display
func toString[T any](v T) string {
	switch val := any(v).(type) {
	case string:
		return val
	case int:
		return itoa(val)
	case int64:
		return itoa(int(val))
	case int32:
		return itoa(int(val))
	case uint:
		return itoa(int(val))
	case uint64:
		return itoa(int(val))
	case uint32:
		return itoa(int(val))
	case float64:
		return floatToStr(val)
	case float32:
		return floatToStr(float64(val))
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return ""
	}
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

// bindMarker is the unified implementation for BindText and BindHTML.
// isHTML=false: binds to text node (nodeType 3), uses nodeValue
// isHTML=true: binds to element (nodeType 1), uses innerHTML
func bindMarker[T any](marker string, store Bindable[T], isHTML bool) {
	// Skip if already bound (comment was removed on first bind)
	if boundMarkers[marker] {
		return
	}
	boundMarkers[marker] = true

	comment := FindComment(marker)
	if comment.IsNull() {
		return
	}

	var node js.Value
	prevSibling := comment.Get("previousSibling")
	nodeType := 3
	prop := "nodeValue"
	if isHTML {
		nodeType = 1
		prop = "innerHTML"
	}
	if !prevSibling.IsNull() && prevSibling.Get("nodeType").Int() == nodeType {
		node = prevSibling
		// Set current value in case it changed after SSR (e.g., fetch results)
		node.Set(prop, toString(store.Get()))
	} else {
		if isHTML {
			node = Document.Call("createElement", "span")
			node.Set("innerHTML", toString(store.Get()))
		} else {
			node = Document.Call("createTextNode", toString(store.Get()))
		}
		comment.Get("parentNode").Call("insertBefore", node, comment)
	}
	comment.Call("remove")
	textNodeRefs[marker] = node
	store.OnChange(func(v T) {
		if n, ok := textNodeRefs[marker]; ok {
			n.Set(prop, toString(v))
		}
	})
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

// ToggleClass adds or removes a class based on a condition
func ToggleClass(el js.Value, class string, add bool) {
	if ok(el) {
		el.Get("classList").Call("toggle", class, add)
	}
}

// ReplaceContent replaces if-block content: removes old, inserts new HTML before anchor comment
func ReplaceContent(anchorMarker string, current js.Value, html string) js.Value {
	anchor := FindComment(anchorMarker)
	if !ok(anchor) {
		return js.Null()
	}
	parentNode := anchor.Get("parentNode")
	newEl := Document.Call("createElement", "span")
	newEl.Set("innerHTML", html)
	if current.Truthy() {
		current.Call("remove")
	}
	parentNode.Call("insertBefore", newEl, anchor)
	return newEl
}

// FindExistingIfContent finds the existing SSR-rendered content before an if-block anchor comment.
// Returns the element if found, or js.Null() if not found.
// This is used during hydration to avoid replacing pre-rendered content.
func FindExistingIfContent(anchorMarker string) js.Value {
	anchor := FindComment(anchorMarker)
	if !ok(anchor) {
		return js.Null()
	}
	// The SSR content is the previous sibling (a span element)
	prev := anchor.Get("previousSibling")
	if !prev.IsNull() && prev.Get("nodeType").Int() == 1 { // Element node
		return prev
	}
	return js.Null()
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

// Txt represents a text binding for batch processing
type Txt[T any] struct {
	Marker string
	Store  Bindable[T]
}

// BindTexts binds multiple text nodes in a loop
func BindTexts[T any](bindings []Txt[T]) {
	for _, b := range bindings {
		BindText(b.Marker, b.Store)
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
