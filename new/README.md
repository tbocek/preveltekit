# PrevelteKit 1.0

A Go framework for building web applications that compile to WebAssembly. Write components in Go with a Svelte-like template syntax, get server-side rendering and client-side hydration.

## Features

- **Reactive stores** - `Store[T]`, `List[T]`, `Map[K,V]` with automatic DOM updates
- **Component templates** - Svelte-inspired syntax with `{#if}`, `{#each}`, `{:else}`
- **Two-way binding** - `bind:value`, `bind:checked` for form inputs
- **Event handling** - `@click`, `@submit.preventDefault`
- **Client-side routing** - SPA navigation with path parameters
- **Typed fetch** - Generic HTTP client with automatic JSON encoding/decoding
- **LocalStorage** - Persistent stores that sync automatically
- **SSR + Hydration** - Pre-rendered HTML, then hydrated with WASM

## Quick Start

```go
package main

import "preveltekit"

type Counter struct {
    Count *preveltekit.Store[int]
}

func (c *Counter) Increment() {
    c.Count.Update(func(v int) int { return v + 1 })
}

func (c *Counter) Template() string {
    return `<div>
        <p>Count: {Count}</p>
        <button @click="Increment()">+1</button>
    </div>`
}
```

Build with:

```bash
# Development (with file watching and auto-rebuild)
./dev.sh myapp/counter.go

# Production build
./build.sh --release myapp/counter.go
```

Output goes to `myapp/dist/` - serve it with any static file server.

## Template Syntax

### Expressions

```html
<p>{Name}</p>
<p>Score: {Score}</p>
```

### Conditionals

```html
{#if Score >= 90}
    <p>Grade: A</p>
{:else if Score >= 80}
    <p>Grade: B</p>
{:else}
    <p>Grade: F</p>
{/if}
```

### Lists

```html
{#each Items as item, i}
    <li>{i}: {item}</li>
{:else}
    <li>No items</li>
{/each}
```

### Two-Way Binding

```html
<input type="text" bind:value="Name">
<input type="checkbox" bind:checked="Agreed">
```

### Events

```html
<button @click="Save()">Save</button>
<button @click="Add(5)">Add 5</button>
<form @submit.preventDefault="Submit()">
```

### Dynamic Classes

```html
<div class:active={IsActive}>
<div class:dark={DarkMode}>
```

### Dynamic Attributes

```html
<div data-state="{State}">
<a href="/user/{UserID}">
```

### Components

```html
<Card title="Hello">
    <p>Content goes in the slot</p>
</Card>

<Button label="Click" @click="HandleClick()" />
```

## Stores

### Store[T]

Basic reactive value:

```go
type App struct {
    Name  *preveltekit.Store[string]
    Count *preveltekit.Store[int]
}

// In methods:
a.Name.Set("Alice")
a.Count.Update(func(v int) int { return v + 1 })
name := a.Name.Get()
```

### List[T]

Reactive list with efficient diff updates:

```go
type App struct {
    Items *preveltekit.List[string]
}

a.Items.Append("new item")
a.Items.RemoveAt(0)
a.Items.Set([]string{"a", "b", "c"})  // computes minimal diff
a.Items.Clear()
```

### LocalStore

String store that persists to localStorage:

```go
type App struct {
    Theme *preveltekit.LocalStore  // auto-saved with key "Theme"
}
```

## Routing

```go
func (a *App) OnMount() {
    router := preveltekit.NewRouter()

    router.Handle("/", func(p map[string]string) {
        a.Page.Set("home")
    })
    router.Handle("/user/:id", func(p map[string]string) {
        a.UserID.Set(p["id"])
    })
    router.NotFound(func() {
        a.Page.Set("404")
    })

    router.Start()
}
```

Links are automatically intercepted for SPA navigation. Use the `external` attribute for full page loads:

```html
<a href="/about">SPA navigation</a>
<a href="/about" external>Full page load</a>
```

## Fetch

Typed HTTP client (must be called from a goroutine):

