package main

import p "github.com/tbocek/preveltekit/v2"

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
	return p.Div(p.Attr("class", "demo"),
		p.H1("Complex"),

		p.Section(
			p.H2("Shared Store"),
			p.P("Theme is shared between parent and all child components:"),
			p.P("Current theme: ", p.Strong(c.Theme)),
			p.Div(p.Attr("class", "buttons"),
				p.Button("Light").On("click", func() { c.Theme.Set("light"); c.Log.Append("Theme: light") }),
				p.Button("Dark").On("click", func() { c.Theme.Set("dark"); c.Log.Append("Theme: dark") }),
				p.Button("Blue").On("click", func() { c.Theme.Set("blue"); c.Log.Append("Theme: blue") }),
			),
		),

		p.Section(
			p.H2("Nested Store[Component] Tabs"),
			p.P("Store[Component] inside a child component (not router-level):"),
			p.Div(p.Attr("class", "tab-bar"),
				p.Button("Dashboard").On("click", func() { c.SetTab(c.tabDashboard); c.Log.Append("Tab: dashboard") }),
				p.Button("Settings").On("click", func() { c.SetTab(c.tabSettings); c.Log.Append("Tab: settings") }),
				p.Button("Activity").On("click", func() { c.SetTab(c.tabActivity); c.Log.Append("Tab: activity") }),
			),
			p.Div(p.Attr("class", "tab-content"), p.Attr("data-theme", c.Theme),
				c.ActiveTab,
			),
		),

		p.Section(
			p.H2("Event Log (list inside if-block)"),
			p.P("Log entries: ", p.Strong(c.Log.Len())),
			p.Div(p.Attr("class", "buttons"),
				p.Button("Toggle Log").On("click", func() { c.ShowLog.Set(!c.ShowLog.Get()) }),
				p.Button("Clear Log").On("click", c.ClearLog),
			),
			p.If(p.Cond(func() bool { return c.ShowLog.Get() }, c.ShowLog),
				p.Ul(p.Attr("class", "log"),
					p.Each(c.Log, func(entry string, i int) p.Node {
						return p.Li(p.Attr("class", "log-entry"), p.Span(p.Attr("class", "log-idx"), p.Itoa(i)), " ", entry)
					}).Else(
						p.Li(p.Attr("class", "empty"), "No log entries"),
					),
				),
			).Else(
				p.P(p.Attr("class", "hint"), "Log hidden"),
			),
		),

		p.Section(
			p.H2("Code"),
			p.Pre(p.Attr("class", "code"), `// shared store: pass same *Store to multiple components
theme := p.New("light")
dashboard := &Dashboard{Theme: theme}
settings  := &Settings{Theme: theme}
// both read and write the same store

// Store[Component] as local tabs (not router):
activeTab := p.New[p.Component](dashboard)
activeTab.WithOptions(dashboard, settings, activity)
// embed in HTML — swaps component on Set()

// Attr(): dynamic attributes from stores
p.Div(p.Attr("data-theme", theme))

// nested reactivity: Each inside If
p.If(p.Cond(func() bool { return showLog.Get() }, showLog),
    p.Ul(
        p.Each(log, func(entry string, i int) p.Node {
            return p.Li(entry)
        }),
    ),
)`),
		),
	)
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
	return p.Div(
		p.H3("Dashboard"),
		p.P("Theme: ", p.Strong(d.Theme)),
		p.P("Stats value: ", p.Strong(d.Stats)),
		p.Div(p.Attr("class", "buttons"),
			p.Button("+10").On("click", func() { d.Stats.Update(func(v int) int { return v + 10 }); d.Log.Append("Stats +10") }),
			p.Button("Reset").On("click", func() { d.Stats.Set(0); d.Log.Append("Stats reset") }),
		),
	)
}

// Settings — tab component that modifies the shared Theme store
type Settings struct {
	Theme    *p.Store[string]
	Log      *p.List[string]
	FontSize *p.Store[int]
}

func (s *Settings) Render() p.Node {
	return p.Div(
		p.H3("Settings"),
		p.P("Font size: ", p.Strong(s.FontSize), "px"),
		p.Div(p.Attr("class", "buttons"),
			p.Button("Small (12)").On("click", func() { s.FontSize.Set(12); s.Log.Append("Font: 12px") }),
			p.Button("Medium (14)").On("click", func() { s.FontSize.Set(14); s.Log.Append("Font: 14px") }),
			p.Button("Large (18)").On("click", func() { s.FontSize.Set(18); s.Log.Append("Font: 18px") }),
		),
		p.P("Change theme from child:"),
		p.Div(p.Attr("class", "buttons"),
			p.Button("Light").On("click", func() { s.Theme.Set("light"); s.Log.Append("Settings: light") }),
			p.Button("Dark").On("click", func() { s.Theme.Set("dark"); s.Log.Append("Settings: dark") }),
		),
	)
}

// Activity — tab component that displays the shared Log as a list
type Activity struct {
	Theme *p.Store[string]
	Log   *p.List[string]
}

func (a *Activity) Render() p.Node {
	return p.Div(
		p.H3("Activity"),
		p.P("Theme: ", p.Strong(a.Theme), " | Entries: ", p.Strong(a.Log.Len())),
		p.If(p.Cond(func() bool { return a.Log.Len().Get() > 0 }, a.Log.Len()),
			p.Ul(p.Attr("class", "activity-list"),
				p.Each(a.Log, func(entry string, i int) p.Node {
					return p.Li(entry)
				}),
			),
		).Else(
			p.P(p.Attr("class", "empty"), "No activity yet"),
		),
	)
}

func (a *Activity) Style() string {
	return `
.activity-list{list-style:none;padding:0;margin:5px 0}
.activity-list li{padding:4px 8px;margin:2px 0;background:rgba(0,0,0,.05);border-radius:3px;font-size:13px}`
}
