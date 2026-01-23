package main

import "reactive"

// PriceResponse is the typed response from the crypto API
type PriceResponse struct {
	RAW struct {
		PRICE      float64 `js:"PRICE"`
		FROMSYMBOL string  `js:"FROMSYMBOL"`
		TOSYMBOL   string  `js:"TOSYMBOL"`
		LASTUPDATE int     `js:"LASTUPDATE"`
	} `js:"RAW"`
}

// Bitcoin is the Bitcoin price tracker component
type Bitcoin struct {
	Price      *reactive.Store[string]
	Symbol     *reactive.Store[string]
	UpdateTime *reactive.Store[string]
	Loading    *reactive.Store[bool]
	Error      *reactive.Store[string]

	stopRefresh func() // cleanup function for interval
}

func (b *Bitcoin) OnMount() {
	b.Loading.Set(true)
	b.Error.Set("")
	b.Price.Set("")
	b.Symbol.Set("")
	b.UpdateTime.Set("")

	// Fetch initial price
	b.FetchPrice()

	// Set up 60-second refresh interval with cleanup
	b.stopRefresh = reactive.SetInterval(60000, func() {
		b.FetchPrice()
	})
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
		resp, err := reactive.Get[PriceResponse]("https://min-api.cryptocompare.com/data/generateAvg?fsym=BTC&tsym=USD&e=coinbase")
		if err != nil {
			b.Error.Set("Failed to fetch: " + err.Error())
			b.Loading.Set(false)
			return
		}

		raw := resp.RAW

		b.Price.Set(raw.TOSYMBOL + " " + formatPrice(raw.PRICE))
		b.Symbol.Set(raw.FROMSYMBOL)

		// Format time (LASTUPDATE is unix timestamp)
		secs := raw.LASTUPDATE % 86400 // seconds since midnight UTC
		h, m, s := secs/3600, (secs%3600)/60, secs%60
		b.UpdateTime.Set(pad2(h) + ":" + pad2(m) + ":" + pad2(s))

		b.Loading.Set(false)
	}()
}

func (b *Bitcoin) Retry() {
	b.FetchPrice()
}

func pad2(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
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
	// Format float with 2 decimal places without fmt
	neg := f < 0
	if neg {
		f = -f
	}
	cents := int(f*100 + 0.5) // round to cents
	dollars := cents / 100
	rem := cents % 100
	s := itoa(dollars) + "." + pad2(rem)
	if neg {
		return "-" + s
	}
	return s
}

func (b *Bitcoin) Template() string {
	return `<div class="bitcoin-container">
	<h2>Bitcoin Price Tracker</h2>

	<div class="bitcoin-card">
		{#if Loading}
			<p class="loading">Loading...</p>
		{:else if Error}
			<p class="error">Error: {Error}</p>
			<button class="retry-btn" @click="Retry()">Retry</button>
		{:else}
			<div class="price-info">
				<span class="symbol">{Symbol}</span>
				<span class="update-time">Updated: {UpdateTime}</span>
			</div>
			<p class="price">{Price}</p>
			<small class="disclaimer">
				Prices are volatile and for reference only. Not financial advice.
			</small>
		{/if}
	</div>

	<p class="refresh-note">Price refreshes automatically every 60 seconds</p>
</div>`
}

func (b *Bitcoin) Style() string {
	return `
.bitcoin-container { max-width: 500px; margin: 2rem auto; padding: 1rem; }
.bitcoin-container h2 { text-align: center; margin-bottom: 1.5rem; color: #1a1a2e; }
.bitcoin-card { background: #fff; padding: 2rem; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); text-align: center; border: 1px solid #e9ecef; }
.price-info { display: flex; justify-content: space-between; margin-bottom: 1rem; color: #666; }
.symbol { font-weight: 600; color: #f7931a; }
.update-time { font-size: 0.875rem; }
.price { font-size: 2.5rem; font-weight: bold; margin: 1rem 0; color: #1a1a2e; }
.disclaimer { display: block; color: #888; margin-top: 1rem; padding-top: 1rem; border-top: 1px solid #e9ecef; font-size: 0.8rem; }
.loading { color: #666; font-size: 1.1rem; }
.error { color: #e53e3e; margin-bottom: 1rem; }
.retry-btn { padding: 0.5rem 1.5rem; background: #e53e3e; color: #fff; border: none; border-radius: 4px; cursor: pointer; font-size: 0.9rem; }
.retry-btn:hover { background: #c53030; }
.refresh-note { text-align: center; color: #888; font-size: 0.85rem; margin-top: 1rem; font-style: italic; }
`
}
