package main

import p "github.com/tbocek/preveltekit/v2"

type Manual struct{}

func (m *Manual) New() p.Component {
	return &Manual{}
}

func (m *Manual) Render() p.Node {
	return p.Div(p.Attr("class", "manual page"),
		p.Div(p.Attr("class", "container"),
			p.H1("Manual"),
			p.P(p.Attr("class", "page-intro"), "API reference for PrevelteKit. Each section shows the API and a code example."),

			p.Nav(p.Attr("class", "toc"),
				p.H2("Contents"),
				p.Ul(
					p.Li(p.A(p.Attr("href", "#stores"), p.Attr("external", ""), "Stores")),
					p.Li(p.A(p.Attr("href", "#elements-and-binding"), p.Attr("external", ""), "Elements & Binding")),
					p.Li(p.A(p.Attr("href", "#events"), p.Attr("external", ""), "Events")),
					p.Li(p.A(p.Attr("href", "#conditionals"), p.Attr("external", ""), "Conditionals")),
					p.Li(p.A(p.Attr("href", "#lists"), p.Attr("external", ""), "Lists")),
					p.Li(p.A(p.Attr("href", "#components"), p.Attr("external", ""), "Components")),
					p.Li(p.A(p.Attr("href", "#routing"), p.Attr("external", ""), "Routing")),
					p.Li(p.A(p.Attr("href", "#fetch"), p.Attr("external", ""), "Fetch")),
					p.Li(p.A(p.Attr("href", "#storage"), p.Attr("external", ""), "Storage")),
					p.Li(p.A(p.Attr("href", "#timers"), p.Attr("external", ""), "Timers")),
					p.Li(p.A(p.Attr("href", "#lifecycle"), p.Attr("external", ""), "Lifecycle")),
				),
			),

			p.Section(p.Attr("id", "stores"),
				p.H2("Stores"),
				p.P("Reactive state containers. When a store value changes, all bound DOM elements update automatically."),
				p.Pre(p.Code(`// Create a store with an initial value
count := p.New(0)              // *Store[int]
name := p.New("hello")         // *Store[string]
dark := p.New(false)           // *Store[bool]

// Read and write
count.Get()                    // returns current value
count.Set(5)                   // set new value
count.Update(func(v int) int { // transform current value
    return v + 1
})

// Subscribe to changes
count.OnChange(func(v int) {
    // called whenever count changes
})`)),
			),

			p.Section(p.Attr("id", "elements-and-binding"),
				p.H2("Elements & Binding"),
				p.P(p.RawHTML("Build DOM trees with typed element functions. Embed stores directly — they become live text nodes.")),
				p.Pre(p.Code("// Typed element functions with reactive store interpolation\np.P(\"Count: \", p.Strong(count))\n\n// Two-way binding for inputs\np.Input(p.Attr(\"type\", \"text\")).Bind(name)       // *Store[string]\np.Input(p.Attr(\"type\", \"text\")).Bind(age)        // *Store[int]\np.Input(p.Attr(\"type\", \"checkbox\")).Bind(dark)    // *Store[bool]\np.Textarea().Bind(notes)                         // *Store[string]\n\n// Dynamic attributes\np.Div(\"content\").Attr(\"data-theme\", theme)\n\n// Conditional attributes (additive for same attribute name)\np.Div(\"content\").AttrIf(\"class\",\n    p.Cond(func() bool { return dark.Get() }, dark), \"active\")\n\n// Raw HTML rendering (not escaped)\np.BindAsHTML(rawHTML)\n\n// Static raw HTML (entities, inline markup)\np.RawHTML(\"&copy; 2024\")")),
			),

			p.Section(p.Attr("id", "events"),
				p.H2("Events"),
				p.P(p.RawHTML("Attach event handlers with <code>.On()</code>. Chain modifiers for common patterns.")),
				p.Pre(p.Code("// Click handler\np.Button(\"Click\").On(\"click\", handler)\n\n// Form submit with preventDefault\np.Form(...).On(\"submit\", handler).PreventDefault()\n\n// Stop event bubbling\np.Button(\"Inner\").On(\"click\", handler).StopPropagation()\n\n// Inline handler\np.Button(\"+5\").On(\"click\", func() {\n    count.Update(func(v int) int { return v + 5 })\n})")),
			),

			p.Section(p.Attr("id", "conditionals"),
				p.H2("Conditionals"),
				p.P(p.RawHTML("Show or hide content reactively with <code>p.If()</code>. Supports <code>ElseIf</code> and <code>Else</code> chains.")),
				p.Pre(p.Code("// Simple if\np.If(p.Cond(func() bool { return count.Get() > 0 }, count),\n    p.P(\"Positive\"),\n)\n\n// If / ElseIf / Else\np.If(p.Cond(func() bool { return score.Get() >= 90 }, score),\n    p.P(\"Grade: A\"),\n).ElseIf(p.Cond(func() bool { return score.Get() >= 80 }, score),\n    p.P(\"Grade: B\"),\n).Else(\n    p.P(\"Grade: F\"),\n)\n\n// p.Cond(predicateFn, ...dependencyStores)")),
			),

			p.Section(p.Attr("id", "lists"),
				p.H2("Lists"),
				p.P(p.RawHTML("Reactive lists with <code>p.NewList()</code>. Render with <code>p.Each()</code>.")),
				p.Pre(p.Code("// Create a reactive list\nitems := p.NewList[string](\"Apple\", \"Banana\", \"Cherry\")\n\n// Mutate — triggers re-render\nitems.Append(\"Date\")\nitems.RemoveAt(0)\nitems.Set([]string{\"Mango\", \"Papaya\"})\nitems.Clear()\n\n// Reactive length\nitems.Len()  // *Store[int] — updates automatically\n\n// Render each item\np.Each(items, func(item string, i int) p.Node {\n    return p.Li(p.Itoa(i), \": \", item)\n}).Else(\n    p.P(\"No items\"),\n)")),
			),

			p.Section(p.Attr("id", "components"),
				p.H2("Components"),
				p.P(p.RawHTML("Components are Go structs implementing <code>Render() p.Node</code>. Use <code>p.Comp()</code> to embed them.")),
				p.Pre(p.Code("// Define a component\ntype Badge struct {\n    Label *p.Store[string]\n}\n\nfunc (b *Badge) Render() p.Node {\n    return p.Span(p.Attr(\"class\", \"badge\"), b.Label)\n}\n\n// Scoped CSS\nfunc (b *Badge) Style() string {\n    return `.badge{background:#007bff;color:#fff}`\n}\n\n// Use a component (props are struct fields)\np.Comp(&Badge{Label: p.New(\"New\")})\n\n// Component with slot (child content)\np.Comp(&Card{Title: p.New(\"Hello\")},\n    p.P(\"Slot content here\"),\n)\n\n// Inside the component, render slot content:\np.Slot()\n\n// Callback props for component events\ntype Button struct {\n    Label   *p.Store[string]\n    OnClick func()\n}\n\n// Shared stores: pass the same *Store to multiple components\ntheme := p.New(\"light\")\np.Comp(&Header{Theme: theme})\np.Comp(&Sidebar{Theme: theme})\n// both components read/write the same store")),
			),

			p.Section(p.Attr("id", "routing"),
				p.H2("Routing"),
				p.P(p.RawHTML("Client-side routing with <code>p.NewRouter()</code>. Each route maps to a component and an SSR HTML file.")),
				p.Pre(p.Code("// Define routes\nroutes := []p.Route{\n    {Path: \"/\", HTMLFile: \"index.html\", SSRPath: \"/\", Component: home},\n    {Path: \"/about\", HTMLFile: \"about.html\", SSRPath: \"/about\", Component: about},\n}\n\n// Start the router (in OnMount)\nrouter := p.NewRouter(currentComponent, routes, \"unique-id\")\nrouter.NotFound(func() { currentComponent.Set(nil) })\nrouter.Start()\n\n// Store[Component] holds the active route component\ncurrentComponent *p.Store[p.Component]\n\n// Links: client-side (default) vs server-side\n// <a href=\"/about\">About</a>           (SPA navigation)\n// <a href=\"/about\" external>About</a>  (full page reload)\n\n// Store[Component] for local tabs (non-router):\nactiveTab := p.New[p.Component](tab1)\nactiveTab.WithOptions(tab1, tab2, tab3)\nactiveTab.Set(tab2)  // switches displayed component")),
			),

			p.Section(p.Attr("id", "fetch"),
				p.H2("Fetch"),
				p.P(p.RawHTML("Type-safe HTTP requests with automatic JSON encoding/decoding via <code>js</code> struct tags.")),
				p.Pre(p.Code("// Define response type with js tags\ntype User struct {\n    ID   int    `js:\"id\"`\n    Name string `js:\"name\"`\n}\n\n// GET\ngo func() {\n    user, err := p.Get[User](\"https://api.example.com/user/1\")\n}()\n\n// POST (send body, decode response)\ngo func() {\n    created, err := p.Post[User](url, newUser)\n}()\n\n// PUT, PATCH, DELETE\nresult, err := p.Put[T](url, body)\nresult, err := p.Patch[T](url, body)\nresult, err := p.Delete[T](url)\n\n// Advanced: custom headers and abort\nsignal, abort := p.NewAbortController()\ngo func() {\n    result, err := p.Fetch[T](url, &p.FetchOptions{\n        Method:  \"GET\",\n        Headers: map[string]string{\"Authorization\": \"Bearer token\"},\n        Signal:  signal,\n    })\n}()\nabort()  // cancel the request")),
			),

			p.Section(p.Attr("id", "storage"),
				p.H2("Storage"),
				p.P(p.RawHTML("Persist state to localStorage. <code>LocalStore</code> auto-syncs on every <code>Set()</code>.")),
				p.Pre(p.Code(`// Auto-persisted store (syncs on every Set)
theme := p.NewLocalStore("theme", "light")
theme.Set("dark")  // automatically saved to localStorage
theme.Store        // *Store[string] — use in any element like any store

// Manual localStorage API
p.SetStorage("notes", "hello")
saved := p.GetStorage("notes")
p.RemoveStorage("notes")
p.ClearStorage()`)),
			),

			p.Section(p.Attr("id", "timers"),
				p.H2("Timers"),
				p.P("Debounce, throttle, setTimeout, and setInterval — all return cleanup functions."),
				p.Pre(p.Code(`// Debounce: fires after idle period
doSearch, cleanup := p.Debounce(300, func() {
    // fires 300ms after last call
})
doSearch()   // call repeatedly — only last one fires
cleanup()    // cancel pending

// Throttle: max once per interval
onClick := p.Throttle(500, func() {
    // max once per 500ms
})

// SetTimeout: fires once after delay
cancel := p.SetTimeout(2000, func() {
    // fires after 2 seconds
})
cancel()  // cancel before it fires

// SetInterval: fires repeatedly
stop := p.SetInterval(60000, func() {
    // fires every 60 seconds
})
stop()  // stop the interval`)),
			),

			p.Section(p.Attr("id", "lifecycle"),
				p.H2("Lifecycle"),
				p.P("Components can implement lifecycle hooks for setup and teardown."),
				p.Pre(p.Code("// OnMount: called when component becomes active\nfunc (c *MyComp) OnMount() {\n    if p.IsBuildTime {\n        return  // skip side effects during SSR\n    }\n    // fetch data, start timers, etc.\n}\n\n// OnDestroy: called when component is removed (e.g. route change)\nfunc (c *MyComp) OnDestroy() {\n    if c.stopTimer != nil {\n        c.stopTimer()\n    }\n}\n\n// New: constructor — create stores, initialize state\nfunc (c *MyComp) New() p.Component {\n    return &MyComp{\n        Count: p.New(0),\n        Name:  p.New(\"\"),\n    }\n}\n\n// GlobalStyle: CSS applied to entire page (on App component)\nfunc (a *App) GlobalStyle() string { return `body{margin:0}` }\n\n// Style: scoped CSS (auto-prefixed to this component)\nfunc (c *MyComp) Style() string { return `.btn{color:red}` }\n\n// p.IsBuildTime: true during SSR, false in WASM\n// Use to guard browser-only code (fetch, timers, DOM access)")),
			),

		),
	)
}

func (m *Manual) Style() string {
	return `
.toc{margin-bottom:40px;padding:20px;background:#f8f9fa;border-radius:8px}
.toc h2{font-size:1em;margin-bottom:12px;color:#1a1a2e}
.toc ul{list-style:none;padding:0;margin:0;display:flex;flex-wrap:wrap;gap:8px}
.toc li a{display:inline-block;padding:6px 14px;background:#1a1a2e;color:#fff;border-radius:4px;font-size:13px;transition:background .2s}
.toc li a:hover{background:#0f3460}

.manual section{margin-bottom:32px;padding:24px;border:1px solid #e5e7eb;border-radius:8px;background:#fff}
.manual section h2{font-size:1.3em;color:#1a1a2e;margin-bottom:8px}
.manual section > p{color:#666;margin-bottom:16px;font-size:.95em}

@media(max-width:768px){
.toc ul{flex-direction:column}
}
`
}
