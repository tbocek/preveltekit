package main

import p "github.com/tbocek/preveltekit/v2"

// UserProfile is the page for the route "/user/:name": the name comes from the URL.
type UserProfile struct {
	Name *p.Store[string] // router.Param("name"), set in App.OnMount
}

func (u *UserProfile) Render() p.Node {
	return p.Div(p.Attr("class", "demo"),
		p.H1("User Profile"),
		p.Section(
			p.H2("Route Parameter"),
			p.P("This page is the route ", p.Code("/user/:name"), ". The name is read from the URL:"),
			p.P(p.Attr("class", "user-name"), p.Strong(u.Name)),
			p.P(
				p.A(p.Attr("href", "/user/alice"), "Alice"), " · ",
				p.A(p.Attr("href", "/user/bob"), "Bob"), " · ",
				p.A(p.Attr("href", "/routing"), "Back to Routing"),
			),
			p.Pre(p.Attr("class", "code"), `// route with a parameter
{Path: "/user/:name", HTMLFile: "user.html", Component: user}

// in OnMount: the router updates the store on every navigation
router := p.NewRouter(a.CurrentComponent, a.routes, "app")
user.Name = router.Param("name")
router.Start()

// bind it like any store
p.Strong(u.Name)`),
		),
	)
}
