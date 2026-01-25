package main

import "reactive"

// Links showcase - demonstrates link handling
type Links struct {
	ClickedLink *reactive.Store[string]
	LinkCount   *reactive.Store[int]
}

func (l *Links) OnMount() {
	l.ClickedLink.Set("")
	l.LinkCount.Set(0)
}

func (l *Links) TrackClick(title string) {
	l.ClickedLink.Set(title)
	l.LinkCount.Set(l.LinkCount.Get() + 1)
}

func (l *Links) ClickGo() {
	l.TrackClick("Go Documentation")
}

func (l *Links) ClickMDN() {
	l.TrackClick("MDN Web Docs")
}

func (l *Links) ClickGitHub() {
	l.TrackClick("GitHub")
}

func (l *Links) ClickWasm() {
	l.TrackClick("WebAssembly")
}

func (l *Links) Template() string {
	return `<div class="demo">
	<h1>Links</h1>

	<section>
		<h2>External Links</h2>
		<p>Links to external websites (open in new tab):</p>

		<div class="link-grid">
			<a href="https://golang.org/doc/" target="_blank" rel="noopener" class="link-card" @click="ClickGo()">
				<strong>Go Documentation</strong>
				<span>Official Go programming language docs</span>
				<span class="link-url">golang.org</span>
			</a>
			<a href="https://developer.mozilla.org/" target="_blank" rel="noopener" class="link-card" @click="ClickMDN()">
				<strong>MDN Web Docs</strong>
				<span>Web technology documentation</span>
				<span class="link-url">developer.mozilla.org</span>
			</a>
			<a href="https://github.com/" target="_blank" rel="noopener" class="link-card" @click="ClickGitHub()">
				<strong>GitHub</strong>
				<span>Code hosting platform</span>
				<span class="link-url">github.com</span>
			</a>
			<a href="https://webassembly.org/" target="_blank" rel="noopener" class="link-card" @click="ClickWasm()">
				<strong>WebAssembly</strong>
				<span>Binary instruction format for the web</span>
				<span class="link-url">webassembly.org</span>
			</a>
		</div>
	</section>

	<section>
		<h2>Link Tracking</h2>
		<p>Track user interactions with links:</p>

		<div class="tracking-info">
			<div class="stat">
				<span class="stat-value">{LinkCount}</span>
				<span class="stat-label">Links Clicked</span>
			</div>
			<div class="stat">
				<span class="stat-value">{#if ClickedLink != ""}{ClickedLink}{:else}None{/if}</span>
				<span class="stat-label">Last Clicked</span>
			</div>
		</div>
	</section>

	<section>
		<h2>Button Links</h2>
		<p>Styled buttons that act as navigation:</p>

		<div class="button-links">
			<a href="#basics" class="btn-link primary">Go to Basics</a>
			<a href="#lists" class="btn-link secondary">Go to Lists</a>
			<a href="#fetch" class="btn-link outline">Go to Fetch</a>
		</div>
	</section>

	<section>
		<h2>Icon Links</h2>
		<p>Links with icons and different states:</p>

		<div class="icon-links">
			<a href="https://github.com/" target="_blank" rel="noopener" class="icon-link">
				<span class="icon">*</span>
				<span>Star on GitHub</span>
			</a>
			<a href="https://twitter.com/" target="_blank" rel="noopener" class="icon-link">
				<span class="icon">@</span>
				<span>Follow on Twitter</span>
			</a>
			<a href="mailto:hello@example.com" class="icon-link">
				<span class="icon">#</span>
				<span>Contact Us</span>
			</a>
		</div>
	</section>

	<section>
		<h2>Breadcrumb Navigation</h2>
		<p>Show current location in site hierarchy:</p>

		<nav class="breadcrumb">
			<a href="#home">Home</a>
			<span class="separator">/</span>
			<a href="#docs">Documentation</a>
			<span class="separator">/</span>
			<a href="#components">Components</a>
			<span class="separator">/</span>
			<span class="current">Links</span>
		</nav>
	</section>

	<section>
		<h2>Pagination Links</h2>
		<p>Navigate through pages of content:</p>

		<nav class="pagination">
			<a href="#" class="page-link disabled">Previous</a>
			<a href="#page1" class="page-link active">1</a>
			<a href="#page2" class="page-link">2</a>
			<a href="#page3" class="page-link">3</a>
			<span class="page-link ellipsis">...</span>
			<a href="#page10" class="page-link">10</a>
			<a href="#page2" class="page-link">Next</a>
		</nav>
	</section>
</div>`
}

