//go:build !wasm

package preveltekit

// ssrPath holds the current SSR path for simulating window.location.pathname
var ssrPath string

// SetSSRPath sets the path that js.Global().Get("location").Get("pathname") will return
func SetSSRPath(path string) {
	ssrPath = path
}

// jsValue is a stub for syscall/js.Value in native builds
type jsValue struct {
	data map[string]any
}

// jsGlobal is the fake global object
var jsGlobal = &jsValue{
	data: map[string]any{
		"location": &jsValue{
			data: map[string]any{},
		},
	},
}

// Global returns a fake JS global object for SSR
func jsGlobalFunc() *jsValue {
	return jsGlobal
}

// Get returns a nested value
func (v *jsValue) Get(key string) *jsValue {
	if key == "pathname" {
		return &jsValue{data: map[string]any{"_str": ssrPath}}
	}
	if val, ok := v.data[key]; ok {
		if jv, ok := val.(*jsValue); ok {
			return jv
		}
	}
	return &jsValue{data: map[string]any{}}
}

// String returns the string value
func (v *jsValue) String() string {
	if str, ok := v.data["_str"].(string); ok {
		return str
	}
	return ""
}

// Call is a no-op for SSR
func (v *jsValue) Call(method string, args ...any) *jsValue {
	return &jsValue{data: map[string]any{}}
}

// Set is a no-op for SSR
func (v *jsValue) Set(key string, val any) {}

// Bool returns false for SSR
func (v *jsValue) Bool() bool { return false }

// Int returns 0 for SSR
func (v *jsValue) Int() int { return 0 }

// IsNull returns false for SSR
func (v *jsValue) IsNull() bool { return false }

// IsUndefined returns false for SSR
func (v *jsValue) IsUndefined() bool { return false }
