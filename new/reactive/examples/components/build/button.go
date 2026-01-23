package main

import "reactive"

// Button is a reusable component with props
type Button struct {
	// Props (set by parent)
	Label   *reactive.Store[string]
	Variant *reactive.Store[string] // "primary" or "secondary"
}

