//go:build wasm

package preveltekit

import (
	"errors"
	"reflect"
	"strings"
	"syscall/js"
)

// Decode converts a js.Value to a Go value using reflection.
// Supports structs with `js` tags, maps, slices, and primitive types.
func Decode(v js.Value, dst any) error {
	if v.IsUndefined() || v.IsNull() {
		return nil
	}

	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("decode: dst must be a non-nil pointer")
	}

	return decodeValue(v, rv.Elem())
}

func decodeValue(v js.Value, rv reflect.Value) error {
	if v.IsUndefined() || v.IsNull() {
		return nil
	}

	switch rv.Kind() {
	case reflect.String:
		rv.SetString(v.String())
	case reflect.Bool:
		rv.SetBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv.SetInt(int64(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		rv.SetUint(uint64(v.Int()))
	case reflect.Float32, reflect.Float64:
		rv.SetFloat(v.Float())
	case reflect.Struct:
		return decodeStruct(v, rv)
	case reflect.Slice:
		return decodeSlice(v, rv)
	case reflect.Map:
		return decodeMap(v, rv)
	case reflect.Ptr:
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		return decodeValue(v, rv.Elem())
	case reflect.Interface:
		rv.Set(reflect.ValueOf(jsToGo(v)))
	default:
		return errors.New("decode: unsupported type " + rv.Kind().String())
	}
	return nil
}

func decodeStruct(v js.Value, rv reflect.Value) error {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldVal := rv.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		jsName := field.Tag.Get("js")
		if jsName == "" {
			jsName = field.Name
		}
		if jsName == "-" {
			continue
		}
		jsVal := v.Get(jsName)
		if jsVal.IsUndefined() || jsVal.IsNull() {
			continue
		}
		if err := decodeValue(jsVal, fieldVal); err != nil {
			return err
		}
	}
	return nil
}

func decodeSlice(v js.Value, rv reflect.Value) error {
	length := v.Length()
	slice := reflect.MakeSlice(rv.Type(), length, length)
	for i := 0; i < length; i++ {
		if err := decodeValue(v.Index(i), slice.Index(i)); err != nil {
			return err
		}
	}
	rv.Set(slice)
	return nil
}

func decodeMap(v js.Value, rv reflect.Value) error {
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rv.Type()))
	}
	keys := js.Global().Get("Object").Call("keys", v)
	length := keys.Length()
	for i := 0; i < length; i++ {
		keyStr := keys.Index(i).String()
		jsVal := v.Get(keyStr)
		keyVal := reflect.New(rv.Type().Key()).Elem()
		if rv.Type().Key().Kind() == reflect.String {
			keyVal.SetString(keyStr)
		}
		elemVal := reflect.New(rv.Type().Elem()).Elem()
		if err := decodeValue(jsVal, elemVal); err != nil {
			return err
		}
		rv.SetMapIndex(keyVal, elemVal)
	}
	return nil
}

func jsToGo(v js.Value) any {
	if v.IsUndefined() || v.IsNull() {
		return nil
	}
	switch v.Type() {
	case js.TypeBoolean:
		return v.Bool()
	case js.TypeNumber:
		return v.Float()
	case js.TypeString:
		return v.String()
	case js.TypeObject:
		if js.Global().Get("Array").Call("isArray", v).Bool() {
			arr := make([]any, v.Length())
			for i := 0; i < v.Length(); i++ {
				arr[i] = jsToGo(v.Index(i))
			}
			return arr
		}
		m := make(map[string]any)
		keys := js.Global().Get("Object").Call("keys", v)
		for i := 0; i < keys.Length(); i++ {
			key := keys.Index(i).String()
			m[key] = jsToGo(v.Get(key))
		}
		return m
	default:
		return nil
	}
}

// Encode converts a Go value to a js.Value using reflection.
func Encode(src any) js.Value {
	if src == nil {
		return js.Null()
	}
	return encodeValue(reflect.ValueOf(src))
}

func encodeValue(rv reflect.Value) js.Value {
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return js.Null()
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.String:
		return js.ValueOf(rv.String())
	case reflect.Bool:
		return js.ValueOf(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return js.ValueOf(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return js.ValueOf(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return js.ValueOf(rv.Float())
	case reflect.Struct:
		return encodeStruct(rv)
	case reflect.Slice, reflect.Array:
		return encodeSlice(rv)
	case reflect.Map:
		return encodeMap(rv)
	case reflect.Interface:
		if rv.IsNil() {
			return js.Null()
		}
		return encodeValue(rv.Elem())
	default:
		return js.Undefined()
	}
}

func encodeStruct(rv reflect.Value) js.Value {
	obj := js.Global().Get("Object").New()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldVal := rv.Field(i)
		if !field.IsExported() {
			continue
		}
		jsName := field.Tag.Get("js")
		if jsName == "-" {
			continue
		}
		if jsName == "" {
			jsName = field.Name
		}
		if strings.HasSuffix(jsName, ",omitempty") {
			jsName = strings.TrimSuffix(jsName, ",omitempty")
			if fieldVal.IsZero() {
				continue
			}
		}
		obj.Set(jsName, encodeValue(fieldVal))
	}
	return obj
}

func encodeSlice(rv reflect.Value) js.Value {
	if rv.IsNil() {
		return js.Null()
	}
	arr := js.Global().Get("Array").New(rv.Len())
	for i := 0; i < rv.Len(); i++ {
		arr.SetIndex(i, encodeValue(rv.Index(i)))
	}
	return arr
}

func encodeMap(rv reflect.Value) js.Value {
	if rv.IsNil() {
		return js.Null()
	}
	obj := js.Global().Get("Object").New()
	iter := rv.MapRange()
	for iter.Next() {
		key := iter.Key()
		if key.Kind() == reflect.String {
			obj.Set(key.String(), encodeValue(iter.Value()))
		}
	}
	return obj
}
