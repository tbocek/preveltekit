package main

import p "github.com/tbocek/preveltekit/v2"

type Home struct{}

func (h *Home) New() p.Component {
	return &Home{}
}

func (h *Home) Render() p.Node {
	return p.Fragment(
		p.Section(p.Attr("class", "hero"),
			p.Div(p.Attr("class", "container"),
				p.H1("PrevelteKit"),
				p.P(p.Attr("class", "tagline"), "A lightweight, high-performance web application framework written in Go, featuring server-side rendering with WebAssembly hydration."),
				p.Div(p.Attr("class", "hero-badges"),
					p.Span(p.Attr("class", "badge"), "Go"),
					p.Span(p.Attr("class", "badge"), "SSR"),
					p.Span(p.Attr("class", "badge"), "WebAssembly"),
				),
				p.P(p.Attr("class", "hero-size"), p.RawHTML("This entire site (3 pages, routing, live Bitcoin API) is <strong>~44 kB</strong> brotli-compressed.")),
			),
		),

		p.Section(p.Attr("class", "highlights"),
			p.Div(p.Attr("class", "container"),
				p.Div(p.Attr("class", "highlight-grid"),
					p.Div(p.Attr("class", "highlight-card"),
						p.H3("Static-first architecture"),
						p.P("Deploy anywhere as pure HTML + WASM. No server runtime needed — works on GitHub Pages, S3, or any CDN."),
					),
					p.Div(p.Attr("class", "highlight-card"),
						p.H3("Lightning fast builds"),
						p.P("Native Go compiler for SSR, TinyGo for WASM. No bundler, no transpiler, no build tool chain."),
					),
					p.Div(p.Attr("class", "highlight-card"),
						p.H3("No JavaScript ecosystem"),
						p.P("Pure Go from component to deployment. No npm, no node_modules, no package.json."),
					),
				),
			),
		),

		p.Section(p.Attr("class", "explore"),
			p.Div(p.Attr("class", "container"),
				p.H2("Explore"),
				p.Div(p.Attr("class", "explore-grid"),
					p.Div(p.Attr("class", "explore-card"),
						p.H3("Manual"),
						p.P("API reference covering stores, components, routing, fetch, and more."),
						p.Div(p.Attr("class", "explore-links"),
							p.A(p.Attr("href", "manual"), "Client-side"),
							p.A(p.Attr("href", "manual"), p.Attr("external", ""), "Server-side"),
						),
					),
					p.Div(p.Attr("class", "explore-card"),
						p.H3("Bitcoin Price Demo"),
						p.P("Live Bitcoin price fetched from an API with auto-refresh — a real-world example."),
						p.Div(p.Attr("class", "explore-links"),
							p.A(p.Attr("href", "bitcoin"), "Client-side"),
							p.A(p.Attr("href", "bitcoin"), p.Attr("external", ""), "Server-side"),
						),
					),
					p.Div(p.Attr("class", "explore-card"),
						p.H3("Examples"),
						p.P("Interactive showcase with 11 examples: basics, lists, routing, fetch, storage, and more."),
						p.Div(p.Attr("class", "explore-links"),
							p.A(p.Attr("href", "https://github.com/tbocek/preveltekit/tree/main/examples"), p.Attr("external", ""), "View on GitHub"),
						),
					),
				),
			),
		),

		p.Section(p.Attr("class", "features"),
			p.Div(p.Attr("class", "container"),
				p.H2("Features"),
				p.Div(p.Attr("class", "feature-grid"),
					p.Div(p.Attr("class", "feature-card"),
						p.Div(p.Attr("class", "feature-icon"), p.RawHTML("&#9889;")),
						p.H3("Lightning Fast Builds"),
						p.P("Go compiles in milliseconds. TinyGo produces compact WASM binaries. No slow bundler step."),
					),
					p.Div(p.Attr("class", "feature-card"),
						p.Div(p.Attr("class", "feature-icon"), p.RawHTML("&#9723;")),
						p.H3("Minimalistic"),
						p.P("Small framework with no code generation, no intermediate files, no complex toolchain."),
					),
					p.Div(p.Attr("class", "feature-card"),
						p.Div(p.Attr("class", "feature-icon"), p.RawHTML("&#9881;")),
						p.H3("SSR + WASM Hydration"),
						p.P("Server-side rendered HTML for instant page loads. WebAssembly hydrates for interactivity."),
					),
					p.Div(p.Attr("class", "feature-card"),
						p.Div(p.Attr("class", "feature-icon"), p.RawHTML("&#9729;")),
						p.H3("Deploy Anywhere"),
						p.P("Output is pure static files. GitHub Pages, S3, Cloudflare, Netlify — any static host works."),
					),
					p.Div(p.Attr("class", "feature-card"),
						p.Div(p.Attr("class", "feature-icon"), p.RawHTML("&#9989;")),
						p.H3("Type-Safe"),
						p.P("Full Go type safety with compile-time error checking. Catch bugs before they reach production."),
					),
					p.Div(p.Attr("class", "feature-card"),
						p.Div(p.Attr("class", "feature-icon"), p.RawHTML("&#9655;")),
						p.H3("Zero Config"),
						p.P(p.RawHTML("Works with sensible defaults. Just <code>go run</code> for SSR and <code>tinygo build</code> for WASM.")),
					),
				),
			),
		),

		p.Section(p.Attr("class", "quickstart"),
			p.Div(p.Attr("class", "container"),
				p.H2("Quick Start"),
				p.P(p.Attr("class", "prereq"), p.RawHTML(`Requires <a href="https://go.dev/dl/" external>Go</a> and <a href="https://docs.docker.com/get-docker/" external>Docker</a>.`)),
				p.Div(p.Attr("class", "steps"),
					p.Div(p.Attr("class", "step"),
						p.Div(p.Attr("class", "step-num"), "1"),
						p.Div(p.Attr("class", "step-content"),
							p.H3("Create a new project"),
							p.Pre(p.Code("mkdir myapp && cd myapp\ngo mod init myapp\ngo run github.com/tbocek/preveltekit/v2/cmd/build@latest init")),
							p.P(p.Attr("class", "step-hint"), p.RawHTML(`Scaffolds <code>main.go</code>, <code>build.sh</code>, <code>Dockerfile</code>, and assets. Or copy from <a href="https://github.com/tbocek/preveltekit/tree/main/cmd/build/assets">cmd/build/assets/</a>.`)),
						),
					),
					p.Div(p.Attr("class", "step"),
						p.Div(p.Attr("class", "step-num"), "2"),
						p.Div(p.Attr("class", "step-content"),
							p.H3("Build and run with Docker"),
							p.Pre(p.Code("docker build -f Dockerfile.dev -t myapp-dev .\ndocker run --init -p 8080:8080 -v $PWD:/app myapp-dev")),
							p.P(p.Attr("class", "step-hint"), p.RawHTML(`Open <code>http://localhost:8080</code> — live reload on file changes, no other dependencies required.`)),
						),
					),
				),
			),
		),
	)
}

