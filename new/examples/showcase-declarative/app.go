package main

import p "preveltekit"

type App struct {
	CurrentPage *p.Store[string]
}

func (a *App) OnMount() {
	router := p.NewRouter()
	router.RegisterRoutes(a.Routes())
	router.NotFound(func() {
		a.CurrentPage.Set("notfound")
	})
	router.Start()
}

func (a *App) Routes() []p.StaticRoute {
	return []p.StaticRoute{
		{Path: "/", HTMLFile: "index.html", Handler: func(params map[string]string) { a.CurrentPage.Set("basics") }},
		{Path: "/basics", HTMLFile: "basics.html", Handler: func(params map[string]string) { a.CurrentPage.Set("basics") }},
		{Path: "/components", HTMLFile: "components.html", Handler: func(params map[string]string) { a.CurrentPage.Set("components") }},
		{Path: "/lists", HTMLFile: "lists.html", Handler: func(params map[string]string) { a.CurrentPage.Set("lists") }},
		{Path: "/routing", HTMLFile: "routing.html", Handler: func(params map[string]string) { a.CurrentPage.Set("routing") }},
		{Path: "/links", HTMLFile: "links.html", Handler: func(params map[string]string) { a.CurrentPage.Set("links") }},
		{Path: "/fetch", HTMLFile: "fetch.html", Handler: func(params map[string]string) { a.CurrentPage.Set("fetch") }},
		{Path: "/storage", HTMLFile: "storage.html", Handler: func(params map[string]string) { a.CurrentPage.Set("storage") }},
		{Path: "/debounce", HTMLFile: "debounce.html", Handler: func(params map[string]string) { a.CurrentPage.Set("debounce") }},
		{Path: "/bitcoin", HTMLFile: "bitcoin.html", Handler: func(params map[string]string) { a.CurrentPage.Set("bitcoin") }},
	}
}

func (a *App) Render() p.Node {
	return p.Div(p.Class("showcase"),
		p.Nav(p.Class("sidebar"),
			p.H2("Reactive"),
			p.Ul(
				p.Li(p.A(p.Href("/basics"), "Basics")),
				p.Li(p.A(p.Href("/components"), "Components")),
				p.Li(p.A(p.Href("/lists"), "Lists")),
				p.Li(p.A(p.Href("/routing"), "Routing")),
				p.Li(p.A(p.Href("/links"), "Links")),
				p.Li(p.A(p.Href("/fetch"), "Fetch")),
				p.Li(p.A(p.Href("/storage"), "Storage")),
				p.Li(p.A(p.Href("/debounce"), "Debounce")),
				p.Li(p.A(p.Href("/bitcoin"), "Bitcoin")),
			),
		),
		p.Main(p.Class("content"),
			p.If(a.CurrentPage.Eq("basics"),
				p.Child("basics"),
			).ElseIf(a.CurrentPage.Eq("components"),
				p.Child("components"),
			).ElseIf(a.CurrentPage.Eq("lists"),
				p.Child("lists"),
			).ElseIf(a.CurrentPage.Eq("routing"),
				p.Child("routing"),
			).ElseIf(a.CurrentPage.Eq("links"),
				p.Child("links"),
			).ElseIf(a.CurrentPage.Eq("fetch"),
				p.Child("fetch"),
			).ElseIf(a.CurrentPage.Eq("storage"),
				p.Child("storage"),
			).ElseIf(a.CurrentPage.Eq("debounce"),
				p.Child("debounce"),
			).ElseIf(a.CurrentPage.Eq("bitcoin"),
				p.Child("bitcoin"),
			).Else(
				p.P("Page not found"),
			),
		),
	)
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

func (a *App) HandleEvent(method string, args string) {
	// App has no event handlers
}
