package main

import p "github.com/tbocek/preveltekit/v2"

// Routing showcase - demonstrates navigation patterns
type Routing struct {
	CurrentTab  *p.Store[string]
	CurrentStep *p.Store[int]
}

func (r *Routing) New() p.Component {
	return &Routing{
		CurrentTab:  p.New("home"),
		CurrentStep: p.New(1),
	}
}

func (r *Routing) GoToTab(tab string) {
	r.CurrentTab.Set(tab)
}

func (r *Routing) GoHome() {
	r.GoToTab("home")
}

func (r *Routing) GoProfile() {
	r.GoToTab("profile")
}

func (r *Routing) GoSettings() {
	r.GoToTab("settings")
}

func (r *Routing) GoNotifications() {
	r.GoToTab("notifications")
}

func (r *Routing) NextStep() {
	if r.CurrentStep.Get() < 4 {
		r.CurrentStep.Set(r.CurrentStep.Get() + 1)
	}
}

func (r *Routing) PrevStep() {
	if r.CurrentStep.Get() > 1 {
		r.CurrentStep.Set(r.CurrentStep.Get() - 1)
	}
}

func (r *Routing) GoToStep(step int) {
	if step >= 1 && step <= 4 {
		r.CurrentStep.Set(step)
	}
}

func (r *Routing) Step1() {
	r.GoToStep(1)
}

func (r *Routing) Step2() {
	r.GoToStep(2)
}

func (r *Routing) Step3() {
	r.GoToStep(3)
}

func (r *Routing) Step4() {
	r.GoToStep(4)
}

func (r *Routing) Render() p.Node {
	return p.Div(p.Attr("class", "demo"),
		p.H1("Routing & Navigation"),

		p.Section(
			p.H2("Tab Navigation"),
			p.P("Simple tab-based navigation pattern:"),

			p.Div(p.Attr("class", "tabs"),
				p.Button("Home").AttrIf("class", func() bool { return r.CurrentTab.Get() == "home" }, "active").On("click", func() { r.GoHome() }),
				p.Button("Profile").AttrIf("class", func() bool { return r.CurrentTab.Get() == "profile" }, "active").On("click", func() { r.GoProfile() }),
				p.Button("Settings").AttrIf("class", func() bool { return r.CurrentTab.Get() == "settings" }, "active").On("click", func() { r.GoSettings() }),
				p.Button("Notifications").AttrIf("class", func() bool { return r.CurrentTab.Get() == "notifications" }, "active").On("click", func() { r.GoNotifications() }),
			),

			p.Div(p.Attr("class", "tab-content"),
				p.If(func() bool { return r.CurrentTab.Get() == "home" },
					p.Div(p.Attr("class", "tab-panel"), p.H3("Home"), p.P("Welcome to the home tab! This is the default view.")),
				).ElseIf(func() bool { return r.CurrentTab.Get() == "profile" },
					p.Div(p.Attr("class", "tab-panel"), p.H3("Profile"), p.P("View and edit your profile information here.")),
				).ElseIf(func() bool { return r.CurrentTab.Get() == "settings" },
					p.Div(p.Attr("class", "tab-panel"), p.H3("Settings"), p.P("Configure your application settings.")),
				).ElseIf(func() bool { return r.CurrentTab.Get() == "notifications" },
					p.Div(p.Attr("class", "tab-panel"), p.H3("Notifications"), p.P("View your recent notifications.")),
				),
			),
		),

		p.Section(
			p.H2("Wizard / Stepper"),
			p.P("Multi-step form navigation:"),

			p.Div(p.Attr("class", "stepper"),
				p.Div(p.Attr("class", "step"),
					p.Button("1"),
					p.Span("Details"),
				).
					AttrIf("class", func() bool { return r.CurrentStep.Get() > 1 }, "completed").
					AttrIf("class", func() bool { return r.CurrentStep.Get() == 1 }, "active").
					On("click", func() { r.Step1() }),
				p.Div(p.Attr("class", "step-line")).AttrIf("class", func() bool { return r.CurrentStep.Get() > 1 }, "completed"),
				p.Div(p.Attr("class", "step"),
					p.Button("2"),
					p.Span("Address"),
				).
					AttrIf("class", func() bool { return r.CurrentStep.Get() > 2 }, "completed").
					AttrIf("class", func() bool { return r.CurrentStep.Get() == 2 }, "active").
					On("click", func() { r.Step2() }),
				p.Div(p.Attr("class", "step-line")).AttrIf("class", func() bool { return r.CurrentStep.Get() > 2 }, "completed"),
				p.Div(p.Attr("class", "step"),
					p.Button("3"),
					p.Span("Payment"),
				).
					AttrIf("class", func() bool { return r.CurrentStep.Get() > 3 }, "completed").
					AttrIf("class", func() bool { return r.CurrentStep.Get() == 3 }, "active").
					On("click", func() { r.Step3() }),
				p.Div(p.Attr("class", "step-line")).AttrIf("class", func() bool { return r.CurrentStep.Get() > 3 }, "completed"),
				p.Div(p.Attr("class", "step"),
					p.Button("4"),
					p.Span("Confirm"),
				).
					AttrIf("class", func() bool { return r.CurrentStep.Get() == 4 }, "active").
					On("click", func() { r.Step4() }),
			),

			p.Div(p.Attr("class", "step-content"),
				p.If(func() bool { return r.CurrentStep.Get() == 1 },
					p.Div(p.Attr("class", "step-panel"), p.H3("Step 1: Personal Details"), p.P("Enter your name and email address.")),
				).ElseIf(func() bool { return r.CurrentStep.Get() == 2 },
					p.Div(p.Attr("class", "step-panel"), p.H3("Step 2: Shipping Address"), p.P("Enter your shipping address.")),
				).ElseIf(func() bool { return r.CurrentStep.Get() == 3 },
					p.Div(p.Attr("class", "step-panel"), p.H3("Step 3: Payment Method"), p.P("Choose your payment method.")),
				).ElseIf(func() bool { return r.CurrentStep.Get() == 4 },
					p.Div(p.Attr("class", "step-panel"), p.H3("Step 4: Confirmation"), p.P("Review and confirm your order.")),
				),
			),

			p.Div(p.Attr("class", "step-buttons"),
				p.Button("Previous").
					AttrIf("class", func() bool { return r.CurrentStep.Get() == 1 }, "disabled").
					On("click", func() { r.PrevStep() }),
				p.Span("Step ", r.CurrentStep, " of 4"),
				p.Button("Next").
					AttrIf("class", func() bool { return r.CurrentStep.Get() == 4 }, "disabled").
					On("click", func() { r.NextStep() }),
			),

			p.Pre(p.Attr("class", "code"), `// AttrIf: conditionally add classes (additive for same attribute)
p.Div(p.Attr("class", "step"), ...).
    AttrIf("class", func() bool { return Step.Get() > 2 }, "completed").
    AttrIf("class", func() bool { return Step.Get() == 2 }, "active").
    On("click", GoToStep2)

// chain multiple AttrIf + On calls on same node`),
		),

		p.Section(
			p.H2("Route Parameters"),
			p.P("The route ", p.Code("/user/:name"), " passes the name from the URL to the page:"),
			p.P(
				p.A(p.Attr("href", "/user/alice"), "Alice"), " · ",
				p.A(p.Attr("href", "/user/bob"), "Bob"),
			),
		),
	)
}

