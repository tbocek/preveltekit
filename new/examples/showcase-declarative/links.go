package main

import p "preveltekit"

// Links showcase - demonstrates link handling
type Links struct {
	LastNavigation *p.Store[string]
}

func (l *Links) OnMount() {
	l.LastNavigation.Set("")
}

func (l *Links) Render() p.Node {
	return p.Div(p.Class("demo"),
		p.H1("Links"),

		p.Section(
			p.H2("Client-Side vs Server-Side"),
			p.P("Same URL, different behavior:"),

			p.Div(p.Class("link-list"),
				p.A(p.Href("/lists"), p.Class("nav-link"),
					p.Span(p.Class("link-icon"), "->"),
					p.Span("/lists"),
					p.Span(p.Class("link-type"), "Client-side"),
				),
				p.A(p.Href("/lists"), p.Attr("external", ""), p.Class("nav-link", "external"),
					p.Span(p.Class("link-icon"), "^"),
					p.Span("/lists"),
					p.Span(p.Class("link-type"), "Server (reload)"),
				),
			),
			p.P(p.Class("hint"), "Click both - first one is instant, second reloads the page."),
		),

		p.Section(
			p.H2("When to use ", p.Code("external")),
			p.Ul(p.Class("info-list"),
				p.Li("Server-side routes (API, auth, downloads)"),
				p.Li("Full page refresh needed"),
				p.Li("Links to other apps on same domain"),
			),
		),

		p.Section(
			p.H2("Try It"),
			p.Div(p.Class("button-links"),
				p.A(p.Href("/basics"), p.Class("btn-link", "primary"), "Basics (SPA)"),
				p.A(p.Href("/basics"), p.Attr("external", ""), p.Class("btn-link", "secondary"), "Basics (Reload)"),
			),
		),
	)
}

func (l *Links) Style() string {
	return `
.demo{max-width:700px}
.demo code{background:#f1f1f1;padding:2px 6px;border-radius:3px;font-size:.9em}
.link-list{display:flex;flex-direction:column;gap:8px}
.nav-link{display:flex;align-items:center;gap:12px;padding:12px 16px;border:1px solid #ddd;border-radius:6px;text-decoration:none;color:#333}
.nav-link:hover{border-color:#007bff;background:#f8f9fa}
.nav-link .link-icon{font-family:monospace;font-weight:700;color:#007bff}
.nav-link .link-type{margin-left:auto;font-size:11px;padding:2px 8px;border-radius:10px;background:#e8f4fd;color:#007bff}
.nav-link.external .link-type{background:#fff3cd;color:#856404}
.nav-link.external .link-icon{color:#856404}
.button-links{display:flex;gap:10px;flex-wrap:wrap}
.btn-link{display:inline-block;padding:10px 20px;border-radius:4px;text-decoration:none;font-weight:500}
.btn-link.primary{background:#007bff;color:#fff}
.btn-link.primary:hover{background:#0056b3}
.btn-link.secondary{background:#6c757d;color:#fff}
.btn-link.secondary:hover{background:#545b62}
.info-list{margin:0;padding-left:20px;color:#555}
.info-list li{margin:5px 0}
`
}

func (l *Links) HandleEvent(method string, args string) {
	// Links has no event handlers
}
