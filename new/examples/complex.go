package main

import p "preveltekit"

// Complex demonstrates nested stores, Store[Component] inside child components,
// shared stores across parent/child, and deeply nested reactive blocks.
type Complex struct {
	// Shared store: parent and children both read/write
	Theme *p.Store[string]

	// Tab panel: Store[Component] inside this component (not at router level)
	ActiveTab    *p.Store[p.Component]
	tabDashboard *Dashboard
	tabSettings  *Settings
	tabActivity  *Activity

	// Nested blocks: list inside if-block
	ShowLog *p.Store[bool]
	Log     *p.List[string]
}

func (c *Complex) New() p.Component {
	theme := p.New("light")
	log := p.NewList[string]("App started")

	dashboard := &Dashboard{Theme: theme, Log: log, Stats: p.New(42)}
	settings := &Settings{Theme: theme, Log: log, FontSize: p.New(14)}
	activity := &Activity{Theme: theme, Log: log}

	activeTab := p.New[p.Component](dashboard)
	activeTab.WithOptions(dashboard, settings, activity)

	return &Complex{
		Theme:        theme,
		ActiveTab:    activeTab,
		tabDashboard: dashboard,
		tabSettings:  settings,
		tabActivity:  activity,
		ShowLog:      p.New(true),
		Log:          log,
	}
}

func (c *Complex) SetTab(comp p.Component) {
	c.ActiveTab.Set(comp)
}

func (c *Complex) ClearLog() {
	c.Log.Clear()
	c.Log.Append("Log cleared")
}

func (c *Complex) Render() p.Node {
	return p.Html(`<div class="demo">
		<h1>Complex</h1>

		<section>
			<h2>Shared Store</h2>
			<p>Theme is shared between parent and all child components:</p>
			<p>Current theme: <strong>`, c.Theme, `</strong></p>
			<div class="buttons">
				`, p.Html(`<button>Light</button>`).On("click", func() { c.Theme.Set("light"); c.Log.Append("Theme: light") }), `
				`, p.Html(`<button>Dark</button>`).On("click", func() { c.Theme.Set("dark"); c.Log.Append("Theme: dark") }), `
				`, p.Html(`<button>Blue</button>`).On("click", func() { c.Theme.Set("blue"); c.Log.Append("Theme: blue") }), `
			</div>
		</section>

		<section>
			<h2>Nested Store[Component] Tabs</h2>
			<p>Store[Component] inside a child component (not router-level):</p>
			<div class="tab-bar">
				`, p.Html(`<button>Dashboard</button>`).On("click", func() { c.SetTab(c.tabDashboard); c.Log.Append("Tab: dashboard") }),
		p.Html(`<button>Settings</button>`).On("click", func() { c.SetTab(c.tabSettings); c.Log.Append("Tab: settings") }),
		p.Html(`<button>Activity</button>`).On("click", func() { c.SetTab(c.tabActivity); c.Log.Append("Tab: activity") }), `
			</div>
			<div class="tab-content" `, p.Attr("data-theme", c.Theme), `>`,
		c.ActiveTab,
		`</div>
		</section>

		<section>
			<h2>Event Log (list inside if-block)</h2>
			<p>Log entries: <strong>`, c.Log.Len(), `</strong></p>
			<div class="buttons">
				`, p.Html(`<button>Toggle Log</button>`).On("click", func() { c.ShowLog.Set(!c.ShowLog.Get()) }), `
				`, p.Html(`<button>Clear Log</button>`).On("click", c.ClearLog), `
			</div>`,
		p.If(p.Cond(func() bool { return c.ShowLog.Get() }, c.ShowLog),
			p.Html(`<ul class="log">`,
				p.Each(c.Log, func(entry string, i int) p.Node {
					return p.Html(`<li class="log-entry"><span class="log-idx">`, p.Itoa(i), `</span> `, entry, `</li>`)
				}).Else(
					p.Html(`<li class="empty">No log entries</li>`),
				),
				`</ul>`),
		).Else(
			p.Html(`<p class="hint">Log hidden</p>`),
		), `
		</section>

		<section>
			<h2>Code</h2>
			<pre class="code">// shared store: pass same *Store to multiple components
theme := p.New("light")
dashboard := &amp;Dashboard{Theme: theme}
settings  := &amp;Settings{Theme: theme}
// both read and write the same store

// Store[Component] as local tabs (not router):
activeTab := p.New[p.Component](dashboard)
activeTab.WithOptions(dashboard, settings, activity)
// embed in HTML — swaps component on Set()

// Attr(): dynamic attributes from stores
p.Html(`+"`"+`&lt;div>`+"`"+`).Attr("data-theme", theme)

// nested reactivity: Each inside If
p.If(p.Cond(func() bool { return showLog.Get() }, showLog),
    p.Html(`+"`"+`&lt;ul>`+"`"+`,
        p.Each(log, func(entry string, i int) p.Node {
            return p.Html(`+"`"+`&lt;li>`+"`"+`, entry, `+"`"+`&lt;/li>`+"`"+`)
        }),
    `+"`"+`&lt;/ul>`+"`"+`),
)</pre>
		</section>
	</div>`)
}

