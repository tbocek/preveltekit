package main

import p "preveltekit"

// Routing showcase - demonstrates navigation patterns
type Routing struct {
	CurrentTab  *p.Store[string]
	CurrentStep *p.Store[int]
}

func (r *Routing) OnMount() {
	r.CurrentTab.Set("home")
	r.CurrentStep.Set(1)
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
	return p.Div(p.Class("demo"),
		p.H1("Routing & Navigation"),

		p.Section(
			p.H2("Tab Navigation"),
			p.P("Simple tab-based navigation pattern:"),

			p.Div(p.Class("tabs"),
				p.Button("Home", p.OnClick(r.GoHome), p.ClassIf("active", r.CurrentTab.Eq("home"))),
				p.Button("Profile", p.OnClick(r.GoProfile), p.ClassIf("active", r.CurrentTab.Eq("profile"))),
				p.Button("Settings", p.OnClick(r.GoSettings), p.ClassIf("active", r.CurrentTab.Eq("settings"))),
				p.Button("Notifications", p.OnClick(r.GoNotifications), p.ClassIf("active", r.CurrentTab.Eq("notifications"))),
			),

			p.Div(p.Class("tab-content"),
				p.If(r.CurrentTab.Eq("home"),
					p.Div(p.Class("tab-panel"),
						p.H3("Home"),
						p.P("Welcome to the home tab! This is the default view."),
					),
				).ElseIf(r.CurrentTab.Eq("profile"),
					p.Div(p.Class("tab-panel"),
						p.H3("Profile"),
						p.P("View and edit your profile information here."),
					),
				).ElseIf(r.CurrentTab.Eq("settings"),
					p.Div(p.Class("tab-panel"),
						p.H3("Settings"),
						p.P("Configure your application settings."),
					),
				).ElseIf(r.CurrentTab.Eq("notifications"),
					p.Div(p.Class("tab-panel"),
						p.H3("Notifications"),
						p.P("View your recent notifications."),
					),
				),
			),
		),

		p.Section(
			p.H2("Wizard / Stepper"),
			p.P("Multi-step form navigation:"),

			p.Div(p.Class("stepper"),
				p.Div(p.Class("step"), p.ClassIf("completed", r.CurrentStep.Gt(1)), p.ClassIf("active", r.CurrentStep.Eq(1)),
					p.Button("1", p.OnClick(r.Step1)),
					p.Span("Details"),
				),
				p.Div(p.Class("step-line"), p.ClassIf("completed", r.CurrentStep.Gt(1))),
				p.Div(p.Class("step"), p.ClassIf("completed", r.CurrentStep.Gt(2)), p.ClassIf("active", r.CurrentStep.Eq(2)),
					p.Button("2", p.OnClick(r.Step2)),
					p.Span("Address"),
				),
				p.Div(p.Class("step-line"), p.ClassIf("completed", r.CurrentStep.Gt(2))),
				p.Div(p.Class("step"), p.ClassIf("completed", r.CurrentStep.Gt(3)), p.ClassIf("active", r.CurrentStep.Eq(3)),
					p.Button("3", p.OnClick(r.Step3)),
					p.Span("Payment"),
				),
				p.Div(p.Class("step-line"), p.ClassIf("completed", r.CurrentStep.Gt(3))),
				p.Div(p.Class("step"), p.ClassIf("active", r.CurrentStep.Eq(4)),
					p.Button("4", p.OnClick(r.Step4)),
					p.Span("Confirm"),
				),
			),

			p.Div(p.Class("step-content"),
				p.If(r.CurrentStep.Eq(1),
					p.Div(p.Class("step-panel"),
						p.H3("Step 1: Personal Details"),
						p.P("Enter your name and email address."),
					),
				).ElseIf(r.CurrentStep.Eq(2),
					p.Div(p.Class("step-panel"),
						p.H3("Step 2: Shipping Address"),
						p.P("Enter your shipping address."),
					),
				).ElseIf(r.CurrentStep.Eq(3),
					p.Div(p.Class("step-panel"),
						p.H3("Step 3: Payment Method"),
						p.P("Choose your payment method."),
					),
				).ElseIf(r.CurrentStep.Eq(4),
					p.Div(p.Class("step-panel"),
						p.H3("Step 4: Confirmation"),
						p.P("Review and confirm your order."),
					),
				),
			),

			p.Div(p.Class("step-buttons"),
				p.Button("Previous", p.OnClick(r.PrevStep), p.ClassIf("disabled", r.CurrentStep.Eq(1))),
				p.Span("Step ", p.Bind(r.CurrentStep), " of 4"),
				p.Button("Next", p.OnClick(r.NextStep), p.ClassIf("disabled", r.CurrentStep.Eq(4))),
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

func (r *Routing) HandleEvent(method string, args string) {
	switch method {
	case "GoHome":
		r.GoHome()
	case "GoProfile":
		r.GoProfile()
	case "GoSettings":
		r.GoSettings()
	case "GoNotifications":
		r.GoNotifications()
	case "NextStep":
		r.NextStep()
	case "PrevStep":
		r.PrevStep()
	case "Step1":
		r.Step1()
	case "Step2":
		r.Step2()
	case "Step3":
		r.Step3()
	case "Step4":
		r.Step4()
	}
}
