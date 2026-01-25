package main

import "reactive"

// Components showcase - demonstrates component features
type Components struct {
	Message      *reactive.Store[string]
	ClickCount   *reactive.Store[int]
	CardTitle    *reactive.Store[string]
	AlertType    *reactive.Store[string]
	AlertMessage *reactive.Store[string]
}

func (c *Components) OnMount() {
	c.Message.Set("Hello from parent!")
	c.ClickCount.Set(0)
	c.CardTitle.Set("Dynamic Card")
	c.AlertType.Set("info")
	c.AlertMessage.Set("This is an alert message")
}

func (c *Components) HandleButtonClick() {
	c.ClickCount.Set(c.ClickCount.Get() + 1)
}

func (c *Components) SetAlertInfo() {
	c.AlertType.Set("info")
	c.AlertMessage.Set("This is an informational message")
}

func (c *Components) SetAlertSuccess() {
	c.AlertType.Set("success")
	c.AlertMessage.Set("Operation completed successfully!")
}

func (c *Components) SetAlertWarning() {
	c.AlertType.Set("warning")
	c.AlertMessage.Set("Please be careful with this action")
}

func (c *Components) SetAlertError() {
	c.AlertType.Set("error")
	c.AlertMessage.Set("Something went wrong!")
}

func (c *Components) Template() string {
	return `<div class="demo">
	<h1>Components</h1>

	<section>
		<h2>Basic Component with Props</h2>
		<p>Pass data to child components via props:</p>
		<Badge label="New" />
		<Badge label="Featured" />
		<Badge label="Sale" />
	</section>

	<section>
		<h2>Dynamic Props</h2>
		<p>Props can be bound to reactive stores:</p>
		<input type="text" bind:value="CardTitle" placeholder="Card title">
		<Card title="{CardTitle}">
			<p>This card's title updates as you type above.</p>
		</Card>
	</section>

	<section>
		<h2>Component with Slot</h2>
		<p>Components can accept child content via slots:</p>
		<Card title="Card with Slot">
			<p>This content is passed through the <strong>slot</strong>.</p>
			<p>You can put any HTML here!</p>
		</Card>
	</section>

	<section>
		<h2>Component Events</h2>
		<p>Child components can emit events to parent:</p>
		<p>Click count: <strong>{ClickCount}</strong></p>
		<Button label="Click Me" @click="HandleButtonClick()" />
		<Button label="Also Click Me" @click="HandleButtonClick()" />
	</section>

	<section>
		<h2>Conditional Styling Component</h2>
		<p>Components with dynamic classes based on props:</p>
		<div class="alert-buttons">
			<button @click="SetAlertInfo()">Info</button>
			<button @click="SetAlertSuccess()">Success</button>
			<button @click="SetAlertWarning()">Warning</button>
			<button @click="SetAlertError()">Error</button>
		</div>
		<Alert type="{AlertType}" message="{AlertMessage}" />
	</section>
</div>`
}

func (c *Components) Style() string {
	return `
.demo { max-width: 700px; }
.demo h1 { color: #1a1a2e; margin-bottom: 20px; }
.demo section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; background: #fff; }
.demo h2 { margin-top: 0; color: #666; font-size: 1.1em; }
.demo input[type="text"] { padding: 8px; width: 250px; border: 1px solid #ccc; border-radius: 4px; margin-bottom: 10px; }
.alert-buttons { display: flex; gap: 8px; margin-bottom: 10px; }
.alert-buttons button { padding: 6px 12px; border: 1px solid #ccc; border-radius: 4px; cursor: pointer; background: #f5f5f5; }
.alert-buttons button:hover { background: #e5e5e5; }
`
}

// Badge - simple component with a label prop
type Badge struct {
	Label *reactive.Store[string]
}

func (b *Badge) Template() string {
	return `<span class="badge">{Label}</span>`
}

func (b *Badge) Style() string {
	return `
.badge { display: inline-block; padding: 4px 8px; margin: 2px; background: #007bff; color: white; border-radius: 12px; font-size: 12px; font-weight: 500; }
`
}

// Card - component with title prop and slot for children
type Card struct {
	Title *reactive.Store[string]
}

func (c *Card) Template() string {
	return `<div class="card">
	<div class="card-header">{Title}</div>
	<div class="card-body"><slot/></div>
</div>`
}

func (c *Card) Style() string {
	return `
.card { border: 1px solid #ddd; border-radius: 8px; overflow: hidden; margin: 10px 0; }
.card-header { padding: 12px 16px; background: #f8f9fa; border-bottom: 1px solid #ddd; font-weight: 600; }
.card-body { padding: 16px; }
`
}

// Button - component that emits click events
type Button struct {
	Label *reactive.Store[string]
}

func (b *Button) Template() string {
	return `<button class="btn">{Label}</button>`
}

func (b *Button) Style() string {
	return `
.btn { padding: 10px 20px; margin: 4px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 14px; }
.btn:hover { background: #0056b3; }
`
}

// Alert - component with type-based styling
type Alert struct {
	Type    *reactive.Store[string]
	Message *reactive.Store[string]
}

func (a *Alert) Template() string {
	return `<div class="alert" data-type="{Type}">
	<strong class="alert-title">{Type}</strong>
	<span class="alert-message">{Message}</span>
</div>`
}

func (a *Alert) Style() string {
	return `
.alert { padding: 12px 16px; border-radius: 4px; margin: 10px 0; display: flex; align-items: center; gap: 10px; }
.alert-title { text-transform: uppercase; font-size: 12px; }
.alert-message { flex: 1; }
.alert[data-type="info"] { background: #e7f3ff; border: 1px solid #b3d7ff; color: #004085; }
.alert[data-type="success"] { background: #d4edda; border: 1px solid #c3e6cb; color: #155724; }
.alert[data-type="warning"] { background: #fff3cd; border: 1px solid #ffeeba; color: #856404; }
.alert[data-type="error"] { background: #f8d7da; border: 1px solid #f5c6cb; color: #721c24; }
`
}
