package main

import p "github.com/tbocek/preveltekit/v2"

type Manual struct{}

func (m *Manual) New() p.Component {
	return &Manual{}
}

func (m *Manual) Render() p.Node {
	return p.Html(`<div class="manual">
	<div class="container">
		<h1>Manual</h1>
		<p class="intro">API reference for PrevelteKit. Each section shows the API and a code example.</p>

		<nav class="toc">
			<h2>Contents</h2>
			<ul>
				<li><a href="#stores" external>Stores</a></li>
				<li><a href="#html-and-binding" external>HTML &amp; Binding</a></li>
				<li><a href="#events" external>Events</a></li>
				<li><a href="#conditionals" external>Conditionals</a></li>
				<li><a href="#lists" external>Lists</a></li>
				<li><a href="#components" external>Components</a></li>
				<li><a href="#routing" external>Routing</a></li>
				<li><a href="#fetch" external>Fetch</a></li>
				<li><a href="#storage" external>Storage</a></li>
				<li><a href="#timers" external>Timers</a></li>
				<li><a href="#lifecycle" external>Lifecycle</a></li>
			</ul>
		</nav>

		<section id="stores">
			<h2>Stores</h2>
			<p>Reactive state containers. When a store value changes, all bound DOM elements update automatically.</p>
			<pre><code>// Create a store with an initial value
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
})</code></pre>
		</section>

		<section id="html-and-binding">
			<h2>HTML &amp; Binding</h2>
			<p>Build DOM trees with <code>p.Html()</code>. Embed stores directly — they become live text nodes.</p>
			<pre><code>// Static HTML with reactive store interpolation
p.Html(` + "`" + `&lt;p>Count: &lt;strong>` + "`" + `, count, ` + "`" + `&lt;/strong>&lt;/p>` + "`" + `)

// Two-way binding for inputs
p.Html(` + "`" + `&lt;input type="text">` + "`" + `).Bind(name)     // *Store[string]
p.Html(` + "`" + `&lt;input type="text">` + "`" + `).Bind(age)      // *Store[int]
p.Html(` + "`" + `&lt;input type="checkbox">` + "`" + `).Bind(dark)  // *Store[bool]
p.Html(` + "`" + `&lt;textarea>` + "`" + `).Bind(notes)              // *Store[string]

// Dynamic attributes
p.Html(` + "`" + `&lt;div>` + "`" + `).Attr("data-theme", theme)

// Conditional attributes (additive for same attribute name)
p.Html(` + "`" + `&lt;div>` + "`" + `).AttrIf("class",
    p.Cond(func() bool { return dark.Get() }, dark), "active")

// Raw HTML rendering (not escaped)
p.BindAsHTML(rawHTML)</code></pre>
		</section>

		<section id="events">
			<h2>Events</h2>
			<p>Attach event handlers with <code>.On()</code>. Chain modifiers for common patterns.</p>
			<pre><code>// Click handler
p.Html(` + "`" + `&lt;button>Click&lt;/button>` + "`" + `).On("click", handler)

// Form submit with preventDefault
p.Html(` + "`" + `&lt;form>...&lt;/form>` + "`" + `).On("submit", handler).PreventDefault()

// Stop event bubbling
p.Html(` + "`" + `&lt;button>Inner&lt;/button>` + "`" + `).On("click", handler).StopPropagation()

// Inline handler
p.Html(` + "`" + `&lt;button>+5&lt;/button>` + "`" + `).On("click", func() {
    count.Update(func(v int) int { return v + 5 })
})</code></pre>
		</section>

		<section id="conditionals">
			<h2>Conditionals</h2>
			<p>Show or hide content reactively with <code>p.If()</code>. Supports <code>ElseIf</code> and <code>Else</code> chains.</p>
			<pre><code>// Simple if
p.If(p.Cond(func() bool { return count.Get() > 0 }, count),
    p.Html(` + "`" + `&lt;p>Positive&lt;/p>` + "`" + `),
)

// If / ElseIf / Else
p.If(p.Cond(func() bool { return score.Get() >= 90 }, score),
    p.Html(` + "`" + `&lt;p>Grade: A&lt;/p>` + "`" + `),
).ElseIf(p.Cond(func() bool { return score.Get() >= 80 }, score),
    p.Html(` + "`" + `&lt;p>Grade: B&lt;/p>` + "`" + `),
).Else(
    p.Html(` + "`" + `&lt;p>Grade: F&lt;/p>` + "`" + `),
)

// p.Cond(predicateFn, ...dependencyStores)</code></pre>
		</section>

		<section id="lists">
			<h2>Lists</h2>
			<p>Reactive lists with <code>p.NewList()</code>. Render with <code>p.Each()</code>.</p>
			<pre><code>// Create a reactive list
items := p.NewList[string]("Apple", "Banana", "Cherry")

// Mutate — triggers re-render
items.Append("Date")
items.RemoveAt(0)
items.Set([]string{"Mango", "Papaya"})
items.Clear()

// Reactive length
items.Len()  // *Store[int] — updates automatically

// Render each item
p.Each(items, func(item string, i int) p.Node {
    return p.Html(` + "`" + `&lt;li>` + "`" + `, p.Itoa(i), ` + "`" + `: ` + "`" + `, item, ` + "`" + `&lt;/li>` + "`" + `)
}).Else(
    p.Html(` + "`" + `&lt;p>No items&lt;/p>` + "`" + `),
)</code></pre>
		</section>

		<section id="components">
			<h2>Components</h2>
			<p>Components are Go structs implementing <code>Render() p.Node</code>. Use <code>p.Comp()</code> to embed them.</p>
			<pre><code>// Define a component
type Badge struct {
    Label *p.Store[string]
}

func (b *Badge) Render() p.Node {
    return p.Html(` + "`" + `&lt;span class="badge">` + "`" + `, b.Label, ` + "`" + `&lt;/span>` + "`" + `)
}

// Scoped CSS
func (b *Badge) Style() string {
    return ` + "`" + `.badge{background:#007bff;color:#fff}` + "`" + `
}

// Use a component (props are struct fields)
p.Comp(&amp;Badge{Label: p.New("New")})

// Component with slot (child content)
p.Comp(&amp;Card{Title: p.New("Hello")},
    p.Html(` + "`" + `&lt;p>Slot content here&lt;/p>` + "`" + `),
)

// Inside the component, render slot content:
p.Slot()

// Callback props for component events
type Button struct {
    Label   *p.Store[string]
    OnClick func()
}

// Shared stores: pass the same *Store to multiple components
theme := p.New("light")
p.Comp(&amp;Header{Theme: theme})
p.Comp(&amp;Sidebar{Theme: theme})
// both components read/write the same store</code></pre>
		</section>

		<section id="routing">
			<h2>Routing</h2>
			<p>Client-side routing with <code>p.NewRouter()</code>. Each route maps to a component and an SSR HTML file.</p>
			<pre><code>// Define routes
routes := []p.Route{
    {Path: "/", HTMLFile: "index.html", SSRPath: "/", Component: home},
    {Path: "/about", HTMLFile: "about.html", SSRPath: "/about", Component: about},
}

// Start the router (in OnMount)
router := p.NewRouter(currentComponent, routes, "unique-id")
router.NotFound(func() { currentComponent.Set(nil) })
router.Start()

// Store[Component] holds the active route component
currentComponent *p.Store[p.Component]

// Links: client-side (default) vs server-side
// &lt;a href="/about">About&lt;/a>           (SPA navigation)
// &lt;a href="/about" external>About&lt;/a>  (full page reload)

// Store[Component] for local tabs (non-router):
activeTab := p.New[p.Component](tab1)
activeTab.WithOptions(tab1, tab2, tab3)
activeTab.Set(tab2)  // switches displayed component</code></pre>
		</section>

		<section id="fetch">
			<h2>Fetch</h2>
			<p>Type-safe HTTP requests with automatic JSON encoding/decoding via <code>js</code> struct tags.</p>
			<pre><code>// Define response type with js tags
type User struct {
    ID   int    ` + "`" + `js:"id"` + "`" + `
    Name string ` + "`" + `js:"name"` + "`" + `
}

// GET
go func() {
    user, err := p.Get[User]("https://api.example.com/user/1")
}()

// POST (send body, decode response)
go func() {
    created, err := p.Post[User](url, newUser)
}()

// PUT, PATCH, DELETE
result, err := p.Put[T](url, body)
result, err := p.Patch[T](url, body)
result, err := p.Delete[T](url)

// Advanced: custom headers and abort
signal, abort := p.NewAbortController()
go func() {
    result, err := p.Fetch[T](url, &amp;p.FetchOptions{
        Method:  "GET",
        Headers: map[string]string{"Authorization": "Bearer token"},
        Signal:  signal,
    })
}()
abort()  // cancel the request</code></pre>
		</section>

		<section id="storage">
			<h2>Storage</h2>
			<p>Persist state to localStorage. <code>LocalStore</code> auto-syncs on every <code>Set()</code>.</p>
			<pre><code>// Auto-persisted store (syncs on every Set)
theme := p.NewLocalStore("theme", "light")
theme.Set("dark")  // automatically saved to localStorage
theme.Store        // *Store[string] — use in Html() like any store

// Manual localStorage API
p.SetStorage("notes", "hello")
saved := p.GetStorage("notes")
p.RemoveStorage("notes")
p.ClearStorage()</code></pre>
		</section>

		<section id="timers">
			<h2>Timers</h2>
			<p>Debounce, throttle, setTimeout, and setInterval — all return cleanup functions.</p>
			<pre><code>// Debounce: fires after idle period
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
stop()  // stop the interval</code></pre>
		</section>

		<section id="lifecycle">
			<h2>Lifecycle</h2>
			<p>Components can implement lifecycle hooks for setup and teardown.</p>
			<pre><code>// OnMount: called when component becomes active
func (c *MyComp) OnMount() {
    if p.IsBuildTime {
        return  // skip side effects during SSR
    }
    // fetch data, start timers, etc.
}

// OnDestroy: called when component is removed (e.g. route change)
func (c *MyComp) OnDestroy() {
    if c.stopTimer != nil {
        c.stopTimer()
    }
}

// New: constructor — create stores, initialize state
func (c *MyComp) New() p.Component {
    return &amp;MyComp{
        Count: p.New(0),
        Name:  p.New(""),
    }
}

// GlobalStyle: CSS applied to entire page (on App component)
func (a *App) GlobalStyle() string { return ` + "`" + `body{margin:0}` + "`" + ` }

// Style: scoped CSS (auto-prefixed to this component)
func (c *MyComp) Style() string { return ` + "`" + `.btn{color:red}` + "`" + ` }

// p.IsBuildTime: true during SSR, false in WASM
// Use to guard browser-only code (fetch, timers, DOM access)</code></pre>
		</section>

	</div>
	</div>`)
}

