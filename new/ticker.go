//go:build wasm

package preveltekit

import (
	"syscall/js"
)

// SetInterval creates a JavaScript interval that calls the callback every ms milliseconds.
// Returns a function to clear the interval.
//
// Example:
//
//	stop := reactive.SetInterval(1000, func() {
//	    fmt.Println("tick")
//	})
//	// Later:
//	stop()
func SetInterval(ms int, callback func()) func() {
	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		callback()
		return nil
	})

	intervalID := js.Global().Call("setInterval", fn, ms)

	return func() {
		js.Global().Call("clearInterval", intervalID)
		fn.Release()
	}
}

// SetTimeout calls the callback after ms milliseconds.
// Returns a function to cancel the timeout before it fires.
//
// Example:
//
//	cancel := reactive.SetTimeout(5000, func() {
//	    fmt.Println("5 seconds passed")
//	})
//	// To cancel before it fires:
//	cancel()
func SetTimeout(ms int, callback func()) func() {
	var fn js.Func
	fn = js.FuncOf(func(this js.Value, args []js.Value) any {
		callback()
		fn.Release()
		return nil
	})

	timeoutID := js.Global().Call("setTimeout", fn, ms)

	return func() {
		js.Global().Call("clearTimeout", timeoutID)
		fn.Release()
	}
}

// Debounce returns a debounced version of the callback that delays execution
// until ms milliseconds have passed without another call.
//
// Example:
//
//	search := reactive.Debounce(300, func() {
//	    // This runs 300ms after the last keystroke
//	    performSearch(input.Get())
//	})
//	input.OnChange(func(_ string) { search() })
func Debounce(ms int, callback func()) func() {
	var timeoutID js.Value
	var fn js.Func

	fn = js.FuncOf(func(this js.Value, args []js.Value) any {
		callback()
		return nil
	})

	return func() {
		if !timeoutID.IsUndefined() && !timeoutID.IsNull() {
			js.Global().Call("clearTimeout", timeoutID)
		}
		timeoutID = js.Global().Call("setTimeout", fn, ms)
	}
}

// Throttle returns a throttled version of the callback that executes at most
// once per ms milliseconds.
//
// Example:
//
//	onScroll := reactive.Throttle(100, func() {
//	    // This runs at most every 100ms during scrolling
//	    updateScrollPosition()
//	})
func Throttle(ms int, callback func()) func() {
	var lastRun float64
	var scheduled bool

	return func() {
		now := js.Global().Get("Date").Call("now").Float()
		if now-lastRun >= float64(ms) {
			lastRun = now
			callback()
		} else if !scheduled {
			scheduled = true
			remaining := int(float64(ms) - (now - lastRun))
			SetTimeout(remaining, func() {
				scheduled = false
				lastRun = js.Global().Get("Date").Call("now").Float()
				callback()
			})
		}
	}
}