func (r *Routing) Style() string {
	return `
.demo{max-width:700px}
.demo button.disabled{opacity:.5;cursor:not-allowed}
.tabs{display:flex;gap:4px;border-bottom:2px solid #ddd}
.tabs button{border:none;border-radius:4px 4px 0 0;background:#f0f0f0;padding:10px 20px;margin:0}
.tabs button:hover{background:#e0e0e0}
.tabs button.active{background:#007bff;color:#fff}
.tab-content{min-height:100px}
.tab-panel{padding:20px;background:#f9f9f9;border-radius:0 0 4px 4px}
.tab-panel h3{margin-top:0;color:#333}
.stepper{display:flex;align-items:center;justify-content:center;margin:20px 0}
.step{display:flex;flex-direction:column;align-items:center;gap:8px}
.step button{width:40px;height:40px;border-radius:50%;border:2px solid #ccc;background:#fff;font-weight:700;cursor:pointer}
.step.active button{border-color:#007bff;background:#007bff;color:#fff}
.step.completed button{border-color:#28a745;background:#28a745;color:#fff}
.step span{font-size:12px;color:#666}
.step-line{width:60px;height:2px;background:#ccc;margin:0 8px 24px}
.step-line.completed{background:#28a745}
.step-content{min-height:80px}
.step-panel{padding:20px;background:#f0f7ff;border-radius:4px;text-align:center}
.step-panel h3{margin-top:0;color:#004085}
.step-buttons{display:flex;justify-content:space-between;align-items:center;margin-top:15px;padding-top:15px;border-top:1px solid #eee}
`
}
