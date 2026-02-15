package main

import p "github.com/tbocek/preveltekit/v2"

type App struct {
	CurrentComponent *p.Store[p.Component]
	routes           []p.Route
}

func (a *App) New() p.Component {
	home := (&Home{}).New()
	manual := (&Manual{}).New()
	bitcoin := (&BitcoinDemo{}).New()

	app := &App{
		CurrentComponent: p.New(home),
		routes: []p.Route{
			{Path: ".", HTMLFile: "index.html", SSRPath: "/", Component: home},
			{Path: "manual", HTMLFile: "manual.html", SSRPath: "/manual", Component: manual},
			{Path: "bitcoin", HTMLFile: "bitcoin.html", SSRPath: "/bitcoin", Component: bitcoin},
		},
	}
	return app
}

func (a *App) OnMount() {
	router := p.NewRouter(a.CurrentComponent, a.routes, "site-router")
	router.NotFound(func() {
		a.CurrentComponent.Set(nil)
	})
	router.Start()
}

func (a *App) Routes() []p.Route {
	return a.routes
}

func (a *App) Render() p.Node {
	return p.Html(`
	<header class="header">
		<div class="header-inner">
			<a href="./" class="brand">PrevelteKit</a>
			<nav class="nav">
				<a href="manual">Manual</a>
				<a href="bitcoin">Bitcoin Demo</a>
				<a href="https://github.com/tbocek/preveltekit" external>GitHub</a>
			</nav>
		</div>
	</header>
	<div id="content">`,
		a.CurrentComponent,
		`</div>
	<footer class="footer">
		<div class="container">
			<p>SSR + WASM hydration &bull; Pure Go &bull; Static deployment &bull; No JavaScript required</p>
		</div>
	</footer>
	`)
}

func (a *App) GlobalStyle() string {
	return `
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:system-ui,-apple-system,sans-serif;color:#333;line-height:1.6}
a{color:inherit;text-decoration:none}
code{font-family:ui-monospace,SFMono-Regular,Menlo,monospace}

.container{max-width:960px;margin:0 auto;padding:0 20px}

.header{background:#1a1a2e;color:#fff;padding:16px 0;position:sticky;top:0;z-index:10}
.header-inner{max-width:960px;margin:0 auto;padding:0 20px;display:flex;align-items:center;justify-content:space-between}
.brand{font-size:1.3em;font-weight:700;color:#fff}
.nav{display:flex;gap:24px}
.nav a{color:#ccc;font-size:14px;transition:color .2s}
.nav a:hover{color:#fff}

.footer{padding:32px 0;background:#1a1a2e;color:#999;text-align:center;font-size:.9em}

@media(max-width:768px){
.nav{gap:12px}
.nav a{font-size:13px}
}
`
}
