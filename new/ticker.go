//go:build wasm

package preveltekit

import (
	"syscall/js"
	"time"
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
// Returns the debounced function and a cleanup function to release resources.
//
// Example:
//
//	search, cleanup := reactive.Debounce(300, func() {
//	    // This runs 300ms after the last keystroke
//	    performSearch(input.Get())
//	})
//	input.OnChange(func(_ string) { search() })
//	// When done:
//	cleanup()
func Debounce(ms int, callback func()) (func(), func()) {
	var timeoutID js.Value
	var released bool

	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		callback()
		return nil
	})

	debounced := func() {
		if released {
			return // Don't schedule if already cleaned up
		}
		if !timeoutID.IsUndefined() && !timeoutID.IsNull() {
			js.Global().Call("clearTimeout", timeoutID)
		}
		timeoutID = js.Global().Call("setTimeout", fn, ms)
	}

	cleanup := func() {
		if released {
			return // Already cleaned up, prevent double-release panic
		}
		released = true
		if !timeoutID.IsUndefined() && !timeoutID.IsNull() {
			js.Global().Call("clearTimeout", timeoutID)
		}
		fn.Release()
	}

	return debounced, cleanup
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
	var lastRun int64
	var scheduled bool
	msInt64 := int64(ms)

	return func() {
		now := time.Now().UnixMilli()
		if now-lastRun >= msInt64 {
			lastRun = now
			callback()
		} else if !scheduled {
			scheduled = true
			remaining := int(msInt64 - (now - lastRun))
			SetTimeout(remaining, func() {
				scheduled = false
				lastRun = time.Now().UnixMilli()
				callback()
			})
		}
	}
}
