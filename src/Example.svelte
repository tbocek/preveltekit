<script lang="ts">
    type BitcoinPrice = {
        RAW: {
            MARKET: string;
            FROMSYMBOL: string;
            TOSYMBOL: string;
            FLAGS: number;
            PRICE: number;
            LASTUPDATE: number;
            LASTVOLUME: number;
            LASTVOLUMETO: number;
            LASTTRADEID: string;
            VOLUME24HOUR: number;
            VOLUME24HOURTO: number;
            OPEN24HOUR: number;
            HIGH24HOUR: number;
            LOW24HOUR: number;
            LASTMARKET: string;
            TOPTIERVOLUME24HOUR: number;
            TOPTIERVOLUME24HOURTO: number;
            CHANGE24HOUR: number;
            CHANGEPCT24HOUR: number;
            CHANGEDAY: number;
            CHANGEPCTDAY: number;
            CHANGEHOUR: number;
            CHANGEPCTHOUR: number;
        };
        DISPLAY: {
            FROMSYMBOL: string;
            TOSYMBOL: string;
            MARKET: string;
            PRICE: string;
            LASTUPDATE: string;
            LASTVOLUME: string;
            LASTVOLUMETO: string;
            LASTTRADEID: string;
            VOLUME24HOUR: string;
            VOLUME24HOURTO: string;
            OPEN24HOUR: string;
            HIGH24HOUR: string;
            LOW24HOUR: string;
            LASTMARKET: string;
            TOPTIERVOLUME24HOUR: string;
            TOPTIERVOLUME24HOURTO: string;
            CHANGE24HOUR: string;
            CHANGEPCT24HOUR: string;
            CHANGEDAY: string;
            CHANGEPCTDAY: string;
            CHANGEHOUR: string;
            CHANGEPCTHOUR: string;
        };
    }

    let priceData = $state<BitcoinPrice | null>(null);
    let loading = $state(true);
    let error = $state<string | null>(null);

    // Demo the SSPR capability
    let renderInfo = $state("Client Rendered");
    if (window?.JSDOM) {
        renderInfo = "Server Pre-Rendered";
    }

    async function fetchBitcoinPrice() {
        try {
            loading = true;
            error = null;
            const response = await fetch('https://min-api.cryptocompare.com/data/generateAvg?fsym=BTC&tsym=USD&e=coinbase');
            if (!response.ok) throw new Error('Failed to fetch data');
            priceData = await response.json();
        } catch (e) {
            error = e instanceof Error ? e.message : 'An error occurred';
        } finally {
            loading = false;
        }
    }

    // Fetch initial data
    $effect(() => {
        if (!window?.JSDOM) {
            fetchBitcoinPrice();
            // Set up refresh interval
            const interval = setInterval(fetchBitcoinPrice, 60000); // Update every minute
            return () => clearInterval(interval);
        }
    });
</script>

<div class="bitcoin-dashboard">
    <div class="header">
        <h2>Bitcoin Price Tracker</h2>
        <p class="render-info">({renderInfo})</p>
    </div>

    <div class="price-display">
        {#if loading && !priceData}
            <div class="loading">Loading Bitcoin prices...</div>
        {:else if error}
            <div class="error">
                Error: {error}
                <button onclick={fetchBitcoinPrice}>Retry</button>
            </div>
        {:else if priceData}
            <div class="price-card">
                <div class="price-header">
                    <span class="currency-code">{@html priceData.RAW.FROMSYMBOL}</span>
                    <span class="update-time">Last Updated: {priceData.RAW.LASTUPDATE}</span>
                </div>
                <div class="current-price">
                    {@html priceData.RAW.TOSYMBOL}
                    {priceData.RAW.PRICE}
                </div>
            </div>
            <div class="disclaimer">
                Cryptocurrency prices are highly volatile and subject to market risks. The displayed price information is for reference only and may not reflect real-time market conditions. Past performance is not indicative of future results. Please conduct your own research and consider your financial situation before making any investment decisions.
            </div>
        {/if}
    </div>
</div>

<style>
    .bitcoin-dashboard {
        max-width: 800px;
        margin: 0 auto;
        padding: 2rem;
    }

    .header {
        text-align: center;
        margin-bottom: 2rem;
    }

    h2 {
        color: #2d3748;
        margin: 0;
        font-size: 2rem;
    }

    .render-info {
        color: #718096;
        font-size: 0.875rem;
        margin-top: 0.5rem;
    }

    .price-display {
        background: white;
        border-radius: 12px;
        box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        overflow: hidden;
    }

    .price-card {
        padding: 2rem;
    }

    .price-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .currency-code {
        font-size: 1.25rem;
        font-weight: bold;
        color: #2d3748;
    }

    .update-time {
        font-size: 0.875rem;
        color: #718096;
    }

    .current-price {
        font-size: 3rem;
        font-weight: bold;
        color: #2d3748;
        margin: 1rem 0;
        text-align: center;
    }

    .disclaimer {
        padding: 1rem;
        background: #f7fafc;
        color: #718096;
        font-size: 0.875rem;
        text-align: center;
        border-top: 1px solid #e2e8f0;
    }

    .loading {
        padding: 2rem;
        text-align: center;
        color: #4a5568;
    }

    .error {
        padding: 2rem;
        text-align: center;
        color: #e53e3e;
    }

    .error button {
        margin-top: 1rem;
        padding: 0.5rem 1rem;
        background: #e53e3e;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }

    .error button:hover {
        background: #c53030;
    }
</style>