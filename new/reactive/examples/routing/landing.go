package main

// Landing is the landing/marketing page
type Landing struct{}

func (l *Landing) Template() string {
	return `<div class="page landing-page">
	<section class="hero">
		<h1>Build Modern Web Apps in Go</h1>
		<p>Reactive brings the power of Go to the browser with WebAssembly. Write type-safe, performant web applications without JavaScript.</p>
	</section>

	<section class="benefits">
		<h2>Why Reactive?</h2>
		<div class="benefit-grid">
			<div class="benefit">
				<div class="benefit-icon">🚀</div>
				<h3>Blazing Fast</h3>
				<p>WebAssembly runs at near-native speed. TinyGo produces tiny binaries that load instantly.</p>
			</div>
			<div class="benefit">
				<div class="benefit-icon">🔒</div>
				<h3>Type Safe</h3>
				<p>Catch errors at compile time. Go's type system prevents entire classes of runtime bugs.</p>
			</div>
			<div class="benefit">
				<div class="benefit-icon">📦</div>
				<h3>Small Bundles</h3>
				<p>Typical apps are under 50KB gzipped. No massive JavaScript framework overhead.</p>
			</div>
			<div class="benefit">
				<div class="benefit-icon">🔄</div>
				<h3>Reactive</h3>
				<p>Automatic DOM updates when state changes. No manual DOM manipulation needed.</p>
			</div>
		</div>
	</section>

	<section class="comparison">
		<h2>Compared to JavaScript Frameworks</h2>
		<table class="comparison-table">
			<thead>
				<tr>
					<th>Feature</th>
					<th>Reactive</th>
					<th>React</th>
					<th>Vue</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td>Bundle Size</td>
					<td class="highlight">~15KB</td>
					<td>~45KB</td>
					<td>~35KB</td>
				</tr>
				<tr>
					<td>Type Safety</td>
					<td class="highlight">Native</td>
					<td>Optional</td>
					<td>Optional</td>
				</tr>
				<tr>
					<td>Runtime</td>
					<td class="highlight">WASM</td>
					<td>JS</td>
					<td>JS</td>
				</tr>
				<tr>
					<td>SSR</td>
					<td class="highlight">Built-in</td>
					<td>Plugin</td>
					<td>Plugin</td>
				</tr>
			</tbody>
		</table>
	</section>
</div>`
}

func (l *Landing) Style() string {
	return `
.landing-page { padding: 2rem 0; }
.hero { text-align: center; padding: 4rem 2rem; background: linear-gradient(135deg, #1a1a2e 0%, #2d2d44 100%); color: #fff; border-radius: 12px; margin-bottom: 3rem; }
.hero h1 { font-size: 2.5rem; margin-bottom: 1rem; }
.hero p { font-size: 1.1rem; max-width: 600px; margin: 0 auto; opacity: 0.9; }
.benefits { margin-bottom: 3rem; }
.benefits h2 { text-align: center; margin-bottom: 2rem; color: #1a1a2e; }
.benefit-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 1.5rem; }
.benefit { background: #fff; padding: 1.5rem; border-radius: 8px; border: 1px solid #e9ecef; text-align: center; }
.benefit-icon { font-size: 2.5rem; margin-bottom: 1rem; }
.benefit h3 { color: #1a1a2e; margin-bottom: 0.5rem; }
.benefit p { color: #666; font-size: 0.9rem; }
.comparison { margin-bottom: 3rem; }
.comparison h2 { text-align: center; margin-bottom: 2rem; color: #1a1a2e; }
.comparison-table { width: 100%; border-collapse: collapse; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.comparison-table th, .comparison-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #e9ecef; }
.comparison-table th { background: #1a1a2e; color: #fff; }
.comparison-table .highlight { color: #28a745; font-weight: 600; }
`
}
