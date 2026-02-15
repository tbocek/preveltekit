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
	return p.Html(`<div class="btc-page">
	<div class="container">
		<h1>Bitcoin Price</h1>
		<p class="subtitle">Live price fetched from CryptoCompare API with auto-refresh every 60 seconds.</p>

		<div class="btc-card">`,
		p.If(p.Cond(func() bool { return b.Loading.Get() }, b.Loading),
			p.Html(`<p class="loading">Loading...</p>`),
		).ElseIf(p.Cond(func() bool { return b.Error.Get() != "" }, b.Error),
			p.Html(`<p class="error">`, b.Error, `</p>`,
				p.Html(`<button class="retry-btn">Retry</button>`).On("click", b.Retry)),
		).Else(
			p.Html(`<div class="price-header">
					<span class="symbol">`, b.Symbol, `</span>
					<span class="update-time">Updated: `, b.UpdateTime, ` UTC</span>
				</div>
				<p class="price">`, b.Price, `</p>
				<small class="disclaimer">Prices are volatile and for reference only.</small>`),
		),
		p.Html(`</div>

		<div class="btc-code">
			<h2>How it works</h2>
			<p>This demo uses <code>p.Get[T]()</code> to fetch JSON, <code>p.SetInterval()</code> for auto-refresh, and lifecycle hooks for setup/teardown.</p>
			<pre><code>type BitcoinDemo struct {
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
}</code></pre>
		</div>
	</div>
	</div>`),
	)
}

func (b *BitcoinDemo) Style() string {
	return `
.btc-page{padding:40px 0}
.btc-page h1{font-size:2.2em;color:#1a1a2e;margin-bottom:8px}
.subtitle{color:#666;margin-bottom:32px;font-size:1.05em}

.btc-card{background:#fff;padding:2rem;border-radius:8px;border:1px solid #e5e7eb;text-align:center;max-width:500px;margin-bottom:32px}
.price-header{display:flex;justify-content:space-between;margin-bottom:1rem;color:#666}
.symbol{font-weight:600;color:#f7931a;font-size:1.1em}
.update-time{font-size:.875rem}
.price{font-size:2.5rem;font-weight:700;margin:1rem 0;color:#1a1a2e}
.disclaimer{display:block;color:#888;margin-top:1rem;padding-top:1rem;border-top:1px solid #e9ecef;font-size:.8rem}
.loading{color:#666;font-size:1.1rem;padding:2rem 0}
.error{color:#e53e3e;margin-bottom:1rem}
.retry-btn{padding:8px 16px;background:#1a1a2e;color:#fff;border:none;border-radius:4px;cursor:pointer}
.retry-btn:hover{background:#0f3460}

.btc-code{max-width:700px}
.btc-code h2{font-size:1.3em;color:#1a1a2e;margin-bottom:8px}
.btc-code > p{color:#555;margin-bottom:16px;font-size:.95em}
.btc-code code{background:#f1f5f9;padding:2px 6px;border-radius:3px;font-size:.85em}
.btc-code pre{background:#1a1a2e;color:#e0e0e0;padding:16px;border-radius:6px;overflow-x:auto;font-size:13px;line-height:1.6}
.btc-code pre code{background:transparent;padding:0;font-size:inherit}
`
}