```go
type User struct {
    ID   int    `js:"id"`
    Name string `js:"name"`
}

func (a *App) LoadUser() {
    go func() {
        user, err := preveltekit.Get[User]("/api/user/1")
        if err != nil {
            a.Error.Set(err.Error())
            return
        }
        a.User.Set(user.Name)
    }()
}
```

Available methods: `Get`, `Post`, `Put`, `Patch`, `Delete`.

## Build

Requirements:
- Go 1.21+
- TinyGo
- wasm-strip (from wabt)

```bash
# Development build
./build.sh app/main.go

# Release build (smaller output)
./build.sh --release app/main.go

# Multiple components
./build.sh app/main.go app/header.go app/sidebar.go
```

The build process:
1. Transforms component syntax to valid Go
2. Pre-renders HTML (SSR)
3. Compiles to WASM with TinyGo
4. Tree-shakes the JS runtime
5. Assembles final output with gzip/brotli compression

## Project Structure

```
myapp/
  app.go           # main component
  header.go        # child component
  build/           # generated Go code (gitignore this)
  dist/            # output files
    index.html
    app.wasm
    wasm_exec.js
```

## Component Lifecycle

```go
type App struct {
    // Stores are auto-initialized
    Count *preveltekit.Store[int]
}

func (a *App) OnMount() {
    // Called after component is mounted to DOM
    // Set initial values, start timers, fetch data
    a.Count.Set(0)
}

func (a *App) Template() string {
    return `<div>...</div>`
}

func (a *App) Style() string {
    return `.app { ... }`  // injected once per component type
}
```

## License

MIT

---

## History: PrevelteKit Origins (Svelte/TypeScript Version)

