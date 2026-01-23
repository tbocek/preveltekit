//go:build !js || !wasm

package main

import "reactive"

// Bitcoin stub for SSR/pre-rendering (no JS runtime)
type Bitcoin struct {
	Price      *reactive.Store[string]
	Symbol     *reactive.Store[string]
	UpdateTime *reactive.Store[string]
	Loading    *reactive.Store[bool]
	Error      *reactive.Store[string]
}

func (b *Bitcoin) OnMount() {
	// Show loading state in SSR
	b.Loading.Set(true)
	b.Error.Set("")
	b.Price.Set("")
	b.Symbol.Set("")
	b.UpdateTime.Set("")
}

func (b *Bitcoin) OnDestroy() {}

func (b *Bitcoin) FetchPrice() {}

func (b *Bitcoin) Retry() {}

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
