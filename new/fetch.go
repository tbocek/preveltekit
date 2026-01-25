//go:build wasm

package reactive

import (
	"errors"
	"syscall/js"
)

// fetchSync performs a synchronous HTTP request and returns the JSON response as js.Value.
// Must be called from a goroutine (not the main thread).
func fetchSync(method, url string, body js.Value) (js.Value, error) {
	done := make(chan struct{})
	var result js.Value
	var fetchErr error

	opts := js.Global().Get("Object").New()
	opts.Set("method", method)
	if !body.IsUndefined() && !body.IsNull() {
		opts.Set("body", js.Global().Get("JSON").Call("stringify", body))
		headers := js.Global().Get("Object").New()
		headers.Set("Content-Type", "application/json")
		opts.Set("headers", headers)
	}

	promise := js.Global().Call("fetch", url, opts)

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
		})).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) any {
			// No JSON body (e.g., 204 No Content)
			result = js.Undefined()
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

// Get fetches JSON from a URL and decodes it into a typed struct.
// Must be called from a goroutine.
func Get[T any](url string) (T, error) {
	var result T
	jsVal, err := fetchSync("GET", url, js.Undefined())
	if err != nil {
		return result, err
	}
	if err := Decode(jsVal, &result); err != nil {
		return result, err
	}
	return result, nil
}

// Post sends a POST request with JSON body and decodes the response.
// Must be called from a goroutine.
func Post[T any](url string, body any) (T, error) {
	var result T
	jsVal, err := fetchSync("POST", url, Encode(body))
	if err != nil {
		return result, err
	}
	if err := Decode(jsVal, &result); err != nil {
		return result, err
	}
	return result, nil
}

// Put sends a PUT request with JSON body and decodes the response.
// Must be called from a goroutine.
func Put[T any](url string, body any) (T, error) {
	var result T
	jsVal, err := fetchSync("PUT", url, Encode(body))
	if err != nil {
		return result, err
	}
	if err := Decode(jsVal, &result); err != nil {
		return result, err
	}
	return result, nil
}

// Patch sends a PATCH request with JSON body and decodes the response.
// Must be called from a goroutine.
func Patch[T any](url string, body any) (T, error) {
	var result T
	jsVal, err := fetchSync("PATCH", url, Encode(body))
	if err != nil {
		return result, err
	}
	if err := Decode(jsVal, &result); err != nil {
		return result, err
	}
	return result, nil
}

// Delete sends a DELETE request and decodes the response.
// Must be called from a goroutine.
func Delete[T any](url string) (T, error) {
	var result T
	jsVal, err := fetchSync("DELETE", url, js.Undefined())
	if err != nil {
		return result, err
	}
	if err := Decode(jsVal, &result); err != nil {
		return result, err
	}
	return result, nil
}