PrevelteKit 1.0 (Go/WASM) is a complete rewrite. The original PrevelteKit was a minimalistic (>500 LoC) web application framework built on [Svelte 5](https://svelte.dev/), featuring single page applications with build-time pre-rendering using [Rsbuild](https://rsbuild.dev/) as the build/bundler tool and [jsdom](https://github.com/jsdom/jsdom) as the DOM environment for pre-rendering components during build.

The inspiration for the original project came from the Vue SSR example in the [Rspack examples repository](https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs). The project adapted those concepts for Svelte, providing a minimal setup.

### Original Motivation

While SvelteKit is the go-to solution for SSR with Svelte, the original PrevelteKit provided a minimalistic solution for build-time pre-rendering without the additional complexity.

From an architectural standpoint, the preference was for clear separation between view code and server code, where the frontend requests data from the backend via dedicated `/api` endpoints. This treats the frontend as purely static assets (HTML/CSS/JS) that can be served from any CDN or simple web server.

Meta-frameworks such as Next.js, Nuxt.js, and SvelteKit blur this separation by requiring a JavaScript runtime (Node.js, Deno, or Bun) for server-side rendering, API routes, and build-time generation. While platforms like Vercel and Netlify can help handle this complex setup, serving just static content is much simpler: deploy anywhere (GitHub Pages, S3, any web server) with predictable performance. You avoid the "full-stack JavaScript" complexity for your deployed frontend - it's just files on a server, nothing more.

### Why Not SvelteKit + adapter-static?

While SvelteKit with adapter-static can achieve similar static site generation, the original PrevelteKit offered a minimalistic alternative using Svelte + jsdom + Rsbuild. At less than 500 lines of code, it was essentially glue code between these libraries rather than a full framework. This provided a lightweight solution for those who wanted static pre-rendering without SvelteKit's additional complexity and features.

### Why Rsbuild and not Vite?

While [benchmarks](https://github.com/rspack-contrib/build-tools-performance) show that Rsbuild and Vite (Rolldown + Oxc) have comparable overall performance in many cases (not for the 10k component case), Rsbuild had a small advantage in producing the smallest compressed bundle size, while Vite (Rolldown + Oxc) had a small advantage in build time performance.

In practice, Rsbuild "just works" after many updates out of the box with minimal configuration, which reduced friction and setup time. However, Vite (Rolldown + Oxc) was being watched closely as it progressed fast.

### Original Key Features

- Lightning Fast: Rsbuild bundles in the range of a couple hundred milliseconds
- Simple Routing: Built-in routing system  
- Layout and static content pre-rendered with Svelte and hydration
- Zero Config: Works out of the box with sensible defaults
- Hot reload in development, production-ready in minutes
- Docker-based development environments to protect against supply chain attacks

### Automatic Fetch Handling

The original PrevelteKit automatically managed fetch requests during build-time pre-rendering:

- Components render with loading states in the pre-rendered HTML
- No need to wrap fetch calls in `window.__isBuildTime` checks
- Use Svelte's `{#await}` blocks for clean loading/error/success states
- If anything went missing, in the worst case, fetch calls timeout after 5 seconds during pre-rendering

### Rendering Comparison

The original project compared three rendering approaches:

**SSR (classic SSR / Next.js / Nuxt)**
- Initial Load: User sees fully rendered content instantly
- After Script Execution: Content remains the same, scripts add interactivity

**SPA (React App / pure Svelte)**
- Initial Load: User sees blank page or loading spinner
- After Script Execution: User sees full interactive content

**SPA + Build-time Pre-Rendering (PrevelteKit approach)**
- Initial Load: User sees pre-rendered static content
- After Script Execution: Content becomes fully interactive

The original repository included SVG diagrams (SSR.svg, SPA.svg, SPAwBR.svg) illustrating these differences visually.

### Prerequisites (Original)

- Node.js (Latest LTS version recommended)
- npm/pnpm or similar

### Original Quick Start

```bash
# Create test directory and go into this directory
mkdir -p preveltekit/src && cd preveltekit 

# Declare dependency and the dev script
echo '{"devDependencies": {"preveltekit": "^1.2.25"},"dependencies": {"svelte": "^5.39.11"},"scripts": {"dev": "preveltekit dev"}}' > package.json 

# Download dependencies
npm install 

# A very simple svelte file
echo '<script>let count = $state(0);</script><h1>Count: {count}</h1><button onclick={() => count++}>Click me</button>' > src/Index.svelte 

# And open a browser with localhost:3000
npm run dev 
```

### Slow Start (Original)

One example was within the project in the example folder, and another example was the [notary example](https://github.com/tbocek/notary-example). The CLI supported: dev/stage/prod.

**Development Server**

```bash
npm run dev
```

This started an Express development server on http://localhost:3000, with:
- Live reloading
- No optimization for faster builds
- Ideal for rapid development

**Build for Production**

```bash
npm run build
```

The production build:
- Generated pre-compressed static files for optimal serving with best compression:
  - Brotli (`.br` files)
  - Zstandard (`.zst` files)
  - Zopfli (`.gz` files)
- Optimized assets for production

**Staging Environment**

```bash
npm run stage
```

The development server prioritized fast rebuilds and developer experience, while the production build focused on optimization and performance. Testing with a stage and production build before deploying was recommended.

### Original Docker Support

To build with docker in production mode:

```bash
docker build . -t preveltekit
docker run -p3000:3000 preveltekit
```

To run in development mode with live reloading:

```bash
docker build -f Dockerfile.dev . -t preveltekit-dev
docker run -p3000:3000 -v./src:/app/src preveltekit-dev
```

### Configuration (Original)

The original PrevelteKit used rsbuild.config.ts for configuration with sensible defaults. To customize settings, you would create an rsbuild.config.ts file in your project - it would merge with the default configuration.

The framework provided fallback files (index.html and index.ts) from the default folder when you didn't supply your own. Once you added your own index.html or index.ts files, PrevelteKit used those instead, ignoring the defaults.

This approach followed a "convention over configuration" pattern where you only needed to specify what differed from the defaults.

### Why the Rewrite?

PrevelteKit 1.0 moves from TypeScript/Svelte to Go/WebAssembly while keeping the same philosophy: minimal framework, static output, clear separation of concerns. The Go version compiles to WASM for client-side execution while still supporting build-time pre-rendering.
