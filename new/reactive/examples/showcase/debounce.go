package main

import "reactive"

type Debounce struct {
	SearchInput   *reactive.Store[string]
	SearchResult  *reactive.Store[string]
	SearchCount   *reactive.Store[int]
	ClickCount    *reactive.Store[int]
	ThrottleCount *reactive.Store[int]
	Status        *reactive.Store[string]

	doSearch      func()
	throttleClick func()
}

func (d *Debounce) OnMount() {
	d.SearchInput.Set("")
	d.SearchResult.Set("")
	d.SearchCount.Set(0)
	d.ClickCount.Set(0)
	d.ThrottleCount.Set(0)
	d.Status.Set("Type to search...")

	d.doSearch = reactive.Debounce(300, func() {
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

	d.throttleClick = reactive.Throttle(500, func() {
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

func (d *Debounce) Template() string {
	return `<div class="demo">
	<h1>Debounce & Throttle</h1>

	<section>
		<h2>Debounced Search</h2>
		<p>Search triggers 300ms after you stop typing.</p>

		<input type="text" bind:value="SearchInput" placeholder="Type to search..." />

		<div class="stats">
			<span>Status: <strong>{Status}</strong></span>
			<span>API calls: <strong>{SearchCount}</strong></span>
		</div>

		{#if SearchResult}
			<div class="result">{SearchResult}</div>
		{/if}

		<p class="hint">Type quickly - search only fires once you pause.</p>
	</section>

	<section>
		<h2>Throttled Clicks</h2>
		<p>Button action throttled to max once per 500ms.</p>

		<button @click="OnClick()">Click me rapidly!</button>

		<div class="stats">
			<span>Total clicks: <strong>{ClickCount}</strong></span>
			<span>Throttled actions: <strong>{ThrottleCount}</strong></span>
		</div>

		<p class="hint">Click fast - throttled count increases slowly.</p>
	</section>

	<section>
		<button @click="Reset()">Reset All</button>
	</section>
</div>`
}

func (d *Debounce) Style() string {
	return `
.demo { max-width: 600px; }
.demo h1 { color: #1a1a2e; margin-bottom: 20px; }
.demo section { margin: 20px 0; padding: 15px; border: 1px solid #ddd; border-radius: 8px; background: #fff; }
.demo h2 { margin-top: 0; color: #666; font-size: 1.1em; }
.demo button { padding: 8px 16px; margin: 4px; cursor: pointer; border: 1px solid #ccc; border-radius: 4px; background: #f5f5f5; }
.demo button:hover { background: #e5e5e5; }
.demo input[type="text"] { width: 100%; padding: 12px; border: 1px solid #ccc; border-radius: 4px; font-size: 16px; }
.stats { display: flex; gap: 20px; margin: 15px 0; padding: 10px; background: #e3f2fd; border-radius: 4px; }
.stats span { color: #1565c0; }
.result { padding: 15px; background: #e8f5e9; border-radius: 4px; color: #2e7d32; margin-top: 10px; }
.hint { font-size: 0.9em; color: #666; font-style: italic; margin-top: 10px; }
`
}
