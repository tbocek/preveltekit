package main

import p "github.com/tbocek/preveltekit/v2"

type PriceResponse struct {
	RAW struct {
		PRICE      float64 `js:"PRICE"`
		FROMSYMBOL string  `js:"FROMSYMBOL"`
		TOSYMBOL   string  `js:"TOSYMBOL"`
		LASTUPDATE int     `js:"LASTUPDATE"`
	} `js:"RAW"`
}

type BitcoinDemo struct {
	Price       *p.Store[string]
	Symbol      *p.Store[string]
	UpdateTime  *p.Store[string]
	Loading     *p.Store[bool]
	Error       *p.Store[string]
	stopRefresh func()
}

func (b *BitcoinDemo) New() p.Component {
	return &BitcoinDemo{
		Price:      p.New(""),
		Symbol:     p.New(""),
		UpdateTime: p.New(""),
		Loading:    p.New(true),
		Error:      p.New(""),
	}
}

func (b *BitcoinDemo) OnMount() {
	if p.IsBuildTime {
		return
	}

	b.FetchPrice()

	b.stopRefresh = p.SetInterval(60000, func() {
		b.FetchPrice()
	})
}

func (b *BitcoinDemo) OnDestroy() {
	if b.stopRefresh != nil {
		b.stopRefresh()
	}
}

func (b *BitcoinDemo) FetchPrice() {
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
		b.Price.Set(raw.TOSYMBOL + " " + btcFormatPrice(raw.PRICE))
		b.Symbol.Set(raw.FROMSYMBOL)

		secs := raw.LASTUPDATE % 86400
		h, m, s := secs/3600, (secs%3600)/60, secs%60
		b.UpdateTime.Set(btcPad2(h) + ":" + btcPad2(m) + ":" + btcPad2(s))

		b.Loading.Set(false)
	}()
}

func (b *BitcoinDemo) Retry() {
	b.FetchPrice()
}

func btcPad2(n int) string {
	if n < 10 {
		return "0" + p.Itoa(n)
	}
	return p.Itoa(n)
}

func btcFormatPrice(f float64) string {
	neg := f < 0
	if neg {
		f = -f
	}
	cents := int(f*100 + 0.5)
	dollars := cents / 100
	rem := cents % 100
	s := p.Itoa(dollars) + "." + btcPad2(rem)
	if neg {
		return "-" + s
	}
	return s
}

func (b *BitcoinDemo) Render() p.Node {
	return p.Div(p.Attr("class", "btc-page page"),
		p.Div(p.Attr("class", "container"),
			p.H1("Bitcoin Price"),
			p.P(p.Attr("class", "page-intro"), "Live price fetched from CryptoCompare API with auto-refresh every 60 seconds."),

			p.Div(p.Attr("class", "btc-card"),
				p.If(p.Cond(func() bool { return b.Loading.Get() }, b.Loading),
					p.P(p.Attr("class", "loading"), "Loading..."),
				).ElseIf(p.Cond(func() bool { return b.Error.Get() != "" }, b.Error),
					p.P(p.Attr("class", "error"), b.Error),
					p.Button(p.Attr("class", "retry-btn"), "Retry").On("click", b.Retry),
				).Else(
					p.Div(p.Attr("class", "price-header"),
						p.Span(p.Attr("class", "symbol"), b.Symbol),
						p.Span(p.Attr("class", "update-time"), "Updated: ", b.UpdateTime, " UTC"),
					),
					p.P(p.Attr("class", "price"), b.Price),
					p.Small(p.Attr("class", "disclaimer"), "Prices are volatile and for reference only."),
				),
			),

			p.Div(p.Attr("class", "btc-code"),
				p.H2("How it works"),
				p.P(p.RawHTML("This demo uses <code>p.Get[T]()</code> to fetch JSON, <code>p.SetInterval()</code> for auto-refresh, and lifecycle hooks for setup/teardown.")),
				p.Pre(p.Code(`type BitcoinDemo struct {
    Price       *p.Store[string]
    Loading     *p.Store[bool]
    Error       *p.Store[string]
    stopRefresh func()
}

func (b *BitcoinDemo) OnMount() {
    if p.IsBuildTime { return }
    b.FetchPrice()
    b.stopRefresh = p.SetInterval(60000, func() {
        b.FetchPrice()
    })
}

func (b *BitcoinDemo) OnDestroy() {
    if b.stopRefresh != nil {
        b.stopRefresh()
    }
}

func (b *BitcoinDemo) FetchPrice() {
    b.Loading.Set(true)
    go func() {
        resp, err := p.Get[PriceResponse](url)
        if err != nil {
            b.Error.Set(err.Error())
            b.Loading.Set(false)
            return
        }
        b.Price.Set(formatPrice(resp.RAW.PRICE))
        b.Loading.Set(false)
    }()
}`)),
			),
		),
	)
}

func (b *BitcoinDemo) Style() string {
	return `
.btc-card{background:#fff;padding:24px;border-radius:8px;border:1px solid #e5e7eb;text-align:center;max-width:500px;margin-bottom:32px}
.price-header{display:flex;justify-content:space-between;margin-bottom:16px;color:#666}
.symbol{font-weight:600;color:#f7931a;font-size:1.1em}
.update-time{font-size:14px}
.price{font-size:2.5em;font-weight:700;margin:16px 0;color:#1a1a2e}
.disclaimer{display:block;color:#999;margin-top:16px;padding-top:16px;border-top:1px solid #e9ecef;font-size:.8em}
.loading{color:#666;font-size:1.1em;padding:32px 0}
.error{color:#e53e3e;margin-bottom:16px}
.retry-btn{padding:8px 16px;background:#1a1a2e;color:#fff;border:none;border-radius:4px;cursor:pointer}
.retry-btn:hover{background:#0f3460}

.btc-code{max-width:700px}
.btc-code h2{font-size:1.3em;color:#1a1a2e;margin-bottom:8px}
.btc-code > p{color:#666;margin-bottom:16px;font-size:.95em}
`
}
