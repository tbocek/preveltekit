# PrevelteKit 2.0

Build reactive web apps in Go. Components compile to WebAssembly, with server-side pre-rendering for instant page loads.

## Features

- **Reactive stores** - `Store[T]`, `List[T]` with automatic DOM updates
- **Go DSL** - build UI trees directly in Go, no template language or code generation
- **Two-way binding** - `.Bind()` for text, number, and checkbox inputs
- **Event handling** - `.On("click", fn)`, `.PreventDefault()`, `.StopPropagation()`
- **Scoped CSS** - per-component styles with automatic class scoping
- **Client-side routing** - SPA navigation with path parameters
- **Typed fetch** - generic HTTP client with automatic JSON encoding/decoding
- **LocalStorage** - persistent stores that sync automatically
- **SSR + Hydration** - pre-rendered HTML at build time, hydrated with WASM at runtime

## Quick Start

```go
import p "github.com/tbocek/preveltekit/v2"
```

### Hello World

```go
type Hello struct{}

func (h *Hello) Render() p.Node {
    return p.Html(`<h1>Hello, World!</h1>`)
}
```

### Reactive Counter

Stores hold reactive state. Embed them in HTML and they update the DOM automatically.

```go
type Counter struct {
    Count *p.Store[int]
}

func (c *Counter) New() p.Component {
    return &Counter{Count: p.New(0)}
}

func (c *Counter) Render() p.Node {
    return p.Html(`<div>
        <p>Count: `, c.Count, `</p>`,
        p.Html(`<button>+1</button>`).On("click", func() {
            c.Count.Update(func(v int) int { return v + 1 })
        }),
    `</div>`)
}
```

`p.New(0)` creates a `*Store[int]` with initial value 0. Drop it into `Html()` and it becomes a live text node. `.On("click", fn)` wires up an event handler.

### Two-Way Binding

Bind a store to an input. Changes flow both ways.

```go
type Greeter struct {
    Name *p.Store[string]
}

func (g *Greeter) New() p.Component {
    return &Greeter{Name: p.New("")}
}

func (g *Greeter) Render() p.Node {
    return p.Html(`<div>
        <label>Name: `, p.Html(`<input type="text">`).Bind(g.Name), `</label>
        <p>Hello, `, g.Name, `!</p>
    </div>`)
}
```

`.Bind()` works with `*Store[string]`, `*Store[int]`, and `*Store[bool]` (checkbox).

### Conditionals

```go
score := p.New(75)

p.If(p.Cond(func() bool { return score.Get() >= 90 }, score),
    p.Html(`<p>Grade: A</p>`),
).ElseIf(p.Cond(func() bool { return score.Get() >= 70 }, score),
    p.Html(`<p>Grade: C</p>`),
).Else(
    p.Html(`<p>Grade: F</p>`),
)
```

`p.Cond(fn, ...stores)` pairs a boolean function with the stores it depends on so the framework knows when to re-evaluate.

### Lists

```go
type Todos struct {
    Items   *p.List[string]
    NewItem *p.Store[string]
}

func (t *Todos) New() p.Component {
    return &Todos{
        Items:   p.NewList[string]("Buy milk", "Write code"),
        NewItem: p.New(""),
    }
}

func (t *Todos) Add() {
    if item := t.NewItem.Get(); item != "" {
        t.Items.Append(item)
        t.NewItem.Set("")
    }
}

func (t *Todos) Render() p.Node {
    return p.Html(`<div>`,
        p.Html(`<input type="text">`).Bind(t.NewItem),
        p.Html(`<button>Add</button>`).On("click", t.Add),
        p.Html(`<ul>`,
            p.Each(t.Items, func(item string, i int) p.Node {
                return p.Html(`<li>`, item, `</li>`)
            }),
        `</ul>`),
    `</div>`)
}
```

`p.NewList` creates a reactive slice. `p.Each` renders each item. The list re-renders when items change.

### Components with Props and Slots

Define a reusable component:

```go
type Card struct {
    Title *p.Store[string]
}

func (c *Card) Render() p.Node {
    return p.Html(`<div class="card">
        <h2>`, c.Title, `</h2>
        <div>`, p.Slot(), `</div>
    </div>`)
}
```

Use it:

```go
p.Comp(&Card{Title: p.New("Welcome")},
    p.Html(`<p>This content fills the slot.</p>`),
)
```

Props are struct fields. `p.Slot()` renders child content passed to `p.Comp()`.

### Component Events (Callbacks)

Pass functions as props for child-to-parent communication:

