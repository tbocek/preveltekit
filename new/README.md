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
