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

func (c *Counter) Template() string {
	return `<div class="counter">
	<span class="value">{Value}</span>
	<button class="counter-btn" @click="Dec()">-{Step}</button>
	<button class="counter-btn" @click="Inc()">+{Step}</button>
	<button class="counter-btn reset" @click="Reset()">Reset</button>
</div>`
}

func (c *Counter) Style() string {
	return `
.counter { display: inline-flex; align-items: center; gap: 8px; padding: 10px; background: #f0f0f0; border-radius: 8px; margin: 10px 0; }
.counter .value { font-size: 24px; font-weight: bold; min-width: 60px; text-align: center; }
.counter-btn { padding: 8px 12px; border: none; border-radius: 4px; cursor: pointer; background: #007bff; color: white; }
.counter-btn:hover { background: #0056b3; }
.counter-btn.reset { background: #6c757d; }
.counter-btn.reset:hover { background: #545b62; }
`
}
