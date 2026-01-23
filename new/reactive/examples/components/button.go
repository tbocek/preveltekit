package main

import "reactive"

// Button is a reusable component with props
type Button struct {
	// Props (set by parent)
	Label   *reactive.Store[string]
	Variant *reactive.Store[string] // "primary" or "secondary"
}

func (b *Button) Template() string {
	return `<button class="btn {Variant}">
		<slot/>
		{Label}
	</button>`
}

func (b *Button) Style() string {
	return `
		.btn { padding: 0.5em 1em; cursor: pointer; border: none; border-radius: 4px; }
		.btn.primary { background: #007bff; color: white; }
		.btn.secondary { background: #6c757d; color: white; }
	`
}
