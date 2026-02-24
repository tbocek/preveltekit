package main

import p "github.com/tbocek/preveltekit/v2"

// Components showcase - demonstrates component features
type Components struct {
	Message      *p.Store[string]
	ClickCount   *p.Store[int]
	CardTitle    *p.Store[string]
	AlertType    *p.Store[string]
	AlertMessage *p.Store[string]
	ViewMode     *p.Store[string]
}

func (c *Components) New() p.Component {
	return &Components{
		Message:      p.New("Hello from parent!"),
		ClickCount:   p.New(0),
		CardTitle:    p.New("Dynamic Card"),
		AlertType:    p.New("info"),
		AlertMessage: p.New("This is an alert message"),
		ViewMode:     p.New("card"),
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

func (c *Components) SetViewMode(mode string) {
	c.ViewMode.Set(mode)
}

func (c *Components) Render() p.Node {
	return p.Div(p.Attr("class", "demo"),
		p.H1("Components"),

		p.Section(
			p.H2("Basic Component with Props"),
			p.P("Pass data to child components via struct fields:"),
			p.Comp(&Badge{Label: p.New("New")}),
			p.Comp(&Badge{Label: p.New("Featured")}),
			p.Comp(&Badge{Label: p.New("Sale")}),
			p.Pre(p.Attr("class", "code"), `type Badge struct {
    Label *p.Store[string]
}

func (b *Badge) Render() p.Node {
    return p.Span(p.Attr("class", "badge"), b.Label)
}

// usage:
p.Comp(&Badge{Label: p.New("New")})`),
		),

		p.Section(
			p.H2("Dynamic Props"),
			p.P("Props can be bound to reactive stores:"),
			p.Input(p.Attr("type", "text"), p.Attr("placeholder", "Card title")).Bind(c.CardTitle),
			p.Comp(&Card{Title: c.CardTitle},
				p.P("This card's title updates as you type above."),
			),
		),

		p.Section(
			p.H2("Component with Slot"),
			p.P("Components can accept child content via slots:"),
			p.Comp(&Card{Title: p.New("Card with Slot")},
				p.P("This content is passed through the ", p.Strong("slot"), "."),
				p.P("You can put any HTML here!"),
			),
			p.Pre(p.Attr("class", "code"), `func (c *Card) Render() p.Node {
    return p.Div(p.Attr("class", "card"),
        p.Div(p.Attr("class", "card-header"), c.Title),
        p.Div(p.Attr("class", "card-body"), p.Slot()),
    )
}

// usage — child content fills the Slot():
p.Comp(&Card{Title: p.New("Title")},
    p.P("Slot content here"),
)`),
		),

		p.Section(
			p.H2("Component Events"),
			p.P("Child components can emit events to parent:"),
			p.P("Click count: ", p.Strong(c.ClickCount)),
			p.Comp(&Button{Label: p.New("Click Me"), OnClick: c.HandleButtonClick}),
			p.Comp(&Button{Label: p.New("Also Click Me"), OnClick: c.HandleButtonClick}),
			p.Pre(p.Attr("class", "code"), `type Button struct {
    Label   *p.Store[string]
    OnClick func() // callback prop — parent passes handler
}

// usage:
p.Comp(&Button{Label: p.New("Click"), OnClick: handler})`),
		),

		p.Section(
			p.H2("Conditional Styling Component"),
			p.P("Components with dynamic classes based on props:"),
			p.Div(p.Attr("class", "alert-buttons"),
				p.Button("Info").On("click", c.SetAlertInfo),
				p.Button("Success").On("click", c.SetAlertSuccess),
				p.Button("Warning").On("click", c.SetAlertWarning),
				p.Button("Error").On("click", c.SetAlertError),
			),
			p.Comp(&Alert{Type: c.AlertType, Message: c.AlertMessage}),
			p.Pre(p.Attr("class", "code"), "// Attr() sets dynamic attributes from stores\n"+
				"p.Div(p.Attr(\"class\", \"alert\"), p.Attr(\"data-type\", alertType))\n\n"+
				"// components get scoped CSS via Style()\n"+
				"func (a *Alert) Style() string { return `.alert{...}` }"),
		),

		p.Section(
			p.H2("Conditional Components"),
			p.P("Components with slots and props inside if-blocks:"),
			p.P("Current view: ", p.Strong(c.ViewMode)),
			p.Div(p.Attr("class", "view-buttons"),
				p.Button("Card").On("click", func() { c.SetViewMode("card") }),
				p.Button("Badge").On("click", func() { c.SetViewMode("badge") }),
				p.Button("Alert").On("click", func() { c.SetViewMode("alert") }),
			),
			p.If(p.Cond(func() bool { return c.ViewMode.Get() == "card" }, c.ViewMode),
				p.Comp(&Card{Title: c.CardTitle},
					p.P("This card appears conditionally."),
					p.P("It receives a ", p.Strong("dynamic prop"), " and ", p.Strong("slot content"), "."),
				),
			).ElseIf(p.Cond(func() bool { return c.ViewMode.Get() == "badge" }, c.ViewMode),
				p.Comp(&Badge{Label: c.Message}),
			).Else(
				p.Comp(&Alert{Type: p.New("success"), Message: c.Message}),
			),
		),
	)
}

func (c *Components) Style() string {
	return `
.demo{max-width:700px}
.demo pre.code{background:#1a1a2e;color:#e0e0e0;font-size:12px;margin-top:12px}
.alert-buttons,.view-buttons{display:flex;gap:8px;margin-bottom:10px}
`
}

// Badge - simple component with a label prop
type Badge struct {
	Label *p.Store[string]
}

func (b *Badge) Render() p.Node {
	return p.Span(p.Attr("class", "badge"), b.Label)
}

func (b *Badge) Style() string {
	return `.badge{display:inline-block;padding:4px 8px;margin:2px;background:#007bff;color:#fff;border-radius:12px;font-size:12px;font-weight:500}`
}

// Card - component with title prop and slot for children
type Card struct {
	Title *p.Store[string]
}

func (c *Card) Render() p.Node {
	return p.Div(p.Attr("class", "card"),
		p.Div(p.Attr("class", "card-header"), c.Title),
		p.Div(p.Attr("class", "card-body"), p.Slot()),
	)
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
		return p.Button(p.Attr("class", "btn"), b.Label).On("click", b.OnClick)
	}
	return p.Button(p.Attr("class", "btn"), b.Label)
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
	return p.Div(p.Attr("class", "alert"), p.Attr("data-type", a.Type),
		p.Strong(p.Attr("class", "alert-title"), a.Type),
		p.Span(p.Attr("class", "alert-message"), a.Message),
	)
}

func (a *Alert) Style() string {
	return `.alert{padding:12px 16px;border-radius:4px;margin:10px 0;display:flex;align-items:center;gap:10px}.alert-title{text-transform:uppercase;font-size:12px}.alert-message{flex:1}.alert[data-type=info]{background:#e7f3ff;border:1px solid #b3d7ff;color:#004085}.alert[data-type=success]{background:#d4edda;border:1px solid #c3e6cb;color:#155724}.alert[data-type=warning]{background:#fff3cd;border:1px solid #ffeeba;color:#856404}.alert[data-type=error]{background:#f8d7da;border:1px solid #f5c6cb;color:#721c24}`
}
