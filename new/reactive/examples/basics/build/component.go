package main

import "reactive"

type Basics struct {
	// Numbers
	Count *reactive.Store[int]

	// Text
	Name    *reactive.Store[string]
	Message *reactive.Store[string]

	// Boolean
	DarkMode *reactive.Store[bool]
	Agreed   *reactive.Store[bool]

	// For conditional demo
	Score *reactive.Store[int]
}

func (b *Basics) OnMount() {
	b.Count.Set(0)
	b.Name.Set("")
	b.Message.Set("Fill out the form above")
	b.DarkMode.Set(false)
	b.Agreed.Set(false)
	b.Score.Set(75)
}

// Count operations
func (b *Basics) Increment() {
	b.Count.Update(func(v int) int { return v + 1 })
}

func (b *Basics) Decrement() {
	b.Count.Update(func(v int) int { return v - 1 })
}

func (b *Basics) Add(n int) {
	b.Count.Update(func(v int) int { return v + n })
}

func (b *Basics) Reset() {
	b.Count.Set(0)
}

// Score operations for conditional demo
func (b *Basics) SetScore(n int) {
	b.Score.Set(n)
}

// Form submission
func (b *Basics) Submit() {
	name := b.Name.Get()
	if name == "" {
		b.Message.Set("Please enter your name")
	} else if !b.Agreed.Get() {
		b.Message.Set("Please agree to terms")
	} else {
		b.Message.Set("Welcome, " + name + "!")
	}
}

