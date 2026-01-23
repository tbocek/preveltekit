package main

// Index is the home page
type Index struct{}

func (i *Index) Template() string {
	return `<div class="page index-page">
	<h1>Welcome to Reactive</h1>
	<p class="subtitle">Build fast, reactive web apps with Go and WebAssembly</p>

	<div class="features">
		<div class="feature">
			<h3>Reactive Stores</h3>
			<p>Generic reactive state management with automatic DOM updates</p>
		</div>
		<div class="feature">
			<h3>Component System</h3>
			<p>Compose UIs with reusable components, props, and slots</p>
		</div>
		<div class="feature">
			<h3>Client-Side Routing</h3>
			<p>SPA navigation with history API and route parameters</p>
		</div>
		<div class="feature">
			<h3>SSR + Hydration</h3>
			<p>Pre-render on server, hydrate on client for fast initial load</p>
		</div>
	</div>

	<div class="cta">
		<a href="/doc" class="btn primary">Get Started</a>
		<a href="/example" class="btn secondary">View Examples</a>
	</div>

	<p class="routing-note">(These buttons use client-side routing - no page reloads)</p>
	<p class="routing-note routing-compare">Compare with the nav links above which use server-side routing (page reloads)</p>
</div>`
}

func (i *Index) Style() string {
	return `
.index-page { text-align: center; padding: 3rem 1rem; }
.index-page h1 { font-size: 3rem; margin-bottom: 0.5rem; color: #1a1a2e; }
.subtitle { font-size: 1.25rem; color: #666; margin-bottom: 3rem; }
.features { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 2rem; margin: 3rem 0; text-align: left; }
.feature { background: #f8f9fa; padding: 1.5rem; border-radius: 8px; border: 1px solid #e9ecef; }
.feature h3 { color: #1a1a2e; margin-bottom: 0.5rem; }
.feature p { color: #666; font-size: 0.95rem; }
.cta { display: flex; gap: 1rem; justify-content: center; margin-top: 2rem; }
.btn { padding: 0.75rem 2rem; border-radius: 6px; text-decoration: none; font-weight: 500; transition: all 0.2s; }
.btn.primary { background: #1a1a2e; color: #fff; }
.btn.primary:hover { background: #2d2d44; }
.btn.secondary { background: #fff; color: #1a1a2e; border: 2px solid #1a1a2e; }
.btn.secondary:hover { background: #f5f5f5; }
.routing-note { margin-top: 1.5rem; color: #888; font-size: 0.9rem; font-style: italic; }
.routing-compare { margin-top: 0.5rem; }
`
}