```go
type Button struct {
    Label   *p.Store[string]
    OnClick func()
}

func (b *Button) Render() p.Node {
    return p.Html(`<button>`, b.Label, `</button>`).On("click", b.OnClick)
}

// parent usage:
p.Comp(&Button{Label: p.New("Save"), OnClick: func() {
    status.Set("Saved!")
}})
```

### Scoped CSS

Return CSS from `Style()` and it's automatically scoped to the component:

```go
func (c *Card) Style() string {
    return `.card { border: 1px solid #ddd; padding: 16px; border-radius: 8px; }`
}
```

No class name collisions across components.

### Conditional Attributes

```go
darkMode := p.New(false)

p.Html(`<div>`).AttrIf("class",
    p.Cond(func() bool { return darkMode.Get() }, darkMode),
    "dark",
)
```

When `darkMode` is true, the `dark` class is added. When false, it's removed.

### Derived Stores

Compute values from other stores:

```go
func Derived1[A, R any](a *p.Store[A], fn func(A) R) *p.Store[R] {
    out := p.New(fn(a.Get()))
    a.OnChange(func(_ A) { out.Set(fn(a.Get())) })
    return out
}

name := p.New("hello")
upper := Derived1(name, strings.ToUpper) // auto-updates when name changes
```

### Fetching Data

Typed HTTP client with automatic JSON encoding/decoding:

```go
type User struct {
    ID   int    `js:"id"`
    Name string `js:"name"`
}

func (c *MyComponent) OnMount() {
    if p.IsBuildTime {
        return // skip during SSR
    }
    go func() {
        user, err := p.Get[User]("/api/user/1")
        if err != nil {
            return
        }
        c.UserName.Set(user.Name)
    }()
}
```

Also available: `p.Post[T]`, `p.Put[T]`, `p.Patch[T]`, `p.Delete[T]`.

### Routing

```go
type App struct {
    CurrentPage *p.Store[p.Component]
}

func (a *App) Routes() []p.Route {
    return []p.Route{
        {Path: "/", HTMLFile: "index.html", SSRPath: "/", Component: &Home{}},
        {Path: "/about", HTMLFile: "about.html", SSRPath: "/about", Component: &About{}},
    }
}

func (a *App) OnMount() {
    router := p.NewRouter(a.CurrentPage, a.Routes(), "app")
    router.Start()
}

func (a *App) Render() p.Node {
    return p.Html(`<div>
        <nav>
            <a href="/">Home</a>
            <a href="/about">About</a>
        </nav>
        <main>`, a.CurrentPage, `</main>
    </div>`)
}
```

Internal `<a>` links are automatically intercepted for SPA navigation. Add the `external` attribute to opt out.

### LocalStorage

```go
// auto-persists on every .Set()
theme := p.NewLocalStore("theme", "light")
theme.Set("dark") // saved to localStorage immediately

// manual localStorage API
p.SetStorage("key", "value")
val := p.GetStorage("key")
p.RemoveStorage("key")
```

### Lifecycle

| Interface | Method | When |
|-----------|--------|------|
| `HasNew` | `New() Component` | Factory -- create stores and child components here |
| `HasOnMount` | `OnMount()` | Component becomes active (fetch data, start timers) |
| `HasOnDestroy` | `OnDestroy()` | Component removed (cleanup) |
| `HasStyle` | `Style() string` | Scoped CSS for this component |
| `HasGlobalStyle` | `GlobalStyle() string` | Global CSS (unscoped) |

### Timers

```go
stop := p.SetInterval(1000, func() { /* runs every second */ })
defer stop()

cancel := p.SetTimeout(3000, func() { /* runs once after 3s */ })

debounced, cleanup := p.Debounce(300, handler)
defer cleanup()
```

## Build

Output goes to `dist/` -- serve with any static file server.

```
dist/
  index.html     # pre-rendered HTML
  main.wasm      # compiled WASM binary
  wasm_exec.js   # Go WASM runtime
```

## Architecture

Both SSR (native Go at build time) and WASM (browser at runtime) execute the same component code. SSR pre-renders HTML with comment markers and element IDs. WASM walks the same `Render()` tree to discover bindings and wire them to the existing DOM. No intermediate binary format, no code generation -- just a direct tree walk.

The critical invariant: SSR and WASM must create stores and register handlers in identical order so counter-based IDs match between pre-rendered HTML and the live WASM runtime.

---

## History

PrevelteKit went through several architectural stages on the way to 2.0. The core philosophy stayed the same throughout: minimal framework, static HTML output, clear separation between frontend and backend.

### 1.x -- Svelte/TypeScript

The original PrevelteKit was a minimalistic (~500 LoC) web framework built on [Svelte 5](https://svelte.dev/), using [Rsbuild](https://rsbuild.dev/) as the bundler and [jsdom](https://github.com/jsdom/jsdom) for build-time pre-rendering. Components were standard Svelte files:

```svelte
<script>
    let count = $state(0);
