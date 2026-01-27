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
	return p.Div(p.Class("demo"),
		p.H1("Debounce & Throttle"),

		p.Section(
			p.H2("Debounced Search"),
			p.P("Search triggers 300ms after you stop typing."),

			p.Input(p.Type("text"), p.BindValue(d.SearchInput), p.Placeholder("Type to search...")),

			p.Div(p.Class("stats"),
				p.Span("Status: ", p.Strong(p.Bind(d.Status))),
				p.Span("API calls: ", p.Strong(p.Bind(d.SearchCount))),
			),

			p.If(d.SearchResult.Ne(""),
				p.Div(p.Class("result"), p.Bind(d.SearchResult)),
			),

			p.P(p.Class("hint"), "Type quickly - search only fires once you pause."),
		),

		p.Section(
			p.H2("Throttled Clicks"),
			p.P("Button action throttled to max once per 500ms."),

			p.Button("Click me rapidly!", p.OnClick(d.OnClick)),

			p.Div(p.Class("stats"),
				p.Span("Total clicks: ", p.Strong(p.Bind(d.ClickCount))),
				p.Span("Throttled actions: ", p.Strong(p.Bind(d.ThrottleCount))),
			),

			p.P(p.Class("hint"), "Click fast - throttled count increases slowly."),
		),

		p.Section(
			p.Button("Reset All", p.OnClick(d.Reset)),
		),
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

func (d *Debounce) HandleEvent(method string, args string) {
	switch method {
	case "OnClick":
		d.OnClick()
	case "Reset":
		d.Reset()
	}
}
