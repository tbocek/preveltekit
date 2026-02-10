package main

import p "github.com/tbocek/preveltekit"

// Links showcase - demonstrates link handling
type Links struct {
	LastNavigation *p.Store[string]
}

func (l *Links) New() p.Component {
	return &Links{
		LastNavigation: p.New(""),
	}
}

func (l *Links) Render() p.Node {
	return p.Html(`<div class="demo">
		<h1>Links</h1>

		<section>
			<h2>Client-Side vs Server-Side</h2>
			<p>Same URL, different behavior:</p>

			<div class="link-list">
				<a href="/lists" class="nav-link">
					<span class="link-icon">-></span>
					<span>/lists</span>
					<span class="link-type">Client-side</span>
				</a>
				<a href="/lists" external class="nav-link external">
					<span class="link-icon">^</span>
					<span>/lists</span>
					<span class="link-type">Server (reload)</span>
				</a>
			</div>
			<p class="hint">Click both - first one is instant, second reloads the page.</p>
		</section>

		<section>
			<h2>When to use <code>external</code></h2>
			<ul class="info-list">
				<li>Server-side routes (API, auth, downloads)</li>
				<li>Full page refresh needed</li>
				<li>Links to other apps on same domain</li>
			</ul>
		</section>

		<section>
			<h2>Try It</h2>
			<div class="button-links">
				<a href="/basics" class="btn-link primary">Basics (SPA)</a>
				<a href="/basics" external class="btn-link secondary">Basics (Reload)</a>
			</div>
		</section>

		<section>
			<h2>Code</h2>
			<pre class="code">// client-side navigation (SPA) — default
&lt;a href="/basics">Basics&lt;/a>

// server-side navigation (full reload)
&lt;a href="/basics" external>Basics&lt;/a>

// router setup in App.OnMount():
router := p.NewRouter(CurrentComponent, routes, "id")
router.NotFound(func() { CurrentComponent.Set(nil) })
router.Start()

// route definition:
p.Route{
    Path:      "/basics",
    HTMLFile:  "basics.html",
    SSRPath:   "/basics",
    Component: basics,
}</pre>
		</section>
	</div>`)
}

func (l *Links) Style() string {
	return `
.demo{max-width:700px}
.demo pre.code{background:#1a1a2e;color:#e0e0e0;font-size:12px;margin-top:12px}
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
