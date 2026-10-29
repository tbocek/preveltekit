//go:build wasm

package preveltekit

import (
	"reflect"
	"syscall/js"

	"github.com/tbocek/preveltekit/v2/internal/codecflags"
)

// JSDecoder is implemented by types that can decode from JS
type JSDecoder interface {
	FromJS(js.Value)
}

// JSEncoder is implemented by types that can encode to JS
type JSEncoder interface {
	ToJS() js.Value
}

// Decode converts a js.Value to a Go value.
// If dst implements JSDecoder, uses that.
// For structs, decodes from JS object using field tags.
// Generic so the JSDecoder check is resolved at compile time: if T implements
// JSDecoder, the reflection-based decoder is not linked in.
func Decode[T any](v js.Value, dst *T) error {
	if v.IsUndefined() || v.IsNull() || dst == nil {
		return nil
	}
	if dec, ok := any(dst).(JSDecoder); ok {
		dec.FromJS(v)
		return nil
	}
	decodeValue(v, reflect.ValueOf(dst).Elem())
	return nil
}

// decodeValue recursively decodes a js.Value into a reflect.Value
func decodeValue(v js.Value, dst reflect.Value) {
	if v.IsUndefined() || v.IsNull() {
		return
	}

	switch dst.Kind() {
	case reflect.Struct:
		t := dst.Type()
		for i := 0; i < dst.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue // skip unexported
			}
			// Get JS field name from tag
			jsName := field.Tag.Get("js")
			if jsName == "" {
				jsName = field.Tag.Get("json")
			}
			if jsName == "" || jsName == "-" {
				// Try field name as-is and lowercase first letter
				jsName = field.Name
			}

			jsVal := v.Get(jsName)
			if jsVal.IsUndefined() {
				// Try lowercase first letter
				if len(field.Name) > 0 {
					lcName := string(field.Name[0]|0x20) + field.Name[1:]
					jsVal = v.Get(lcName)
				}
			}
			if !jsVal.IsUndefined() && !jsVal.IsNull() {
				decodeValue(jsVal, dst.Field(i))
			}
		}

	case reflect.Ptr:
		if !codecflags.Ptr {
			codecDisabled("ptr")
			return
		}
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		decodeValue(v, dst.Elem())

	case reflect.Slice:
		if !codecflags.Slice {
			codecDisabled("slice")
			return
		}
		length := v.Length()
		slice := reflect.MakeSlice(dst.Type(), length, length)
		for i := range length {
			decodeValue(v.Index(i), slice.Index(i))
		}
		dst.Set(slice)

	case reflect.Map:
		if !codecflags.Map {
			codecDisabled("map")
			return
		}
		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		// Get keys using Object.keys()
		keys := js.Global().Get("Object").Call("keys", v)
		for i := 0; i < keys.Length(); i++ {
			key := keys.Index(i).String()
			val := reflect.New(dst.Type().Elem()).Elem()
			decodeValue(v.Get(key), val)
			dst.SetMapIndex(reflect.ValueOf(key), val)
		}

	case reflect.String:
		if !codecflags.String {
			codecDisabled("string")
			return
		}
		dst.SetString(v.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !codecflags.Int {
			codecDisabled("int")
			return
		}
		dst.SetInt(int64(v.Int()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !codecflags.Uint {
			codecDisabled("uint")
			return
		}
		dst.SetUint(uint64(v.Int()))

	case reflect.Float32, reflect.Float64:
		if !codecflags.Float {
			codecDisabled("float")
			return
		}
		dst.SetFloat(v.Float())

	case reflect.Bool:
		if !codecflags.Bool {
			codecDisabled("bool")
			return
		}
		dst.SetBool(v.Bool())

	case reflect.Interface:
		if !codecflags.Iface {
			codecDisabled("iface")
			return
		}
		// For interface{}, store the raw value as best we can
		switch v.Type() {
		case js.TypeString:
			dst.Set(reflect.ValueOf(v.String()))
		case js.TypeNumber:
			dst.Set(reflect.ValueOf(v.Float()))
		case js.TypeBoolean:
			dst.Set(reflect.ValueOf(v.Bool()))
		}
	}
}

// Encode converts a Go value to a js.Value.
// If src implements JSEncoder, uses that.
// For structs, converts to a JS object using field tags.
func Encode(src any) js.Value {
	if src == nil {
		return js.Null()
	}
	if enc, ok := src.(JSEncoder); ok {
		return enc.ToJS()
	}
	return encodeValue(reflect.ValueOf(src))
}

// encodeValue recursively encodes a reflect.Value to js.Value
func encodeValue(v reflect.Value) js.Value {
	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if !codecflags.Ptr {
			codecDisabled("ptr")
			return js.Undefined()
		}
		if v.IsNil() {
			return js.Null()
		}
		return encodeValue(v.Elem())
	}

	switch v.Kind() {
	case reflect.Struct:
		obj := js.Global().Get("Object").New()
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" {
				continue // skip unexported
			}
			// Get JS field name from tag or use lowercase first letter
			jsName := field.Tag.Get("js")
			if jsName == "" {
				jsName = field.Tag.Get("json")
			}
			if jsName == "" || jsName == "-" {
				// Use field name with lowercase first letter
				name := field.Name
				if len(name) > 0 {
					jsName = string(name[0]|0x20) + name[1:]
				}
			}
			obj.Set(jsName, encodeValue(v.Field(i)))
		}
		return obj

	case reflect.Slice, reflect.Array:
		if !codecflags.Slice {
			codecDisabled("slice")
			return js.Undefined()
		}
		arr := js.Global().Get("Array").New(v.Len())
		for i := 0; i < v.Len(); i++ {
			arr.SetIndex(i, encodeValue(v.Index(i)))
		}
		return arr

	case reflect.Map:
		if !codecflags.Map {
			codecDisabled("map")
			return js.Undefined()
		}
		obj := js.Global().Get("Object").New()
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			if key.Kind() == reflect.String {
				obj.Set(key.String(), encodeValue(iter.Value()))
			}
		}
		return obj

	case reflect.String:
		if !codecflags.String {
			codecDisabled("string")
			return js.Undefined()
		}
		return js.ValueOf(v.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !codecflags.Int {
			codecDisabled("int")
			return js.Undefined()
		}
		return js.ValueOf(v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !codecflags.Uint {
			codecDisabled("uint")
			return js.Undefined()
		}
		return js.ValueOf(v.Uint())

	case reflect.Float32, reflect.Float64:
		if !codecflags.Float {
			codecDisabled("float")
			return js.Undefined()
		}
		return js.ValueOf(v.Float())

	case reflect.Bool:
		if !codecflags.Bool {
			codecDisabled("bool")
			return js.Undefined()
		}
		return js.ValueOf(v.Bool())

	case reflect.Interface:
		if !codecflags.Iface {
			codecDisabled("iface")
			return js.Undefined()
		}
		if v.IsNil() {
			return js.Null()
		}
		return encodeValue(v.Elem())

	default:
		return js.Undefined()
	}
}

// codecDisabled reports a value whose kind was compiled out via a pk_no<kind>
// build tag (see internal/codecflags).
func codecDisabled(kind string) {
	js.Global().Get("console").Call("error", "preveltekit: "+kind+" support disabled by build tag pk_no"+kind)
}
