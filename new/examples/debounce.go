package main

import p "preveltekit"

type Debounce struct {
	SearchInput   *p.Store[string]
	SearchResult  *p.Store[string]
	SearchCount   *p.Store[int]
	ClickCount    *p.Store[int]
	ThrottleCount *p.Store[int]
	Status        *p.Store[string]

	doSearch        func()
	cleanupDebounce func()
	throttleClick   func()
}

func (d *Debounce) OnMount() {
	d.SearchInput.Set("")
	d.SearchResult.Set("")
	d.SearchCount.Set(0)
	d.ClickCount.Set(0)
	d.ThrottleCount.Set(0)
	d.Status.Set("Type to search...")

	d.doSearch, d.cleanupDebounce = p.Debounce(300, func() {
		query := d.SearchInput.Get()
		if query == "" {
			d.SearchResult.Set("")
			d.Status.Set("Type to search...")
			return
		}

		d.SearchCount.Set(d.SearchCount.Get() + 1)
		d.SearchResult.Set("Results for: " + query)
		d.Status.Set("Search complete!")
	})

	d.SearchInput.OnChange(func(_ string) {
		d.Status.Set("Waiting...")
		d.doSearch()
	})

	d.throttleClick = p.Throttle(500, func() {
		d.ThrottleCount.Set(d.ThrottleCount.Get() + 1)
	})
}

func (d *Debounce) OnClick() {
	d.ClickCount.Set(d.ClickCount.Get() + 1)
	d.throttleClick()
}

func (d *Debounce) Reset() {
	d.SearchInput.Set("")
	d.SearchResult.Set("")
	d.SearchCount.Set(0)
	d.ClickCount.Set(0)
	d.ThrottleCount.Set(0)
	d.Status.Set("Type to search...")
}

func (d *Debounce) Render() p.Node {
	return p.Html(`<div class="demo">
		<h1>Debounce &amp; Throttle</h1>

		<section>
			<h2>Debounced Search</h2>
			<p>Search triggers 300ms after you stop typing.</p>

			`, p.BindValue(`<input type="text" placeholder="Type to search...">`, d.SearchInput), `

			<div class="stats">
				<span>Status: <strong>`, p.Bind(d.Status), `</strong></span>
				<span>API calls: <strong>`, p.Bind(d.SearchCount), `</strong></span>
			</div>`,

		p.If(d.SearchResult.Ne(""),
			p.Html(`<div class="result">`, p.Bind(d.SearchResult), `</div>`),
		),

		p.Html(`<p class="hint">Type quickly - search only fires once you pause.</p>
		</section>

		<section>
			<h2>Throttled Clicks</h2>
			<p>Button action throttled to max once per 500ms.</p>

			`, p.Html(`<button>Click me rapidly!</button>`).WithOn("click", d.OnClick), `

			<div class="stats">
				<span>Total clicks: <strong>`, p.Bind(d.ClickCount), `</strong></span>
				<span>Throttled actions: <strong>`, p.Bind(d.ThrottleCount), `</strong></span>
			</div>

			<p class="hint">Click fast - throttled count increases slowly.</p>
		</section>

		<section>
			`, p.Html(`<button>Reset All</button>`).WithOn("click", d.Reset), `
		</section>
	</div>`),
	)
}

func (d *Debounce) Style() string {
	return `
.demo input[type=text]{width:100%;padding:12px;font-size:16px}
.stats{display:flex;gap:20px;margin:15px 0;padding:10px;background:#e3f2fd;border-radius:4px}
.stats span{color:#1565c0}
.result{padding:15px;background:#e8f5e9;border-radius:4px;color:#2e7d32;margin-top:10px}
`
}
