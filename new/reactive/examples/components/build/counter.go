package main

import "reactive"

// Counter is a component with its own internal state
type Counter struct {
	// Props (set by parent)
	Initial *reactive.Store[int]
	Step    *reactive.Store[int]

	// Internal state
	Value *reactive.Store[int]
}

func (c *Counter) OnMount() {
	// Initialize value from initial prop
	c.Value.Set(c.Initial.Get())
}

func (c *Counter) Inc() {
	c.Value.Set(c.Value.Get() + c.Step.Get())
}

func (c *Counter) Dec() {
	c.Value.Set(c.Value.Get() - c.Step.Get())
}

func (c *Counter) Reset() {
	c.Value.Set(c.Initial.Get())
}

