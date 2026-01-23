//go:build js && wasm

package reactive

import (
	"errors"
	"syscall/js"
)

// Fetch performs an HTTP GET request and calls the callback with the result.
// On success: callback(responseText, nil)
// On error: callback("", error)
func Fetch(url string, callback func(data string, err error)) {
	promise := js.Global().Call("fetch", url)

	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		response := args[0]
		if !response.Get("ok").Bool() {
			callback("", js.Error{Value: response.Get("statusText")})
			return nil
		}
		response.Call("text").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			callback(args[0].String(), nil)
			return nil
		}))
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
		callback("", js.Error{Value: args[0]})
		return nil
	}))
}

// FetchJSON performs an HTTP GET request and parses JSON response.
// On success: callback(jsonString, nil)
// On error: callback("", error)
func FetchJSON(url string, callback func(data string, err error)) {
	promise := js.Global().Call("fetch", url)

	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		response := args[0]
		if !response.Get("ok").Bool() {
			callback("", js.Error{Value: response.Get("statusText")})
			return nil
		}
		response.Call("json").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			// Convert JS object back to JSON string
			jsonStr := js.Global().Get("JSON").Call("stringify", args[0], nil, 2).String()
			callback(jsonStr, nil)
			return nil
		}))
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
		callback("", js.Error{Value: args[0]})
		return nil
	}))
}

// FetchSync performs a synchronous HTTP GET request and returns the JSON as js.Value.
// Must be called from a goroutine (not the main thread).
// Example:
//
//	go func() {
//	    data, err := reactive.FetchSync(url)
//	    if err != nil { ... }
//	    price := data.Get("RAW").Get("PRICE").Float()
//	}()
func FetchSync(url string) (js.Value, error) {
	done := make(chan struct{})
	var result js.Value
	var fetchErr error

	promise := js.Global().Call("fetch", url)

	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		response := args[0]
		if !response.Get("ok").Bool() {
			fetchErr = errors.New(response.Get("statusText").String())
			close(done)
			return nil
		}
		response.Call("json").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			result = args[0]
			close(done)
			return nil
		}))
		return nil
	})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
		fetchErr = errors.New(args[0].Get("message").String())
		close(done)
		return nil
	}))

	<-done
	return result, fetchErr
}
