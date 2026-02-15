package main

import p "github.com/tbocek/preveltekit/v2"

type Home struct{}

func (h *Home) New() p.Component {
	return &Home{}
}

func (h *Home) Render() p.Node {
	return p.Html(`
	<section class="hero">
		<div class="container">
			<h1>PrevelteKit</h1>
			<p class="tagline">A lightweight, high-performance web application framework written in Go, featuring server-side rendering with WebAssembly hydration.</p>
			<div class="hero-badges">
				<span class="badge">Go</span>
				<span class="badge">SSR</span>
				<span class="badge">WebAssembly</span>
			</div>
			<p class="hero-size">This entire site (3 pages, routing, live Bitcoin API) is <strong>~44 kB</strong> brotli-compressed.</p>
		</div>
	</section>

	<section class="highlights">
		<div class="container">
			<div class="highlight-grid">
				<div class="highlight-card">
					<h3>Static-first architecture</h3>
					<p>Deploy anywhere as pure HTML + WASM. No server runtime needed — works on GitHub Pages, S3, or any CDN.</p>
				</div>
				<div class="highlight-card">
					<h3>Lightning fast builds</h3>
					<p>Native Go compiler for SSR, TinyGo for WASM. No bundler, no transpiler, no build tool chain.</p>
				</div>
				<div class="highlight-card">
					<h3>No JavaScript ecosystem</h3>
					<p>Pure Go from component to deployment. No npm, no node_modules, no package.json.</p>
				</div>
			</div>
		</div>
	</section>

	<section class="explore">
		<div class="container">
			<h2>Explore</h2>
			<div class="explore-grid">
				<div class="explore-card">
					<h3>Manual</h3>
					<p>API reference covering stores, components, routing, fetch, and more.</p>
					<div class="explore-links">
						<a href="manual">Client-side</a>
						<a href="manual" external>Server-side</a>
					</div>
				</div>
				<div class="explore-card">
					<h3>Bitcoin Price Demo</h3>
					<p>Live Bitcoin price fetched from an API with auto-refresh — a real-world example.</p>
					<div class="explore-links">
						<a href="bitcoin">Client-side</a>
						<a href="bitcoin" external>Server-side</a>
					</div>
				</div>
				<div class="explore-card">
					<h3>Examples</h3>
					<p>Interactive showcase with 11 examples: basics, lists, routing, fetch, storage, and more.</p>
					<div class="explore-links">
						<a href="https://github.com/tbocek/preveltekit/tree/main/examples" external>View on GitHub</a>
					</div>
				</div>
			</div>
		</div>
	</section>

	<section class="features">
		<div class="container">
			<h2>Features</h2>
			<div class="feature-grid">
				<div class="feature-card">
					<div class="feature-icon">&#9889;</div>
					<h3>Lightning Fast Builds</h3>
					<p>Go compiles in milliseconds. TinyGo produces compact WASM binaries. No slow bundler step.</p>
				</div>
				<div class="feature-card">
					<div class="feature-icon">&#9723;</div>
					<h3>Minimalistic</h3>
					<p>Small framework with no code generation, no intermediate files, no complex toolchain.</p>
				</div>
				<div class="feature-card">
					<div class="feature-icon">&#9881;</div>
					<h3>SSR + WASM Hydration</h3>
					<p>Server-side rendered HTML for instant page loads. WebAssembly hydrates for interactivity.</p>
				</div>
				<div class="feature-card">
					<div class="feature-icon">&#9729;</div>
					<h3>Deploy Anywhere</h3>
					<p>Output is pure static files. GitHub Pages, S3, Cloudflare, Netlify — any static host works.</p>
				</div>
				<div class="feature-card">
					<div class="feature-icon">&#9989;</div>
					<h3>Type-Safe</h3>
					<p>Full Go type safety with compile-time error checking. Catch bugs before they reach production.</p>
				</div>
				<div class="feature-card">
					<div class="feature-icon">&#9655;</div>
					<h3>Zero Config</h3>
					<p>Works with sensible defaults. Just <code>go run</code> for SSR and <code>tinygo build</code> for WASM.</p>
				</div>
			</div>
		</div>
	</section>

	<section class="quickstart">
		<div class="container">
			<h2>Quick Start</h2>
			<p class="prereq">Requires <a href="https://go.dev/dl/" external>Go</a> and <a href="https://docs.docker.com/get-docker/" external>Docker</a>.</p>
			<div class="steps">
				<div class="step">
					<div class="step-num">1</div>
					<div class="step-content">
						<h3>Create a new project</h3>
						<pre><code>mkdir myapp && cd myapp
go mod init myapp
go run github.com/tbocek/preveltekit/v2/cmd/build@latest init</code></pre>
						<p class="step-hint">Scaffolds <code>main.go</code>, <code>build.sh</code>, <code>Dockerfile</code>, and assets. Or copy from <a href="https://github.com/tbocek/preveltekit/tree/main/cmd/build/assets">cmd/build/assets/</a>.</p>
					</div>
				</div>
				<div class="step">
					<div class="step-num">2</div>
					<div class="step-content">
						<h3>Build and run with Docker</h3>
						<pre><code>docker build -f Dockerfile.dev -t myapp-dev .
docker run --init -p 8080:8080 -v $PWD:/app myapp-dev</code></pre>
						<p class="step-hint">Open <code>http://localhost:8080</code> — live reload on file changes, no other dependencies required.</p>
					</div>
				</div>
			</div>
		</div>
	</section>
	`)
}

