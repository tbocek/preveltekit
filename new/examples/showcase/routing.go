package main

import "reactive"

// Routing showcase - demonstrates navigation patterns
type Routing struct {
	CurrentTab  *reactive.Store[string]
	CurrentStep *reactive.Store[int]
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

func (r *Routing) Template() string {
	return `<div class="demo">
	<h1>Routing & Navigation</h1>

	<section>
		<h2>Tab Navigation</h2>
		<p>Simple tab-based navigation pattern:</p>

		<div class="tabs">
			<button @click="GoHome()" class:active={CurrentTab == "home"}>Home</button>
			<button @click="GoProfile()" class:active={CurrentTab == "profile"}>Profile</button>
			<button @click="GoSettings()" class:active={CurrentTab == "settings"}>Settings</button>
			<button @click="GoNotifications()" class:active={CurrentTab == "notifications"}>Notifications</button>
		</div>

		<div class="tab-content">
			{#if CurrentTab == "home"}
				<div class="tab-panel">
					<h3>Home</h3>
					<p>Welcome to the home tab! This is the default view.</p>
				</div>
			{:else if CurrentTab == "profile"}
				<div class="tab-panel">
					<h3>Profile</h3>
					<p>View and edit your profile information here.</p>
				</div>
			{:else if CurrentTab == "settings"}
				<div class="tab-panel">
					<h3>Settings</h3>
					<p>Configure your application settings.</p>
				</div>
			{:else if CurrentTab == "notifications"}
				<div class="tab-panel">
					<h3>Notifications</h3>
					<p>View your recent notifications.</p>
				</div>
			{/if}
		</div>
	</section>

	<section>
		<h2>Wizard / Stepper</h2>
		<p>Multi-step form navigation:</p>

		<div class="stepper">
			<div class="step" class:completed={CurrentStep > 1} class:active={CurrentStep == 1}>
				<button @click="Step1()">1</button>
				<span>Details</span>
			</div>
			<div class="step-line" class:completed={CurrentStep > 1}></div>
			<div class="step" class:completed={CurrentStep > 2} class:active={CurrentStep == 2}>
				<button @click="Step2()">2</button>
				<span>Address</span>
			</div>
			<div class="step-line" class:completed={CurrentStep > 2}></div>
			<div class="step" class:completed={CurrentStep > 3} class:active={CurrentStep == 3}>
				<button @click="Step3()">3</button>
				<span>Payment</span>
			</div>
			<div class="step-line" class:completed={CurrentStep > 3}></div>
			<div class="step" class:active={CurrentStep == 4}>
				<button @click="Step4()">4</button>
				<span>Confirm</span>
			</div>
		</div>

		<div class="step-content">
			{#if CurrentStep == 1}
				<div class="step-panel">
					<h3>Step 1: Personal Details</h3>
					<p>Enter your name and email address.</p>
				</div>
			{:else if CurrentStep == 2}
				<div class="step-panel">
					<h3>Step 2: Shipping Address</h3>
					<p>Enter your shipping address.</p>
				</div>
			{:else if CurrentStep == 3}
				<div class="step-panel">
					<h3>Step 3: Payment Method</h3>
					<p>Choose your payment method.</p>
				</div>
			{:else if CurrentStep == 4}
				<div class="step-panel">
					<h3>Step 4: Confirmation</h3>
					<p>Review and confirm your order.</p>
				</div>
			{/if}
		</div>

		<div class="step-buttons">
			<button @click="PrevStep()" class:disabled={CurrentStep == 1}>Previous</button>
			<span>Step {CurrentStep} of 4</span>
			<button @click="NextStep()" class:disabled={CurrentStep == 4}>Next</button>
		</div>
	</section>
</div>`
}

func (r *Routing) Style() string {
	return `
.demo { max-width: 700px; }
.demo h1 { color: #1a1a2e; margin-bottom: 20px; }
.demo section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; background: #fff; }
.demo h2 { margin-top: 0; color: #666; font-size: 1.1em; }
.demo button { padding: 8px 16px; margin: 4px; cursor: pointer; border: 1px solid #ccc; border-radius: 4px; background: #f5f5f5; }
.demo button:hover { background: #e5e5e5; }
.demo button.disabled { opacity: 0.5; cursor: not-allowed; }

.tabs { display: flex; gap: 4px; border-bottom: 2px solid #ddd; padding-bottom: 0; }
.tabs button { border: none; border-radius: 4px 4px 0 0; background: #f0f0f0; padding: 10px 20px; margin: 0; }
.tabs button:hover { background: #e0e0e0; }
.tabs button.active { background: #007bff; color: white; }

.tab-content { min-height: 100px; }
.tab-panel { padding: 20px; background: #f9f9f9; border-radius: 0 0 4px 4px; }
.tab-panel h3 { margin-top: 0; color: #333; }

.stepper { display: flex; align-items: center; justify-content: center; margin: 20px 0; }
.step { display: flex; flex-direction: column; align-items: center; gap: 8px; }
.step button { width: 40px; height: 40px; border-radius: 50%; border: 2px solid #ccc; background: white; font-weight: bold; cursor: pointer; }
.step.active button { border-color: #007bff; background: #007bff; color: white; }
.step.completed button { border-color: #28a745; background: #28a745; color: white; }
.step span { font-size: 12px; color: #666; }
.step-line { width: 60px; height: 2px; background: #ccc; margin: 0 8px; margin-bottom: 24px; }
.step-line.completed { background: #28a745; }

.step-content { min-height: 80px; }
.step-panel { padding: 20px; background: #f0f7ff; border-radius: 4px; text-align: center; }
.step-panel h3 { margin-top: 0; color: #004085; }

.step-buttons { display: flex; justify-content: space-between; align-items: center; margin-top: 15px; padding-top: 15px; border-top: 1px solid #eee; }
`
}
