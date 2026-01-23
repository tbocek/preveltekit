package main

// Docs is the documentation page
type Docs struct{}

func (d *Docs) Template() string {
	return `<div class="page doc-page">
	<h1>Documentation</h1>

	<nav class="doc-nav">
		<a href="#getting-started">Getting Started</a>
		<a href="#stores">Reactive Stores</a>
		<a href="#components">Components</a>
		<a href="#routing">Routing</a>
		<a href="#ssr">SSR</a>
	</nav>

	<div class="doc-content">
	<section id="getting-started" class="doc-section">
		<h2>Getting Started</h2>
		<p>Reactive is a Go framework for building web applications with WebAssembly. It provides reactive state management, a component system, and client-side routing.</p>

		<h3>Installation</h3>
		<pre><code>go get github.com/user/reactive</code></pre>

		<h3>Requirements</h3>
		<ul>
			<li>Go 1.21 or later</li>
			<li>TinyGo (for small WASM builds)</li>
		</ul>
	</section>

	<section id="stores" class="doc-section">
		<h2>Reactive Stores</h2>
		<p>Stores are the core of Reactive's state management. They hold values and notify subscribers when changed.</p>

		<h3>Store[T]</h3>
		<pre><code>// Create a store
count := reactive.New(0)

// Get value
fmt.Println(count.Get()) // 0

// Set value (triggers updates)
count.Set(5)

// Update with function
count.Update(func(v int) int { return v + 1 })

// Subscribe to changes
count.OnChange(func(v int) {
    fmt.Println("Count is now:", v)
})</code></pre>

		<h3>List[T]</h3>
		<pre><code>// Create a list
items := reactive.NewList("a", "b", "c")

// Append items
items.Append("d", "e")

// Remove by index
items.RemoveAt(0)

// Subscribe to edits
items.OnEdit(func(edit reactive.Edit[string]) {
    if edit.Op == reactive.EditInsert {
        // Handle insert
    }
})</code></pre>

		<h3>Map[K, V]</h3>
		<pre><code>// Create a map
users := reactive.NewMap[string, User]()

// Set values
users.Set("alice", User{Name: "Alice"})

// Get values
user, ok := users.Get("alice")

// Delete
users.Delete("alice")</code></pre>
	</section>

	<section id="components" class="doc-section">
		<h2>Components</h2>
		<p>Components are Go structs with Template() and optional Style() methods.</p>

		<h3>Basic Component</h3>
		<pre><code>type Counter struct {
    Count *reactive.Store[int]
}

func (c *Counter) OnMount() {
    c.Count.Set(0)
}

func (c *Counter) Increment() {
    c.Count.Update(func(v int) int { return v + 1 })
}

func (c *Counter) Template() string {
    return ` + "`" + `<div>
    <p>Count: {Count}</p>
    <button @click="Increment()">+1</button>
</div>` + "`" + `
}</code></pre>

		<h3>Props</h3>
		<pre><code><Button label="Click me" variant="primary" /></code></pre>

		<h3>Slots</h3>
		<pre><code><Card title="My Card">
    <p>This content goes in the slot</p>
</Card></code></pre>

		<h3>Events</h3>
		<pre><code><Button @click="HandleClick()" /></code></pre>
	</section>

	<section id="routing" class="doc-section">
		<h2>Routing</h2>
		<p>Client-side routing with history API support.</p>

		<h3>Basic Router</h3>
		<pre><code>router := reactive.NewRouter()

router.Handle("/", func(params map[string]string) {
    // Show home page
})

router.Handle("/user/:id", func(params map[string]string) {
    userID := params["id"]
    // Show user page
})

router.NotFound(func() {
    // Show 404
})

router.Start()</code></pre>

		<h3>Navigation</h3>
		<pre><code>// Programmatic navigation
router.Navigate("/about")

// Replace without history
router.Replace("/login")

// Get current path
path := router.CurrentPath().Get()</code></pre>

		<h3>Navigation Guards</h3>
		<pre><code>router.BeforeNavigate(func(from, to string) bool {
    if !isLoggedIn && to == "/admin" {
        router.Replace("/login")
        return false // cancel navigation
    }
    return true
})</code></pre>
	</section>

	<section id="ssr" class="doc-section">
		<h2>Server-Side Rendering</h2>
		<p>Reactive supports SSR for fast initial page loads and SEO.</p>

		<h3>Build Process</h3>
		<ol>
			<li>Generate code with reactivebuild</li>
			<li>Pre-render HTML with native Go binary</li>
			<li>Build WASM for client hydration</li>
			<li>Assemble final HTML with embedded content</li>
		</ol>

		<pre><code># Build everything
./build.sh app.go

# Output structure
dist/
  index.html    # Pre-rendered HTML
  app.wasm      # Client WASM
  wasm_exec.js  # WASM loader</code></pre>
	</section>
	</div>
</div>`
}

func (d *Docs) Style() string {
	return `
.doc-page { display: grid; grid-template-columns: 200px 1fr; gap: 2rem; }
.doc-page h1 { grid-column: 1 / -1; margin-bottom: 1rem; color: #1a1a2e; }
.doc-nav { position: sticky; top: 1rem; align-self: start; display: flex; flex-direction: column; gap: 0.5rem; padding: 1rem; background: #f8f9fa; border-radius: 8px; }
.doc-nav a { color: #666; text-decoration: none; padding: 0.5rem; border-radius: 4px; transition: all 0.2s; }
.doc-nav a:hover { color: #1a1a2e; background: #e9ecef; }
.doc-section { margin-bottom: 3rem; }
.doc-section h2 { color: #1a1a2e; border-bottom: 2px solid #e9ecef; padding-bottom: 0.5rem; margin-bottom: 1rem; }
.doc-section h3 { color: #444; margin: 1.5rem 0 0.75rem; }
.doc-section p { color: #555; margin-bottom: 1rem; }
.doc-section ul, .doc-section ol { margin: 1rem 0; padding-left: 1.5rem; color: #555; }
.doc-section li { margin-bottom: 0.5rem; }
.doc-section pre { background: #1a1a2e; color: #e9ecef; padding: 1rem; border-radius: 8px; overflow-x: auto; margin: 1rem 0; }
.doc-section code { font-family: 'Fira Code', monospace; font-size: 0.9rem; }
@media (max-width: 768px) { .doc-page { grid-template-columns: 1fr; } .doc-nav { position: static; flex-direction: row; flex-wrap: wrap; } }
`
}
