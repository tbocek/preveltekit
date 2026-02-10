//go:build wasm

package preveltekit

import (
	"syscall/js"
)

// FetchError provides detailed HTTP error information
type FetchError struct {
	Status     int
	StatusText string
	URL        string
}

func (e *FetchError) Error() string {
	s := "fetch " + e.URL + ": "
	if e.StatusText != "" {
		return s + itoa(e.Status) + " " + e.StatusText
	}
	return s + "HTTP " + itoa(e.Status)
}

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

// fetchSync performs a synchronous HTTP request and returns the JSON response as js.Value.
// Must be called from a goroutine (not the main thread).
func fetchSync(method, url string, body js.Value) (js.Value, error) {
	return fetchSyncWithOpts(url, buildFetchOpts(method, body, js.Undefined()))
}

// buildFetchOpts creates JS fetch options object
func buildFetchOpts(method string, body, signal js.Value) js.Value {
	opts := js.Global().Get("Object").New()
	opts.Set("method", method)
	if !body.IsUndefined() && !body.IsNull() {
		opts.Set("body", js.Global().Get("JSON").Call("stringify", body))
		headers := js.Global().Get("Object").New()
		headers.Set("Content-Type", "application/json")
		opts.Set("headers", headers)
	}
	if !signal.IsUndefined() && !signal.IsNull() {
		opts.Set("signal", signal)
	}
	return opts
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
	var result T
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

	jsVal, err := fetchSyncWithOpts(url, jsOpts)
	if err != nil {
		return result, err
	}
	if err := Decode(jsVal, &result); err != nil {
		return result, err
	}
	return result, nil
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
