//go:build js && wasm

package reactive

import (
	"syscall/js"
)

// Ticker holds a JavaScript interval that fires at regular intervals.
// Similar to time.Ticker but works with JavaScript's event loop.
//
// Example:
//
//	ticker := reactive.NewTicker(60000) // 60 seconds
//	go func() {
//	    for range ticker.C {
//	        // do something every 60 seconds
//	    }
//	}()
//	// Later, to stop:
//	ticker.Stop()
type Ticker struct {
	C          chan struct{}
	intervalID js.Value
	callback   js.Func
	stopped    bool
}

// NewTicker creates a new Ticker that sends to its channel every ms milliseconds.
// The ticker must be stopped with Stop() when no longer needed to prevent leaks.
func NewTicker(ms int) *Ticker {
	t := &Ticker{
		C: make(chan struct{}, 1), // buffered to prevent blocking JS
	}

	t.callback = js.FuncOf(func(this js.Value, args []js.Value) any {
		if !t.stopped {
			select {
			case t.C <- struct{}{}:
			default:
				// channel full, skip this tick
			}
		}
		return nil
	})

	t.intervalID = js.Global().Call("setInterval", t.callback, ms)
	return t
}

// Stop stops the ticker and releases resources.
// After Stop, no more values will be sent on the channel.
func (t *Ticker) Stop() {
	if t.stopped {
		return
	}
	t.stopped = true
	if !t.intervalID.IsUndefined() && !t.intervalID.IsNull() {
		js.Global().Call("clearInterval", t.intervalID)
	}
	t.callback.Release()
	close(t.C)
}

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
