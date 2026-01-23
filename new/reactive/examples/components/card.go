package main

import "reactive"

// Card is a container component with a title and slot content
type Card struct {
	Title *reactive.Store[string]
}

func (c *Card) Template() string {
	return `<div class="card">
	<div class="card-header">{Title}</div>
	<div class="card-body"><slot></slot></div>
</div>`
}

func (c *Card) Style() string {
	return `
.card { border: 1px solid #ddd; border-radius: 8px; margin: 10px 0; overflow: hidden; }
.card-header { background: #f5f5f5; padding: 10px 15px; font-weight: bold; border-bottom: 1px solid #ddd; }
.card-body { padding: 15px; }
`
}
