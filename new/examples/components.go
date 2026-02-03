package main

import p "preveltekit"

// Components showcase - demonstrates component features
type Components struct {
	Message      *p.Store[string]
	ClickCount   *p.Store[int]
	CardTitle    *p.Store[string]
	AlertType    *p.Store[string]
	AlertMessage *p.Store[string]
}

func (c *Components) New() p.Component {
	return &Components{
		Message:      p.New("components.Message", "Hello from parent!"),
		ClickCount:   p.New("components.ClickCount", 0),
		CardTitle:    p.New("components.CardTitle", "Dynamic Card"),
		AlertType:    p.New("components.AlertType", "info"),
		AlertMessage: p.New("components.AlertMessage", "This is an alert message"),
	}
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

func (c *Components) Render() p.Node {
	return p.Html(`<div class="demo">
		<h1>Components</h1>

		<section>
			<h2>Basic Component with Props</h2>
			<p>Pass data to child components via props:</p>`,
		p.Comp(&Badge{Label: p.New("badge.New", "New")}),
		p.Comp(&Badge{Label: p.New("badge.Featured", "Featured")}),
		p.Comp(&Badge{Label: p.New("badge.Sale", "Sale")}),
		`</section>

		<section>
			<h2>Dynamic Props</h2>
			<p>Props can be bound to reactive stores:</p>
			`, p.BindValue(`<input type="text" placeholder="Card title">`, c.CardTitle),
		p.Comp(&Card{Title: c.CardTitle},
			p.Html(`<p>This card's title updates as you type above.</p>`),
		),
		`</section>

		<section>
			<h2>Component with Slot</h2>
			<p>Components can accept child content via slots:</p>`,
		p.Comp(&Card{Title: p.New("card.Slot", "Card with Slot")},
			p.Html(`<p>This content is passed through the <strong>slot</strong>.</p>
				<p>You can put any HTML here!</p>`),
		),
		`</section>

		<section>
			<h2>Component Events</h2>
			<p>Child components can emit events to parent:</p>
			<p>Click count: <strong>`, p.Bind(c.ClickCount), `</strong></p>`,
		p.Comp(&Button{Label: p.New("button.ClickMe", "Click Me"), OnClick: c.HandleButtonClick}),
		p.Comp(&Button{Label: p.New("button.AlsoClickMe", "Also Click Me"), OnClick: c.HandleButtonClick}),
		`</section>

		<section>
			<h2>Conditional Styling Component</h2>
			<p>Components with dynamic classes based on props:</p>
			<div class="alert-buttons">
				`, p.Html(`<button>Info</button>`).WithOn("click", "components.SetAlertInfo", c.SetAlertInfo), `
				`, p.Html(`<button>Success</button>`).WithOn("click", "components.SetAlertSuccess", c.SetAlertSuccess), `
				`, p.Html(`<button>Warning</button>`).WithOn("click", "components.SetAlertWarning", c.SetAlertWarning), `
				`, p.Html(`<button>Error</button>`).WithOn("click", "components.SetAlertError", c.SetAlertError), `
			</div>`,
		p.Comp(&Alert{Type: c.AlertType, Message: c.AlertMessage}),
		`</section>
	</div>`)
}

func (c *Components) Style() string {
	// Nested component styles (Badge, Card, Button, Alert) are auto-collected during SSR
	return `
.demo{max-width:700px}
.alert-buttons{display:flex;gap:8px;margin-bottom:10px}
`
}

// Badge - simple component with a label prop
type Badge struct {
	Label *p.Store[string]
}

func (b *Badge) Render() p.Node {
	return p.Html(`<span class="badge">`, p.Bind(b.Label), `</span>`)
}

func (b *Badge) Style() string {
	return `.badge{display:inline-block;padding:4px 8px;margin:2px;background:#007bff;color:#fff;border-radius:12px;font-size:12px;font-weight:500}`
}

// Card - component with title prop and slot for children
type Card struct {
	Title *p.Store[string]
}

func (c *Card) Render() p.Node {
	return p.Html(`<div class="card">
		<div class="card-header">`, p.Bind(c.Title), `</div>
		<div class="card-body">`, p.Slot(), `</div>
	</div>`)
}

func (c *Card) Style() string {
	return `.card{border:1px solid #ddd;border-radius:8px;overflow:hidden;margin:10px 0}.card-header{padding:12px 16px;background:#f8f9fa;border-bottom:1px solid #ddd;font-weight:600}.card-body{padding:16px}`
}

// Button - component that emits click events
type Button struct {
	Label   *p.Store[string]
	OnClick func()
}

func (b *Button) Render() p.Node {
	if b.OnClick != nil {
		// Use the label's ID as part of the handler ID for uniqueness
		handlerID := b.Label.ID() + ".click"
		return p.Html(`<button class="btn">`, p.Bind(b.Label), `</button>`).WithOn("click", handlerID, b.OnClick)
	}
	return p.Html(`<button class="btn">`, p.Bind(b.Label), `</button>`)
}

func (b *Button) Style() string {
	return `.btn{padding:10px 20px;margin:4px;background:#007bff;color:#fff;border:none;border-radius:4px;cursor:pointer}.btn:hover{background:#0056b3}`
}

// Alert - component with type-based styling
type Alert struct {
	Type    *p.Store[string]
	Message *p.Store[string]
}

func (a *Alert) Render() p.Node {
	return p.Html(`<div class="alert" `, p.DynAttr("data-type", "{0}", a.Type), `>
		<strong class="alert-title">`, p.Bind(a.Type), `</strong>
		<span class="alert-message">`, p.Bind(a.Message), `</span>
	</div>`)
}

func (a *Alert) Style() string {
	return `.alert{padding:12px 16px;border-radius:4px;margin:10px 0;display:flex;align-items:center;gap:10px}.alert-title{text-transform:uppercase;font-size:12px}.alert-message{flex:1}.alert[data-type=info]{background:#e7f3ff;border:1px solid #b3d7ff;color:#004085}.alert[data-type=success]{background:#d4edda;border:1px solid #c3e6cb;color:#155724}.alert[data-type=warning]{background:#fff3cd;border:1px solid #ffeeba;color:#856404}.alert[data-type=error]{background:#f8d7da;border:1px solid #f5c6cb;color:#721c24}`
}
