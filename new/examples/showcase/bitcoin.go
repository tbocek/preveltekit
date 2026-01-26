package main

import "preveltekit"

type PriceResponse struct {
	RAW struct {
		PRICE      float64 `js:"PRICE"`
		FROMSYMBOL string  `js:"FROMSYMBOL"`
		TOSYMBOL   string  `js:"TOSYMBOL"`
		LASTUPDATE int     `js:"LASTUPDATE"`
	} `js:"RAW"`
}

type Bitcoin struct {
	Price       *preveltekit.Store[string]
	Symbol      *preveltekit.Store[string]
	UpdateTime  *preveltekit.Store[string]
	Loading     *preveltekit.Store[bool]
	Error       *preveltekit.Store[string]
	stopRefresh func()
}

func (b *Bitcoin) OnCreate() {
	// Called once - start the refresh timer
	b.FetchPrice()

	b.stopRefresh = preveltekit.SetInterval(60000, func() {
		b.FetchPrice()
	})
}

func (b *Bitcoin) OnMount() {
	// Called every navigation - nothing needed here since we cache data
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
		resp, err := preveltekit.Get[PriceResponse]("https://min-api.cryptocompare.com/data/generateAvg?fsym=BTC&tsym=USD&e=coinbase")
		if err != nil {
			b.Error.Set("Failed to fetch: " + err.Error())
			b.Loading.Set(false)
			return
		}

		raw := resp.RAW

		b.Price.Set(raw.TOSYMBOL + " " + formatPrice(raw.PRICE))
		b.Symbol.Set(raw.FROMSYMBOL)

		secs := raw.LASTUPDATE % 86400
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

func (b *Bitcoin) Template() string {
	return `<div class="demo">
	<h1>Bitcoin Price</h1>

	<section class="bitcoin-card">
		{#if Loading}
			<p class="loading">Loading...</p>
		{:else if Error}
			<p class="error">Error: {Error}</p>
			<button @click="Retry()">Retry</button>
		{:else}
			<div class="price-info">
				<span class="symbol">{Symbol}</span>
				<span class="update-time">Updated: {UpdateTime} UTC</span>
			</div>
			<p class="price">{Price}</p>
			<small class="disclaimer">
				Prices are volatile and for reference only.
			</small>
		{/if}
	</section>

	<p class="hint">Price refreshes automatically every 60 seconds</p>
</div>`
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
