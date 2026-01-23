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

func (a *App) Template() string {
	return `<div class="app">
	<h1>Component Composition</h1>

	<section>
		<h2>1. Multiple Instances with Different Props</h2>
		<p>Same Button component, different configurations:</p>
		<div class="button-row">
			<Button label="Primary" variant="primary" @click="Increment()" />
			<Button label="Secondary" variant="secondary" @click="Decrement()" />
			<Button label="Success" variant="success" @click="Add(5)" />
			<Button label="Danger" variant="danger" @click="Reset()" />
		</div>
		<p>Count: <strong>{Count}</strong></p>
	</section>

	<section>
		<h2>2. Dynamic Props</h2>
		<p>Button label updates when parent state changes:</p>
		<Button label="{Message}" variant="primary" @click="Increment()" />
	</section>

	<section>
		<h2>3. Slots with Parent-Bound Content</h2>
		<p>Content inside component tags becomes slot content:</p>
		<Button variant="primary" @click="Increment()">
			Count is {Count}
		</Button>
		<Card title="Status Card" @click="CardClicked()">
			<p>Current count: <strong>{Count}</strong></p>
			<p>Theme: <strong>{Theme}</strong></p>
		</Card>
	</section>

	<section>
		<h2>4. Component with Internal State</h2>
		<p>Counter has its own internal state:</p>
		<Counter initial="10" step="5" />
	</section>

	<section>
		<h2>5. Conditional Rendering</h2>
		<Button label="Toggle Card" variant="secondary" @click="ToggleCard()" />
		{#if ShowCard}
			<Card title="Toggleable Card" @click="CardClicked()">
				<p>This card can be shown/hidden.</p>
				<p>Parent count: {Count}</p>
			</Card>
		{:else}
			<p class="hidden-note">Card is hidden</p>
		{/if}
	</section>

	<section>
		<h2>Event Log</h2>
		<p class="message">{Message}</p>
	</section>
</div>`
}

func (a *App) Style() string {
	return `
.app { font-family: system-ui, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; }
h1 { color: #333; }
h2 { margin-top: 0; color: #666; font-size: 1.1em; }
.button-row { display: flex; gap: 10px; flex-wrap: wrap; margin: 10px 0; }
.message { padding: 10px; background: #e7f3ff; border-radius: 4px; }
.hidden-note { color: #999; font-style: italic; }
`
}
