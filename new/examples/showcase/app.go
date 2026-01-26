package main

import "preveltekit"

type App struct {
	CurrentPage *preveltekit.Store[string]
}

func (a *App) OnMount() {
	router := preveltekit.NewRouter()
	router.RegisterRoutes(a.Routes())
	router.NotFound(func() {
		a.CurrentPage.Set("notfound")
	})
	router.Start()
}

func (a *App) Routes() []preveltekit.StaticRoute {
	return []preveltekit.StaticRoute{
		{Path: "/", HTMLFile: "index.html", Handler: func(p map[string]string) { a.CurrentPage.Set("basics") }},
		{Path: "/basics", HTMLFile: "basics.html", Handler: func(p map[string]string) { a.CurrentPage.Set("basics") }},
		{Path: "/components", HTMLFile: "components.html", Handler: func(p map[string]string) { a.CurrentPage.Set("components") }},
		{Path: "/lists", HTMLFile: "lists.html", Handler: func(p map[string]string) { a.CurrentPage.Set("lists") }},
		{Path: "/routing", HTMLFile: "routing.html", Handler: func(p map[string]string) { a.CurrentPage.Set("routing") }},
		{Path: "/links", HTMLFile: "links.html", Handler: func(p map[string]string) { a.CurrentPage.Set("links") }},
		{Path: "/fetch", HTMLFile: "fetch.html", Handler: func(p map[string]string) { a.CurrentPage.Set("fetch") }},
		{Path: "/storage", HTMLFile: "storage.html", Handler: func(p map[string]string) { a.CurrentPage.Set("storage") }},
		{Path: "/debounce", HTMLFile: "debounce.html", Handler: func(p map[string]string) { a.CurrentPage.Set("debounce") }},
		{Path: "/bitcoin", HTMLFile: "bitcoin.html", Handler: func(p map[string]string) { a.CurrentPage.Set("bitcoin") }},
	}
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
*{box-sizing:border-box}
body{margin:0;font-family:system-ui,-apple-system,sans-serif}
.showcase{display:flex;min-height:100vh}
.sidebar{width:200px;background:#1a1a2e;color:#fff;padding:20px;flex-shrink:0}
.sidebar h2{margin:0 0 20px;font-size:1.5em}
.sidebar ul{list-style:none;padding:0;margin:0}
.sidebar li{margin:5px 0}
.sidebar a{display:block;padding:10px 15px;color:#ccc;text-decoration:none;border-radius:4px;font-size:14px}
.sidebar a:hover{background:#2a2a4e;color:#fff}
.content{flex:1;padding:20px;background:#f5f5f5;overflow-y:auto}
.demo{max-width:600px}
.demo h1{color:#1a1a2e;margin-bottom:20px}
.demo h2{margin-top:0;color:#666;font-size:1.1em}
.demo section{margin:20px 0;padding:15px;border:1px solid #ddd;border-radius:8px;background:#fff}
.demo button{padding:8px 16px;margin:4px;cursor:pointer;border:1px solid #ccc;border-radius:4px;background:#f5f5f5}
.demo button:hover{background:#e5e5e5}
.demo input[type=text]{padding:8px;border:1px solid #ccc;border-radius:4px}
.demo pre{background:#f5f5f5;padding:15px;border-radius:4px;overflow-x:auto;font-size:12px;white-space:pre-wrap}
.demo textarea{width:100%;padding:10px;border:1px solid #ccc;border-radius:4px;font-family:inherit;resize:vertical}
.hint{font-size:.9em;color:#666;font-style:italic}
.buttons{display:flex;gap:10px;flex-wrap:wrap;margin:10px 0}
`
}
