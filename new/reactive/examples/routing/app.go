package main

import "reactive"

// App is the main router component
type App struct {
	CurrentPage *reactive.Store[string]
	router      *reactive.Router
}

func (a *App) OnMount() {
	a.CurrentPage.Set("index")

	// Log when CurrentPage changes
	a.CurrentPage.OnChange(func(page string) {
		println("[App] CurrentPage changed to:", page)
	})

	a.router = reactive.NewRouter()

	a.router.Handle("/", func(params map[string]string) {
		println("[App] Route / matched, setting page to index")
		a.CurrentPage.Set("index")
	})

	a.router.Handle("/index", func(params map[string]string) {
		println("[App] Route /index matched, setting page to index")
		a.CurrentPage.Set("index")
	})

	a.router.Handle("/landing", func(params map[string]string) {
		println("[App] Route /landing matched, setting page to landing")
		a.CurrentPage.Set("landing")
	})

	a.router.Handle("/doc", func(params map[string]string) {
		println("[App] Route /doc matched, setting page to doc")
		a.CurrentPage.Set("doc")
	})

	a.router.Handle("/example", func(params map[string]string) {
		println("[App] Route /example matched, setting page to example")
		a.CurrentPage.Set("example")
	})

	a.router.NotFound(func() {
		println("[App] No route matched, setting page to notfound")
		a.CurrentPage.Set("notfound")
	})

	a.router.Start()
}

func (a *App) Navigate(page string) {
	a.router.Navigate("/" + page)
}

func (a *App) GoHome() {
	a.router.Navigate("/")
}

func (a *App) GoLanding() {
	a.router.Navigate("/landing")
}

func (a *App) GoDoc() {
	a.router.Navigate("/doc")
}

func (a *App) GoExample() {
	a.router.Navigate("/example")
}

func (a *App) Template() string {
	return `<div class="app">
	<nav class="nav">
		<a href="/" class="nav-brand">Reactive</a>
		<div class="nav-links">
			<a href="/index" external class="nav-link">Home</a>
			<a href="/landing" external class="nav-link">Landing</a>
			<a href="/doc" external class="nav-link">Docs</a>
			<a href="/example" external class="nav-link">Examples</a>
		</div>
		<span class="routing-indicator">(Server-side routing - page reloads)</span>
	</nav>

	<main class="content">
		{#if CurrentPage == "index"}
			<Index />
		{:else if CurrentPage == "landing"}
			<Landing />
		{:else if CurrentPage == "doc"}
			<Docs />
		{:else if CurrentPage == "example"}
			<Example />
		{:else}
			<NotFound />
		{/if}
	</main>

	<footer class="footer">
		<p>Built with Reactive - Go + WebAssembly</p>
	</footer>
</div>`
}

func (a *App) Style() string {
	return `
* { box-sizing: border-box; margin: 0; padding: 0; }
body { font-family: system-ui, -apple-system, sans-serif; line-height: 1.6; color: #333; }
.app { min-height: 100vh; display: flex; flex-direction: column; }
.nav { background: #1a1a2e; padding: 1rem 2rem; display: flex; align-items: center; gap: 2rem; }
.nav-brand { color: #fff; font-size: 1.5rem; font-weight: bold; text-decoration: none; }
.nav-links { display: flex; gap: 1.5rem; }
.nav-link { color: #a0a0a0; text-decoration: none; transition: color 0.2s; }
.nav-link:hover { color: #fff; }
.routing-indicator { color: #888; font-size: 0.75rem; font-style: italic; margin-left: auto; }
.content { flex: 1; padding: 2rem; max-width: 1200px; margin: 0 auto; width: 100%; }
.footer { background: #f5f5f5; padding: 1rem 2rem; text-align: center; color: #666; }
`
}