</script>

<h1>Count: {count}</h1>
<button onclick={() => count++}>Click me</button>
```

The motivation was simple: SvelteKit is powerful but heavy. PrevelteKit offered build-time pre-rendering without the meta-framework complexity. The output was purely static assets -- HTML, CSS, JS -- deployable to any CDN or web server with no server runtime required.

This version worked well, but the dependency on the JavaScript ecosystem (Node.js, npm, bundlers) remained a friction point. The idea of writing the entire frontend in Go and compiling to WASM started to take shape.

### 1.9.1 -- Go/WASM with DSL and Code Generation

The first Go rewrite introduced a Svelte-inspired template DSL. Components had a `Template()` method returning a string with special syntax:

```go
func (c *Counter) Template() string {
    return `<div>
        <p>Count: {Count}</p>
        <button @click="Increment()">+1</button>
        {#if Count > 10}
            <p>That's a lot!</p>
        {:else}
            <p>Keep clicking</p>
        {/if}
    </div>`
}
```

A build step parsed these templates and generated Go code -- transforming `{Count}` into store reads, `@click` into handler registrations, `{#if}` / `{#each}` into conditional/iteration logic. The generated code was then compiled to WASM.

This approach worked but had significant drawbacks:
- A custom parser and code generator added complexity and maintenance burden
- Template errors surfaced at generation time, not compile time -- debugging was indirect
- The generated Go code was hard to read and harder to debug
- Two languages in one file (Go + template DSL) felt awkward

### 1.9.2 -- Bindings Binary

The next iteration removed the template DSL in favor of writing UI trees directly in Go. But it introduced a different separation: SSR rendered HTML at build time, and a `bindings.bin` file was generated to tell the WASM runtime where all the reactive bindings, event handlers, and dynamic blocks lived. The WASM binary didn't contain any HTML -- it only read the bindings file and wired up interactivity.

This reduced WASM binary size since no HTML strings were compiled in, but added its own complexity:
- A custom binary format had to be designed, serialized at build time, and deserialized at runtime
- The bindings file was another artifact to generate, serve, and keep in sync
- Any mismatch between the HTML and the bindings file caused subtle, hard-to-diagnose bugs
- The indirection made the system harder to reason about

Many intermediate prototypes were built and discarded between 1.9.1 and 1.9.2 (and between 1.9.2 and 2.0), often with the help of LLMs for rapid exploration of different approaches.

### 2.0 -- Direct Tree Walk

The current version eliminates both code generation and the bindings binary. Components define their UI with a Go DSL using `Html()`, `If()`, `Each()`, `Comp()`, etc. The same `Render()` method runs at build time (native Go, SSR) and at runtime (WASM, hydration). Both walks advance the same global counters in the same order, so comment markers and element IDs match without any intermediate format.

What changed:
- **No code generation** -- the Go DSL is plain Go, checked by the compiler
- **No bindings.bin** -- WASM discovers bindings by walking the same tree SSR walked
- **No template language** -- conditionals, loops, and components are Go function calls
- **Simpler mental model** -- one `Render()` method, two execution contexts

The tradeoff is that WASM binaries include HTML string literals, making them slightly larger. In practice (~60kb gzipped) this is acceptable.

What started as ~500 lines of glue code between Svelte, jsdom, and Rsbuild is now a ~4k LoC self-contained framework with no external dependencies beyond the Go standard library and the WASM runtime.

### Why Static Output?

Throughout all versions, the preference has been clear separation: the frontend is static assets (HTML/CSS/JS or HTML/CSS/WASM) served from any CDN. The backend is a separate service with `/api` endpoints. No server-side rendering runtime, no Node.js in production, no blurred boundaries between view code and server code.

Meta-frameworks like Next.js, Nuxt, and SvelteKit blur this separation by requiring a JavaScript runtime for SSR, API routes, and build-time generation. Serving just static content is simpler: deploy anywhere (GitHub Pages, S3, any web server) with predictable performance.

**Classic SSR** (Next.js, Nuxt): Server renders HTML on every request. Requires a runtime.

**SPA** (React, Vue): Browser renders everything. User sees a blank page until JS loads.

**Build-time Pre-rendering** (PrevelteKit): HTML is rendered once at build time. User sees content instantly. WASM hydrates for interactivity. No server runtime needed.

### Inspiration

- https://github.com/serge-hulne/Golid
- https://github.com/maxence-charriere/go-app
- https://github.com/hexops/vecty
- https://github.com/vugu/vugu

## License

MIT
