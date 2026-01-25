//go:build !js || !wasm

package reactive

// Stub implementations for non-WASM builds (SSR/pre-rendering)
// These are no-ops since DOM manipulation only happens in the browser

type jsValue struct{}

func (jsValue) IsUndefined() bool { return true }
func (jsValue) IsNull() bool      { return true }

// Document is a no-op stub for SSR
var Document = jsValue{}

func InjectStyle(name, css string)                              {}
func GetEl(id string) jsValue                                   { return jsValue{} }
func SetText(el jsValue, text string)                           {}
func SetAttr(el jsValue, attr, val string)                      {}
func On(el jsValue, event string, handler func())               {}
func OnEvent(el jsValue, event string, handler func(e jsValue)) {}

// Bindable is implemented by types that can be bound to DOM elements.
type Bindable interface {
	Get() string
	OnChange(func(string))
}

// BindableInt is implemented by int stores.
type BindableInt interface {
	Get() int
	OnChange(func(int))
}

func Bind(id string, store Bindable)                                    {}
func BindInt(id string, store BindableInt)                              {}
func BindAttr(selector, attr, tmpl, field string, store *Store[string]) {}
func QuerySelector(selector string) jsValue                             { return jsValue{} }
func CreateElement(tag string) jsValue                                  { return jsValue{} }
