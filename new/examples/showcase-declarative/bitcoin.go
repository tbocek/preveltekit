package main

import p "preveltekit"

type PriceResponse struct {
	RAW struct {
		PRICE      float64 `js:"PRICE"`
		FROMSYMBOL string  `js:"FROMSYMBOL"`
		TOSYMBOL   string  `js:"TOSYMBOL"`
		LASTUPDATE int     `js:"LASTUPDATE"`
	} `js:"RAW"`
}

type Bitcoin struct {
	Price       *p.Store[string]
	Symbol      *p.Store[string]
	UpdateTime  *p.Store[string]
	Loading     *p.Store[bool]
	Error       *p.Store[string]
	stopRefresh func()
}

func (b *Bitcoin) OnCreate() {
	if p.IsBuildTime {
		// During SSR, show loading state - the actual fetch happens in WASM
		b.Loading.Set(true)
		b.Error.Set("")
		return
	}

	b.FetchPrice()

	b.stopRefresh = p.SetInterval(60000, func() {
		b.FetchPrice()
	})
}

func (b *Bitcoin) OnMount() {
}

func (b *Bitcoin) OnDestroy() {
	if b.stopRefresh != nil {
		b.stopRefresh()
	}
}

func (b *Bitcoin) FetchPrice() {
	b.Loading.Set(true)
	b.Error.Set("")

	go func() {
		resp, err := p.Get[PriceResponse]("https://min-api.cryptocompare.com/data/generateAvg?fsym=BTC&tsym=USD&e=coinbase")
		if err != nil {
			b.Error.Set("Failed to fetch: " + err.Error())
			b.Loading.Set(false)
			return
		}

		raw := resp.RAW

		println("DEBUG: PRICE=", raw.PRICE, "SYMBOL=", raw.FROMSYMBOL, "LASTUPDATE=", raw.LASTUPDATE)

		b.Price.Set(raw.TOSYMBOL + " " + formatPrice(raw.PRICE))
		b.Symbol.Set(raw.FROMSYMBOL)

		// Convert Unix timestamp to UTC time
		secs := raw.LASTUPDATE % 86400
		h, m, s := secs/3600, (secs%3600)/60, secs%60
		b.UpdateTime.Set(pad2(h) + ":" + pad2(m) + ":" + pad2(s))
		println("DEBUG: time=", pad2(h)+":"+pad2(m)+":"+pad2(s))

		b.Loading.Set(false)
	}()
}

func (b *Bitcoin) Retry() {
	b.FetchPrice()
}

func pad2(n int) string {
	if n < 10 {
		return "0" + btcItoa(n)
	}
	return btcItoa(n)
}

func btcItoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(byte('0'+n%10)) + s
		n /= 10
	}
	return s
}

func formatPrice(f float64) string {
	neg := f < 0
	if neg {
		f = -f
	}
	cents := int(f*100 + 0.5)
	dollars := cents / 100
	rem := cents % 100
	s := btcItoa(dollars) + "." + pad2(rem)
	if neg {
		return "-" + s
	}
	return s
}

func (b *Bitcoin) Render() p.Node {
	return p.Div(p.Class("demo"),
		p.H1("Bitcoin Price"),

		p.Section(p.Class("bitcoin-card"),
			p.If(p.IsTrue(b.Loading),
				p.P(p.Class("loading"), "Loading..."),
			).ElseIf(b.Error.Ne(""),
				p.P(p.Class("error"), "Error: ", p.Bind(b.Error)),
				p.Button("Retry", p.OnClick(b.Retry)),
			).Else(
				p.Div(p.Class("price-info"),
					p.Span(p.Class("symbol"), p.Bind(b.Symbol)),
					p.Span(p.Class("update-time"), "Updated: ", p.Bind(b.UpdateTime), " UTC"),
				),
				p.P(p.Class("price"), p.Bind(b.Price)),
				p.Small(p.Class("disclaimer"),
					"Prices are volatile and for reference only.",
				),
			),
		),

		p.P(p.Class("hint"), "Price refreshes automatically every 60 seconds"),
	)
}

func (b *Bitcoin) Style() string {
	return `
.demo{max-width:500px}
.bitcoin-card{background:#fff;padding:2rem;border-radius:8px;border:1px solid #ddd;text-align:center}
.price-info{display:flex;justify-content:space-between;margin-bottom:1rem;color:#666}
.symbol{font-weight:600;color:#f7931a}
.update-time{font-size:.875rem}
.price{font-size:2.5rem;font-weight:700;margin:1rem 0;color:#1a1a2e}
.disclaimer{display:block;color:#888;margin-top:1rem;padding-top:1rem;border-top:1px solid #e9ecef;font-size:.8rem}
.loading{color:#666;font-size:1.1rem}
.error{color:#e53e3e;margin-bottom:1rem}
`
}

func (b *Bitcoin) HandleEvent(method string, args string) {
	switch method {
	case "Retry":
		b.Retry()
	}
}
