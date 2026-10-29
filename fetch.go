//go:build wasm

package preveltekit

import (
	"syscall/js"
)

// FetchOptions configures a fetch request
type FetchOptions struct {
	Method  string
	Body    any
	Headers map[string]string
	Signal  js.Value // AbortController.signal for cancellation
}

// NewAbortController creates a JS AbortController for request cancellation.
// Returns the signal to pass to FetchOptions and an abort function to cancel the request.
func NewAbortController() (signal js.Value, abort func()) {
	controller := js.Global().Get("AbortController").New()
	signal = controller.Get("signal")
	abort = func() {
		controller.Call("abort")
	}
	return
}

// fetchSyncWithOpts performs fetch with pre-built options
func fetchSyncWithOpts(url string, opts js.Value) (js.Value, error) {
	done := make(chan struct{})
	var result js.Value
	var fetchErr error

	// Track all funcs for cleanup to prevent memory leaks
	var funcs []js.Func
	cleanup := func() {
		for _, fn := range funcs {
			fn.Release()
		}
	}

	// Create callbacks before use to ensure they're tracked
	var jsonThen, jsonCatch js.Func

	jsonThen = js.FuncOf(func(this js.Value, args []js.Value) any {
		result = args[0]
		close(done)
		return nil
	})
	funcs = append(funcs, jsonThen)

	jsonCatch = js.FuncOf(func(this js.Value, args []js.Value) any {
		// No JSON body (e.g., 204 No Content)
		result = js.Undefined()
		close(done)
		return nil
	})
	funcs = append(funcs, jsonCatch)

	fetchThen := js.FuncOf(func(this js.Value, args []js.Value) any {
		response := args[0]
		if !response.Get("ok").Bool() {
			fetchErr = &FetchError{
				Status:     response.Get("status").Int(),
				StatusText: response.Get("statusText").String(),
				URL:        url,
			}
			close(done)
			return nil
		}
		response.Call("json").Call("then", jsonThen).Call("catch", jsonCatch)
		return nil
	})
	funcs = append(funcs, fetchThen)

	fetchCatch := js.FuncOf(func(this js.Value, args []js.Value) any {
		fetchErr = &FetchError{
			URL:        url,
			StatusText: args[0].Get("message").String(),
		}
		close(done)
		return nil
	})
	funcs = append(funcs, fetchCatch)

	js.Global().Call("fetch", url, opts).Call("then", fetchThen).Call("catch", fetchCatch)

	<-done
	cleanup()
	return result, fetchErr
}

// Fetch performs an HTTP request with full options including cancellation support.
// Must be called from a goroutine.
func Fetch[T any](url string, opts *FetchOptions) (T, error) {
	if opts == nil {
		opts = &FetchOptions{Method: "GET"}
	}
	if opts.Method == "" {
		opts.Method = "GET"
	}

	jsOpts := js.Global().Get("Object").New()
	jsOpts.Set("method", opts.Method)

	if opts.Body != nil {
		jsOpts.Set("body", js.Global().Get("JSON").Call("stringify", Encode(opts.Body)))
		if opts.Headers == nil {
			opts.Headers = make(map[string]string)
		}
		if _, exists := opts.Headers["Content-Type"]; !exists {
			opts.Headers["Content-Type"] = "application/json"
		}
	}

	if len(opts.Headers) > 0 {
		headers := js.Global().Get("Object").New()
		for k, v := range opts.Headers {
			headers.Set(k, v)
		}
		jsOpts.Set("headers", headers)
	}

	if !opts.Signal.IsUndefined() && !opts.Signal.IsNull() {
		jsOpts.Set("signal", opts.Signal)
	}

	return decodeResponse[T](fetchSyncWithOpts(url, jsOpts))
}

// request performs a request without body. Unlike Fetch it builds no header
// map, which keeps Get/Delete small.
func request[T any](method, url string) (T, error) {
	var result T
	v, err := fetchSync(method, url, js.Undefined())
	if err == nil {
		err = Decode(v, &result)
	}
	return result, err
}

// requestBody performs a request with a JSON body.
func requestBody[T any](method, url string, body any) (T, error) {
	var result T
	v, err := fetchSync(method, url, Encode(body))
	if err == nil {
		err = Decode(v, &result)
	}
	return result, err
}

// fetchSync performs a request with an optional JSON body.
func fetchSync(method, url string, body js.Value) (js.Value, error) {
	opts := js.Global().Get("Object").New()
	opts.Set("method", method)
	if !body.IsUndefined() && !body.IsNull() {
		opts.Set("body", js.Global().Get("JSON").Call("stringify", body))
		headers := js.Global().Get("Object").New()
		headers.Set("Content-Type", "application/json")
		opts.Set("headers", headers)
	}
	return fetchSyncWithOpts(url, opts)
}

// decodeResponse decodes the JSON response of a fetch into T.
func decodeResponse[T any](v js.Value, err error) (T, error) {
	var result T
	if err == nil {
		err = Decode(v, &result)
	}
	return result, err
}
