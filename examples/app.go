package main

import p "github.com/tbocek/preveltekit"

type App struct {
	CurrentComponent *p.Store[p.Component]
	routes           []p.Route
}

func (a *App) New() p.Component {
	// Initialize child components
	basics := (&Basics{}).New()
	derived := (&Derived{}).New()
	complex := (&Complex{}).New()
	components := (&Components{}).New()
	lists := (&Lists{}).New()
	routing := (&Routing{}).New()
	links := (&Links{}).New()
	fetch := (&Fetch{}).New()
	storage := (&Storage{}).New()
	debounce := (&Debounce{}).New()
	bitcoin := (&Bitcoin{}).New()

	// Create fresh app with initialized stores
	app := &App{
		CurrentComponent: p.New(basics),
		routes: []p.Route{
			{Path: "/", HTMLFile: "index.html", SSRPath: "/", Component: basics},
			{Path: "/basics", HTMLFile: "basics.html", SSRPath: "/basics", Component: basics},
			{Path: "/derived", HTMLFile: "derived.html", SSRPath: "/derived", Component: derived},
			{Path: "/complex", HTMLFile: "complex.html", SSRPath: "/complex", Component: complex},
			{Path: "/components", HTMLFile: "components.html", SSRPath: "/components", Component: components},
			{Path: "/lists", HTMLFile: "lists.html", SSRPath: "/lists", Component: lists},
			{Path: "/routing", HTMLFile: "routing.html", SSRPath: "/routing", Component: routing},
			{Path: "/links", HTMLFile: "links.html", SSRPath: "/links", Component: links},
			{Path: "/fetch", HTMLFile: "fetch.html", SSRPath: "/fetch", Component: fetch},
			{Path: "/storage", HTMLFile: "storage.html", SSRPath: "/storage", Component: storage},
			{Path: "/debounce", HTMLFile: "debounce.html", SSRPath: "/debounce", Component: debounce},
			{Path: "/bitcoin", HTMLFile: "bitcoin.html", SSRPath: "/bitcoin", Component: bitcoin},
		},
	}
	return app
}

func (a *App) OnMount() {
	router := p.NewRouter(a.CurrentComponent, a.routes, "a unique id")
	router.NotFound(func() {
		a.CurrentComponent.Set(nil)
	})
	router.Start()
}

func (a *App) Routes() []p.Route {
	return a.routes
}

func (a *App) Render() p.Node {
	return p.Html(`<div class="showcase">
		<nav class="sidebar">
			<h2>Reactive</h2>
			<ul>
				<li><a href="/basics">Basics</a></li>
				<li><a href="/derived">Derived</a></li>
				<li><a href="/complex">Complex</a></li>
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
		<main id="content" class="content">`,
		a.CurrentComponent,
		`</main>
	</div>`)
}

func (a *App) GlobalStyle() string {
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
.demo pre.code{background:#1a1a2e;color:#e0e0e0}
.demo textarea{width:100%;padding:10px;border:1px solid #ccc;border-radius:4px;font-family:inherit;resize:vertical}
.hint{font-size:.9em;color:#666;font-style:italic}
.buttons{display:flex;gap:10px;flex-wrap:wrap;margin:10px 0}
`
}
