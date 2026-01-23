package main

import "reactive"

// App demonstrates component composition patterns
type App struct {
	// State
	Count    *reactive.Store[int]
	Message  *reactive.Store[string]
	Theme    *reactive.Store[string]
	ShowCard *reactive.Store[bool]
}

func (a *App) OnMount() {
	a.Count.Set(0)
	a.Message.Set("Click a button to see events in action")
	a.Theme.Set("light")
	a.ShowCard.Set(true)
}

// Event handlers
func (a *App) Increment() {
	a.Count.Set(a.Count.Get() + 1)
	a.Message.Set("Incremented!")
}

func (a *App) Decrement() {
	a.Count.Set(a.Count.Get() - 1)
	a.Message.Set("Decremented!")
}

func (a *App) Add(n int) {
	a.Count.Set(a.Count.Get() + n)
	a.Message.Set("Added " + string(rune('0'+n)) + "!")
}

func (a *App) Reset() {
	a.Count.Set(0)
	a.Message.Set("Reset to zero!")
}

func (a *App) ToggleTheme() {
	if a.Theme.Get() == "light" {
		a.Theme.Set("dark")
	} else {
		a.Theme.Set("light")
	}
}

func (a *App) ToggleCard() {
	a.ShowCard.Set(!a.ShowCard.Get())
}

func (a *App) CardClicked() {
	a.Message.Set("Card was clicked!")
}

