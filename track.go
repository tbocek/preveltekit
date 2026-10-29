package preveltekit

import "slices"

// Automatic dependency tracking: while a function runs under
// dependencies.track, every Store.Get (and List.Get/At) records its store.
// Conditions (If, AttrIf) and Computed use this to subscribe to exactly the
// stores they read, so they need no list of dependencies.

// subscriber is implemented by Store and List: subscribe without knowing the
// value type, with a function to unsubscribe again.
type subscriber interface {
	subscribeAny(func()) (unsubscribe func())
}

func (s *Store[T]) subscribeAny(fn func()) func() {
	return s.subscribe(func(T) { fn() })
}

// reads collects the stores read by the function being tracked (nil: none).
var reads *[]subscriber

// read records s if a function is being tracked.
func read(s subscriber) {
	if reads != nil && !slices.Contains(*reads, s) {
		*reads = append(*reads, s)
	}
}

// dependencies keeps a callback subscribed to the stores a function read the
// last time it ran.
type dependencies struct {
	stores []subscriber
	unsubs []func()
}

// track runs fn, records the stores it reads with Get and subscribes
// onChange to exactly those. The stores are recorded on every run, so
// dependencies that change with the data (a && b only reads b once a is
// true) are followed; the subscriptions only change if the stores do.
func (d *dependencies) track(fn func(), onChange func()) {
	prev := reads
	var stores []subscriber
	reads = &stores
	fn()
	reads = prev

	if slices.Equal(stores, d.stores) {
		return
	}
	d.stop()
	d.stores = stores
	for _, s := range stores {
		d.unsubs = append(d.unsubs, s.subscribeAny(onChange))
	}
}

// stop removes all subscriptions.
func (d *dependencies) stop() {
	for _, unsub := range d.unsubs {
		unsub()
	}
	d.stores, d.unsubs = nil, nil
}

// Computed returns a store holding fn's result, recomputed whenever a store
// that fn reads with Get changes.
//
//	fullName := p.Computed(func() string { return first.Get() + " " + last.Get() })
func Computed[T any](fn func() T) *Store[T] {
	var out *Store[T]
	deps := &dependencies{}
	var run func()
	run = func() {
		var v T
		deps.track(func() { v = fn() }, run)
		if out == nil {
			out = New(v)
		} else {
			out.Set(v)
		}
	}
	run()
	return out
}
