package main

import "preveltekit"

type App struct {
	CurrentPage *preveltekit.Store[string]
}

func (a *App) OnMount() {
	router := preveltekit.NewRouter()

	router.Handle("/", func(p map[string]string) {
		a.CurrentPage.Set("basics")
	})
	router.Handle("/basics", func(p map[string]string) {
		a.CurrentPage.Set("basics")
	})
	router.Handle("/components", func(p map[string]string) {
		a.CurrentPage.Set("components")
	})
	router.Handle("/lists", func(p map[string]string) {
		a.CurrentPage.Set("lists")
	})
	router.Handle("/routing", func(p map[string]string) {
		a.CurrentPage.Set("routing")
	})
	router.Handle("/links", func(p map[string]string) {
		a.CurrentPage.Set("links")
	})
	router.Handle("/fetch", func(p map[string]string) {
		a.CurrentPage.Set("fetch")
	})
	router.Handle("/storage", func(p map[string]string) {
		a.CurrentPage.Set("storage")
	})
	router.Handle("/debounce", func(p map[string]string) {
		a.CurrentPage.Set("debounce")
	})
	router.Handle("/bitcoin", func(p map[string]string) {
		a.CurrentPage.Set("bitcoin")
	})

	router.NotFound(func() {
		a.CurrentPage.Set("notfound")
	})

	router.Start()
}

func (a *App) Template() string {
	return `<div class="showcase">
	<nav class="sidebar">
		<h2>Reactive</h2>
		<ul>
			<li><a href="/basics">Basics</a></li>
			<li><a href="/components">Components</a></li>
			<li><a href="/lists">Lists</a></li>
			<li><a href="/routing">Routing</a></li>
			<li><a href="/links">Links</a></li>
			<li><a href="/fetch">Fetch</a></li>
			<li><a href="/storage">Storage</a></li>
			<li><a href="/debounce">Debounce</a></li>
			<li><a href="/bitcoin">Bitcoin</a></li>
		</ul>
	</nav>
	<main class="content">
		{#if CurrentPage == "basics"}
			<Basics />
		{:else if CurrentPage == "components"}
			<Components />
		{:else if CurrentPage == "lists"}
			<Lists />
		{:else if CurrentPage == "routing"}
			<Routing />
		{:else if CurrentPage == "links"}
			<Links />
		{:else if CurrentPage == "fetch"}
			<Fetch />
		{:else if CurrentPage == "storage"}
			<Storage />
		{:else if CurrentPage == "debounce"}
			<Debounce />
		{:else if CurrentPage == "bitcoin"}
			<Bitcoin />
		{:else if CurrentPage == "notfound"}
			<div class="notfound">
				<h1>404</h1>
				<p>Page not found</p>
				<a href="/">Go Home</a>
			</div>
		{/if}
	</main>
</div>`
}

func (a *App) Style() string {
	return `
* { box-sizing: border-box; }
body { margin: 0; font-family: system-ui, -apple-system, sans-serif; }
.showcase { display: flex; min-height: 100vh; }
.sidebar { width: 200px; background: #1a1a2e; color: #fff; padding: 20px; flex-shrink: 0; }
.sidebar h2 { margin: 0 0 20px; font-size: 1.5em; color: #fff; }
.sidebar ul { list-style: none; padding: 0; margin: 0; }
.sidebar li { margin: 5px 0; }
.sidebar a { display: block; padding: 10px 15px; color: #ccc; text-decoration: none; border-radius: 4px; font-size: 14px; }
.sidebar a:hover { background: #2a2a4e; color: #fff; }
.content { flex: 1; padding: 20px; background: #f5f5f5; overflow-y: auto; }
`
}