func (m *Manual) Style() string {
	return `
.manual{padding:40px 0}
.manual h1{font-size:2.2em;color:#1a1a2e;margin-bottom:8px}
.intro{color:#666;margin-bottom:32px;font-size:1.05em}

.toc{margin-bottom:40px;padding:20px;background:#f8f9fa;border-radius:8px}
.toc h2{font-size:1em;margin-bottom:12px;color:#1a1a2e}
.toc ul{list-style:none;padding:0;margin:0;display:flex;flex-wrap:wrap;gap:8px}
.toc li a{display:inline-block;padding:6px 14px;background:#1a1a2e;color:#fff;border-radius:4px;font-size:13px;transition:background .2s}
.toc li a:hover{background:#0f3460}

.manual section{margin-bottom:32px;padding:24px;border:1px solid #e5e7eb;border-radius:8px;background:#fff}
.manual section h2{font-size:1.3em;color:#1a1a2e;margin-bottom:8px}
.manual section > p{color:#555;margin-bottom:16px;font-size:.95em}
.manual section code{background:#f1f5f9;padding:2px 6px;border-radius:3px;font-size:.85em}
.manual section pre{background:#1a1a2e;color:#e0e0e0;padding:16px;border-radius:6px;overflow-x:auto;font-size:13px;line-height:1.6}
.manual section pre code{background:transparent;padding:0;font-size:inherit}

@media(max-width:768px){
.toc ul{flex-direction:column}
}
`
}
