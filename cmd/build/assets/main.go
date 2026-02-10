package main

import p "github.com/tbocek/preveltekit/v2"

type App struct {
	CurrentComponent *p.Store[p.Component]
}

func (a *App) New() p.Component {
	return &App{
		CurrentComponent: p.New[p.Component](&Hello{}),
	}
}

func (a *App) Routes() []p.Route {
	return []p.Route{
		{Path: "/", HTMLFile: "index.html", SSRPath: "/", Component: &Hello{}},
	}
}

func (a *App) OnMount() {
	router := p.NewRouter(a.CurrentComponent, a.Routes(), "app")
	router.Start()
}

func (a *App) Render() p.Node {
	return p.Html(`<main>`, a.CurrentComponent, `</main>`)
}

type Hello struct{}

func (h *Hello) Render() p.Node {
	return p.Html(`<h1>Hello, World!</h1>`)
}

func main() {
	p.Hydrate(&App{})
}