func (c *Complex) Style() string {
	return `
.demo pre.code{background:#1a1a2e;color:#e0e0e0;font-size:12px;margin-top:12px}
.tab-bar{display:flex;gap:4px;margin-bottom:0}
.tab-bar button{border-bottom:none;border-radius:4px 4px 0 0}
.tab-content{border:1px solid #ddd;border-radius:0 4px 4px 4px;padding:15px;min-height:120px}
.tab-content[data-theme=dark]{background:#2d2d2d;color:#eee}
.tab-content[data-theme=blue]{background:#e8f0fe;color:#1a237e}
.tab-content[data-theme=light]{background:#fff;color:#333}
.log{list-style:none;padding:0;margin:10px 0;max-height:200px;overflow-y:auto}
.log-entry{padding:4px 8px;margin:2px 0;background:#f8f8f8;border-radius:3px;font-size:13px;font-family:monospace}
.log-idx{display:inline-block;width:20px;color:#999;font-size:11px}
.empty{color:#999;font-style:italic;padding:10px}
`
}

// Dashboard — tab component that reads shared Theme store and writes to shared Log
type Dashboard struct {
	Theme *p.Store[string]
	Log   *p.List[string]
	Stats *p.Store[int]
}

func (d *Dashboard) Render() p.Node {
	return p.Html(`<div>
		<h3>Dashboard</h3>
		<p>Theme: <strong>`, d.Theme, `</strong></p>
		<p>Stats value: <strong>`, d.Stats, `</strong></p>
		<div class="buttons">
			`, p.Html(`<button>+10</button>`).On("click", func() { d.Stats.Update(func(v int) int { return v + 10 }); d.Log.Append("Stats +10") }), `
			`, p.Html(`<button>Reset</button>`).On("click", func() { d.Stats.Set(0); d.Log.Append("Stats reset") }), `
		</div>
	</div>`)
}

// Settings — tab component that modifies the shared Theme store
type Settings struct {
	Theme    *p.Store[string]
	Log      *p.List[string]
	FontSize *p.Store[int]
}

func (s *Settings) Render() p.Node {
	return p.Html(`<div>
		<h3>Settings</h3>
		<p>Font size: <strong>`, s.FontSize, `</strong>px</p>
		<div class="buttons">
			`, p.Html(`<button>Small (12)</button>`).On("click", func() { s.FontSize.Set(12); s.Log.Append("Font: 12px") }), `
			`, p.Html(`<button>Medium (14)</button>`).On("click", func() { s.FontSize.Set(14); s.Log.Append("Font: 14px") }), `
			`, p.Html(`<button>Large (18)</button>`).On("click", func() { s.FontSize.Set(18); s.Log.Append("Font: 18px") }), `
		</div>
		<p>Change theme from child:</p>
		<div class="buttons">
			`, p.Html(`<button>Light</button>`).On("click", func() { s.Theme.Set("light"); s.Log.Append("Settings: light") }), `
			`, p.Html(`<button>Dark</button>`).On("click", func() { s.Theme.Set("dark"); s.Log.Append("Settings: dark") }), `
		</div>
	</div>`)
}

// Activity — tab component that displays the shared Log as a list
type Activity struct {
	Theme *p.Store[string]
	Log   *p.List[string]
}

func (a *Activity) Render() p.Node {
	return p.Html(`<div>
		<h3>Activity</h3>
		<p>Theme: <strong>`, a.Theme, `</strong> | Entries: <strong>`, a.Log.Len(), `</strong></p>
		`, p.If(p.Cond(func() bool { return a.Log.Len().Get() > 0 }, a.Log.Len()),
		p.Html(`<ul class="activity-list">`,
			p.Each(a.Log, func(entry string, i int) p.Node {
				return p.Html(`<li>`, entry, `</li>`)
			}),
			`</ul>`),
	).Else(
		p.Html(`<p class="empty">No activity yet</p>`),
	), `
	</div>`)
}

func (a *Activity) Style() string {
	return `
.activity-list{list-style:none;padding:0;margin:5px 0}
.activity-list li{padding:4px 8px;margin:2px 0;background:rgba(0,0,0,.05);border-radius:3px;font-size:13px}
`
}
