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
	return p.Html(`<div class="demo">
		<h1>Routing & Navigation</h1>

		<section>
			<h2>Tab Navigation</h2>
			<p>Simple tab-based navigation pattern:</p>

			<div class="tabs">`,
		p.ClassIf(`<button>Home</button>`, "active", r.CurrentTab.Eq("home")).WithOnClick(func() { r.GoHome() }),
		p.ClassIf(`<button>Profile</button>`, "active", r.CurrentTab.Eq("profile")).WithOnClick(func() { r.GoProfile() }),
		p.ClassIf(`<button>Settings</button>`, "active", r.CurrentTab.Eq("settings")).WithOnClick(func() { r.GoSettings() }),
		p.ClassIf(`<button>Notifications</button>`, "active", r.CurrentTab.Eq("notifications")).WithOnClick(func() { r.GoNotifications() }),
		`</div>

			<div class="tab-content">`,
		p.If(r.CurrentTab.Eq("home"),
			p.Html(`<div class="tab-panel"><h3>Home</h3><p>Welcome to the home tab! This is the default view.</p></div>`),
		).ElseIf(r.CurrentTab.Eq("profile"),
			p.Html(`<div class="tab-panel"><h3>Profile</h3><p>View and edit your profile information here.</p></div>`),
		).ElseIf(r.CurrentTab.Eq("settings"),
			p.Html(`<div class="tab-panel"><h3>Settings</h3><p>Configure your application settings.</p></div>`),
		).ElseIf(r.CurrentTab.Eq("notifications"),
			p.Html(`<div class="tab-panel"><h3>Notifications</h3><p>View your recent notifications.</p></div>`),
		),
		`</div>
		</section>

		<section>
			<h2>Wizard / Stepper</h2>
			<p>Multi-step form navigation:</p>

			<div class="stepper">`,
		p.ClassIf(`<div class="step">`, "completed", r.CurrentStep.Gt(1), "active", r.CurrentStep.Eq(1)),
		`<button `, p.OnClick(func() { r.Step1() }), `>1</button>
					<span>Details</span>
				</div>`,
		p.ClassIf(`<div class="step-line">`, "completed", r.CurrentStep.Gt(1)), `</div>`,
		p.ClassIf(`<div class="step">`, "completed", r.CurrentStep.Gt(2), "active", r.CurrentStep.Eq(2)),
		`<button `, p.OnClick(func() { r.Step2() }), `>2</button>
					<span>Address</span>
				</div>`,
		p.ClassIf(`<div class="step-line">`, "completed", r.CurrentStep.Gt(2)), `</div>`,
		p.ClassIf(`<div class="step">`, "completed", r.CurrentStep.Gt(3), "active", r.CurrentStep.Eq(3)),
		`<button `, p.OnClick(func() { r.Step3() }), `>3</button>
					<span>Payment</span>
				</div>`,
		p.ClassIf(`<div class="step-line">`, "completed", r.CurrentStep.Gt(3)), `</div>`,
		p.ClassIf(`<div class="step">`, "active", r.CurrentStep.Eq(4)),
		`<button `, p.OnClick(func() { r.Step4() }), `>4</button>
					<span>Confirm</span>
				</div>
			</div>

			<div class="step-content">`,
		p.If(r.CurrentStep.Eq(1),
			p.Html(`<div class="step-panel"><h3>Step 1: Personal Details</h3><p>Enter your name and email address.</p></div>`),
		).ElseIf(r.CurrentStep.Eq(2),
			p.Html(`<div class="step-panel"><h3>Step 2: Shipping Address</h3><p>Enter your shipping address.</p></div>`),
		).ElseIf(r.CurrentStep.Eq(3),
			p.Html(`<div class="step-panel"><h3>Step 3: Payment Method</h3><p>Choose your payment method.</p></div>`),
		).ElseIf(r.CurrentStep.Eq(4),
			p.Html(`<div class="step-panel"><h3>Step 4: Confirmation</h3><p>Review and confirm your order.</p></div>`),
		),
		`</div>

			<div class="step-buttons">
				<button `, p.OnClick(func() { r.PrevStep() }), ` `, p.ClassIf("disabled", r.CurrentStep.Eq(1)), `>Previous</button>
				<span>Step `, p.Bind(r.CurrentStep), ` of 4</span>
				<button `, p.OnClick(func() { r.NextStep() }), ` `, p.ClassIf("disabled", r.CurrentStep.Eq(4)), `>Next</button>
			</div>
		</section>
	</div>`)
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
