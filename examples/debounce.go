package main

import p "github.com/tbocek/preveltekit/v2"

type Debounce struct {
	SearchInput   *p.Store[string]
	SearchResult  *p.Store[string]
	SearchCount   *p.Store[int]
	ClickCount    *p.Store[int]
	ThrottleCount *p.Store[int]
	Status        *p.Store[string]
	TimerStatus   *p.Store[string]

	doSearch        func()
	cleanupDebounce func()
	throttleClick   func()
}

func (d *Debounce) New() p.Component {
	return &Debounce{
		SearchInput:   p.New(""),
		SearchResult:  p.New(""),
		SearchCount:   p.New(0),
		ClickCount:    p.New(0),
		ThrottleCount: p.New(0),
		Status:        p.New("Type to search..."),
		TimerStatus:   p.New("Ready"),
	}
}

func (d *Debounce) OnMount() {

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

func (d *Debounce) StartTimer() {
	d.TimerStatus.Set("Waiting 2 seconds...")
	p.SetTimeout(2000, func() {
		d.TimerStatus.Set("Timer fired!")
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
	return p.Div(p.Attr("class", "demo"),
		p.H1("Debounce & Throttle"),

		p.Section(
			p.H2("Debounced Search"),
			p.P("Search triggers 300ms after you stop typing."),
			p.Input(p.Attr("type", "text"), p.Attr("placeholder", "Type to search...")).Bind(d.SearchInput),
			p.Div(p.Attr("class", "stats"),
				p.Span("Status: ", p.Strong(d.Status)),
				p.Span("API calls: ", p.Strong(d.SearchCount)),
			),
			p.If(p.Cond(func() bool { return d.SearchResult.Get() != "" }, d.SearchResult),
				p.Div(p.Attr("class", "result"), d.SearchResult),
			),
			p.P(p.Attr("class", "hint"), "Type quickly - search only fires once you pause."),
		),

		p.Section(
			p.H2("Throttled Clicks"),
			p.P("Button action throttled to max once per 500ms."),
			p.Button("Click me rapidly!").On("click", d.OnClick),
			p.Div(p.Attr("class", "stats"),
				p.Span("Total clicks: ", p.Strong(d.ClickCount)),
				p.Span("Throttled actions: ", p.Strong(d.ThrottleCount)),
			),
			p.P(p.Attr("class", "hint"), "Click fast - throttled count increases slowly."),
		),

		p.Section(
			p.H2("SetTimeout \u2014 One-Shot Timer"),
			p.P("Fires once after a delay."),
			p.Button("Start 2s Timer").On("click", d.StartTimer),
			p.P("Timer: ", p.Strong(d.TimerStatus)),
		),

		p.Section(
			p.Button("Reset All").On("click", d.Reset),
		),

		p.Section(
			p.H2("Code"),
			p.Pre(p.Attr("class", "code"),
				`// debounce: fires after idle period
doSearch, cleanup := p.Debounce(300, func() {
    // fires 300ms after last call
})
doSearch()  // call repeatedly — only last one fires
cleanup()   // cancel pending

// throttle: fires at most once per interval
onClick := p.Throttle(500, func() {
    // max once per 500ms
})

// setTimeout: fires once after delay
cancel := p.SetTimeout(2000, func() {
    // fires after 2 seconds
})
cancel() // cancel before it fires

// setInterval: fires repeatedly
stop := p.SetInterval(60000, func() {
    // fires every 60 seconds
})
stop() // stop the interval`),
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
