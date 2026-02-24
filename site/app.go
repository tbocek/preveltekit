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
	return p.Fragment(
		p.Header(p.Attr("class", "header"),
			p.Div(p.Attr("class", "header-inner"),
				p.A(p.Attr("href", "./"), p.Attr("class", "brand"), "PrevelteKit"),
				p.Nav(p.Attr("class", "nav"),
					p.A(p.Attr("href", "manual"), "Manual"),
					p.A(p.Attr("href", "bitcoin"), "Bitcoin Demo"),
					p.A(p.Attr("href", "https://github.com/tbocek/preveltekit"), p.Attr("external", ""), "GitHub"),
				),
			),
		),
		p.Div(p.Attr("id", "content"),
			a.CurrentComponent,
		),
		p.Footer(p.Attr("class", "footer"),
			p.Div(p.Attr("class", "container"),
				p.P(p.RawHTML("SSR + WASM hydration &bull; Pure Go &bull; Static deployment &bull; No JavaScript required")),
			),
		),
	)
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

.page{padding:40px 0}
.page h1{font-size:2.2em;color:#1a1a2e;margin-bottom:8px}
.page-intro{color:#666;margin-bottom:32px;font-size:1.05em}

pre{background:#1a1a2e;color:#e0e0e0;padding:16px;border-radius:6px;overflow-x:auto;font-size:13px;line-height:1.6}
pre code{background:transparent;padding:0;font-size:inherit}
code{background:#f1f5f9;padding:2px 6px;border-radius:3px;font-size:.85em}

.footer{padding:32px 0;background:#1a1a2e;color:#999;text-align:center;font-size:.9em}

@media(max-width:768px){
.nav{gap:12px}
.nav a{font-size:13px}
}
`
}
