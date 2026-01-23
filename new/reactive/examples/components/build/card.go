package main

import "reactive"

// Card is a container component with a title and slot content
type Card struct {
	Title *reactive.Store[string]
}

