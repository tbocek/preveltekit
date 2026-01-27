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

func (c *Components) Render() p.Node {
	return p.Div(p.Class("demo"),
		p.H1("Components"),

		// Basic Component with Props
		p.Section(
			p.H2("Basic Component with Props"),
			p.P("Pass data to child components via props:"),
			p.Comp("Badge", p.Prop("label", "New")),
			p.Comp("Badge", p.Prop("label", "Featured")),
			p.Comp("Badge", p.Prop("label", "Sale")),
		),

		// Dynamic Props
		p.Section(
			p.H2("Dynamic Props"),
			p.P("Props can be bound to reactive stores:"),
			p.Input(p.Type("text"), p.BindValue(c.CardTitle), p.Placeholder("Card title")),
			p.Comp("Card", p.Prop("title", c.CardTitle),
				p.P("This card's title updates as you type above."),
			),
		),

		// Component with Slot
		p.Section(
			p.H2("Component with Slot"),
			p.P("Components can accept child content via slots:"),
			p.Comp("Card", p.Prop("title", "Card with Slot"),
				p.P("This content is passed through the ", p.Strong("slot"), "."),
				p.P("You can put any HTML here!"),
			),
		),

		// Component Events
		p.Section(
			p.H2("Component Events"),
			p.P("Child components can emit events to parent:"),
			p.P("Click count: ", p.Strong(p.Bind(c.ClickCount))),
			p.Comp("Button", p.Prop("label", "Click Me"), p.OnClick(c.HandleButtonClick)),
			p.Comp("Button", p.Prop("label", "Also Click Me"), p.OnClick(c.HandleButtonClick)),
		),

		// Conditional Styling Component
		p.Section(
			p.H2("Conditional Styling Component"),
			p.P("Components with dynamic classes based on props:"),
			p.Div(p.Class("alert-buttons"),
				p.Button("Info", p.OnClick(c.SetAlertInfo)),
				p.Button("Success", p.OnClick(c.SetAlertSuccess)),
				p.Button("Warning", p.OnClick(c.SetAlertWarning)),
				p.Button("Error", p.OnClick(c.SetAlertError)),
			),
			p.Comp("Alert", p.Prop("type", c.AlertType), p.Prop("message", c.AlertMessage)),
		),
	)
}

func (c *Components) Style() string {
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
	return p.Span(p.Class("badge"), p.Bind(b.Label))
}

func (b *Badge) Style() string {
	return `.badge{display:inline-block;padding:4px 8px;margin:2px;background:#007bff;color:#fff;border-radius:12px;font-size:12px;font-weight:500}`
}

// Card - component with title prop and slot for children
type Card struct {
	Title *p.Store[string]
}

func (c *Card) Render() p.Node {
	return p.Div(p.Class("card"),
		p.Div(p.Class("card-header"), p.Bind(c.Title)),
		p.Div(p.Class("card-body"), p.Slot()),
	)
}

func (c *Card) Style() string {
	return `.card{border:1px solid #ddd;border-radius:8px;overflow:hidden;margin:10px 0}.card-header{padding:12px 16px;background:#f8f9fa;border-bottom:1px solid #ddd;font-weight:600}.card-body{padding:16px}`
}

// Button - component that emits click events
type Button struct {
	Label *p.Store[string]
}

func (b *Button) Render() p.Node {
	return p.Button(p.Class("btn"), p.Bind(b.Label))
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
	return p.Div(p.Class("alert"), p.DynAttr("data-type", "{0}", a.Type),
		p.Strong(p.Class("alert-title"), p.Bind(a.Type)),
		p.Span(p.Class("alert-message"), p.Bind(a.Message)),
	)
}

func (a *Alert) Style() string {
	return `.alert{padding:12px 16px;border-radius:4px;margin:10px 0;display:flex;align-items:center;gap:10px}.alert-title{text-transform:uppercase;font-size:12px}.alert-message{flex:1}.alert[data-type=info]{background:#e7f3ff;border:1px solid #b3d7ff;color:#004085}.alert[data-type=success]{background:#d4edda;border:1px solid #c3e6cb;color:#155724}.alert[data-type=warning]{background:#fff3cd;border:1px solid #ffeeba;color:#856404}.alert[data-type=error]{background:#f8d7da;border:1px solid #f5c6cb;color:#721c24}`
}

func (c *Components) HandleEvent(method string, args string) {
	switch method {
	case "HandleButtonClick":
		c.HandleButtonClick()
	case "SetAlertInfo":
		c.SetAlertInfo()
	case "SetAlertSuccess":
		c.SetAlertSuccess()
	case "SetAlertWarning":
		c.SetAlertWarning()
	case "SetAlertError":
		c.SetAlertError()
	}
}

func (b *Badge) HandleEvent(method string, args string)  {}
func (c *Card) HandleEvent(method string, args string)   {}
func (b *Button) HandleEvent(method string, args string) {}
func (a *Alert) HandleEvent(method string, args string)  {}