func (h *Home) Style() string {
	return `
.hero{background:linear-gradient(135deg,#1a1a2e 0%,#16213e 50%,#0f3460 100%);color:#fff;padding:80px 0;text-align:center}
.hero h1{font-size:3em;margin-bottom:16px;font-weight:800}
.tagline{font-size:1.2em;color:#b8c5d6;max-width:600px;margin:0 auto 24px}
.hero-badges{display:flex;gap:12px;justify-content:center}
.badge{background:rgba(255,255,255,.12);padding:6px 16px;border-radius:20px;font-size:13px;font-weight:500;border:1px solid rgba(255,255,255,.2)}
.hero-size{margin-top:20px;font-size:.95em;color:#8899b0}

.highlights{padding:60px 0;background:#f8f9fa}
.highlight-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:24px}
.highlight-card{padding:24px;background:#fff;border-radius:8px;box-shadow:0 1px 3px rgba(0,0,0,.08)}
.highlight-card h3{font-size:1.1em;margin-bottom:8px;color:#1a1a2e}
.highlight-card p{font-size:.95em;color:#555}

.explore{padding:60px 0}
.explore h2{text-align:center;font-size:2em;margin-bottom:40px;color:#1a1a2e}
.explore-grid{display:grid;grid-template-columns:repeat(3,1fr);gap:24px}
.explore-card{padding:24px;border:1px solid #e5e7eb;border-radius:8px}
.explore-card h3{font-size:1.1em;margin-bottom:8px;color:#1a1a2e}
.explore-card p{font-size:.9em;color:#666;margin-bottom:16px}
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
.feature-card h3{font-size:1em;margin-bottom:8px;color:#1a1a2e}
.feature-card p{font-size:.9em;color:#666}
.feature-card code{background:#f1f5f9;padding:2px 6px;border-radius:3px;font-size:.85em}

.quickstart{padding:60px 0}
.quickstart h2{text-align:center;font-size:2em;margin-bottom:12px;color:#1a1a2e}
.prereq{text-align:center;color:#666;margin-bottom:32px;font-size:.95em}
.prereq a{color:#4a7fff;text-decoration:underline}
.steps{max-width:700px;margin:0 auto}
.step{display:flex;gap:20px;margin-bottom:32px}
.step-num{width:36px;height:36px;background:#1a1a2e;color:#fff;border-radius:50%;display:flex;align-items:center;justify-content:center;font-weight:700;flex-shrink:0}
.step-content{flex:1}
.step-content h3{font-size:1.1em;margin-bottom:12px;color:#1a1a2e}
.step-content pre{background:#1a1a2e;color:#e0e0e0;padding:16px;border-radius:6px;overflow-x:auto;font-size:13px;line-height:1.5}
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