func (l *Links) Style() string {
	return `
.demo { max-width: 700px; }
.demo h1 { color: #1a1a2e; margin-bottom: 20px; }
.demo section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; background: #fff; }
.demo h2 { margin-top: 0; color: #666; font-size: 1.1em; }

.link-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; }
.link-card { display: flex; flex-direction: column; padding: 16px; border: 1px solid #ddd; border-radius: 8px; text-decoration: none; color: inherit; transition: all 0.2s; }
.link-card:hover { border-color: #007bff; box-shadow: 0 2px 8px rgba(0,123,255,0.15); }
.link-card strong { color: #007bff; margin-bottom: 4px; }
.link-card span { font-size: 13px; color: #666; }
.link-card .link-url { margin-top: 8px; font-size: 11px; color: #999; }

.tracking-info { display: flex; gap: 20px; padding: 15px; background: #f8f9fa; border-radius: 8px; }
.stat { display: flex; flex-direction: column; align-items: center; flex: 1; }
.stat-value { font-size: 24px; font-weight: bold; color: #007bff; }
.stat-label { font-size: 12px; color: #666; margin-top: 4px; }

.button-links { display: flex; gap: 10px; flex-wrap: wrap; }
.btn-link { display: inline-block; padding: 10px 20px; border-radius: 4px; text-decoration: none; font-weight: 500; transition: all 0.2s; }
.btn-link.primary { background: #007bff; color: white; }
.btn-link.primary:hover { background: #0056b3; }
.btn-link.secondary { background: #6c757d; color: white; }
.btn-link.secondary:hover { background: #545b62; }
.btn-link.outline { border: 2px solid #007bff; color: #007bff; background: transparent; }
.btn-link.outline:hover { background: #007bff; color: white; }

.icon-links { display: flex; flex-direction: column; gap: 8px; }
.icon-link { display: flex; align-items: center; gap: 12px; padding: 10px 16px; border: 1px solid #eee; border-radius: 4px; text-decoration: none; color: #333; transition: all 0.2s; }
.icon-link:hover { background: #f8f9fa; border-color: #007bff; }
.icon-link .icon { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; background: #e9ecef; border-radius: 50%; font-weight: bold; }

.breadcrumb { display: flex; align-items: center; gap: 8px; padding: 10px 16px; background: #f8f9fa; border-radius: 4px; }
.breadcrumb a { color: #007bff; text-decoration: none; }
.breadcrumb a:hover { text-decoration: underline; }
.breadcrumb .separator { color: #999; }
.breadcrumb .current { color: #666; font-weight: 500; }

.pagination { display: flex; gap: 4px; justify-content: center; }
.page-link { display: flex; align-items: center; justify-content: center; min-width: 36px; height: 36px; padding: 0 12px; border: 1px solid #ddd; border-radius: 4px; text-decoration: none; color: #333; transition: all 0.2s; }
.page-link:hover:not(.disabled):not(.active):not(.ellipsis) { background: #f0f0f0; border-color: #007bff; }
.page-link.active { background: #007bff; color: white; border-color: #007bff; }
.page-link.disabled { color: #999; cursor: not-allowed; }
.page-link.ellipsis { border: none; cursor: default; }
`
}