func (h *Home) Style() string {
	return `
.hero{background:linear-gradient(135deg,#1a1a2e 0%,#16213e 50%,#0f3460 100%);color:#fff;padding:80px 0;text-align:center}
.hero h1{font-size:3em;margin-bottom:16px;font-weight:800}
.tagline{font-size:1.2em;color:#b8c5d6;max-width:600px;margin:0 auto 24px}
.hero-badges{display:flex;gap:12px;justify-content:center}
.badge{background:rgba(255,255,255,.12);padding:6px 16px;border-radius:20px;font-size:13px;font-weight:500;border:1px solid rgba(255,255,255,.2)}
.hero-size{margin-top:24px;font-size:.95em;color:#8899b0}

.highlights{padding:60px 0;background:#f8f9fa}
.highlight-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:24px}
.highlight-card{padding:24px;background:#fff;border-radius:8px;box-shadow:0 1px 3px rgba(0,0,0,.08)}
.highlight-card h3{font-size:1.1em;margin-bottom:8px;color:#1a1a2e}
.highlight-card p{font-size:.95em;color:#666}

.explore{padding:60px 0}
.explore h2{text-align:center;font-size:2em;margin-bottom:40px;color:#1a1a2e}
.explore-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:24px}
.explore-card{padding:24px;border:1px solid #e5e7eb;border-radius:8px}
.explore-card h3{font-size:1.1em;margin-bottom:8px;color:#1a1a2e}
.explore-card p{font-size:.95em;color:#666;margin-bottom:16px}
.explore-links{display:flex;gap:12px}
.explore-links a{display:inline-block;padding:8px 16px;border-radius:4px;font-size:13px;font-weight:500;background:#1a1a2e;color:#fff;transition:background .2s}
.explore-links a:hover{background:#0f3460}
.explore-links a[external]{background:#f8f9fa;color:#1a1a2e;border:1px solid #ddd}
.explore-links a[external]:hover{background:#e5e7eb}

.features{padding:60px 0;background:#f8f9fa}
.features h2{text-align:center;font-size:2em;margin-bottom:40px;color:#1a1a2e}
.feature-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:24px}
.feature-card{padding:24px;border:1px solid #e5e7eb;border-radius:8px;text-align:center;background:#fff}
.feature-icon{font-size:2em;margin-bottom:12px}
.feature-card h3{font-size:1.1em;margin-bottom:8px;color:#1a1a2e}
.feature-card p{font-size:.95em;color:#666}

.quickstart{padding:60px 0}
.quickstart h2{text-align:center;font-size:2em;margin-bottom:12px;color:#1a1a2e}
.prereq{text-align:center;color:#666;margin-bottom:32px;font-size:.95em}
.prereq a{color:#4a7fff;text-decoration:underline}
.steps{max-width:700px;margin:0 auto}
.step{display:flex;gap:20px;margin-bottom:32px}
.step-num{width:36px;height:36px;background:#1a1a2e;color:#fff;border-radius:50%;display:flex;align-items:center;justify-content:center;font-weight:700;flex-shrink:0}
.step-content{flex:1}
.step-content h3{font-size:1.1em;margin-bottom:12px;color:#1a1a2e}
.step-hint{font-size:.9em;color:#666;margin-top:8px}
.step-hint a{color:#4a7fff;text-decoration:underline}
.step-hint code{background:#e5e7eb;padding:2px 6px;border-radius:3px;font-size:.85em}

@media(max-width:768px){
.highlight-grid,.feature-grid,.explore-grid{grid-template-columns:1fr}
.hero h1{font-size:2em}
.tagline{font-size:1em}
}
`
}
