<script lang="ts">
    type CurrencyData = {
        code: string;
        symbol: string;
        rate: string;
        description: string;
        rate_float: number;
    };

    type BitcoinPrice = {
        time: {
            updated: string;
            updatedISO: string;
            updateduk: string;
        };
        disclaimer: string;
        bpi: Record<'USD' | 'GBP' | 'EUR', CurrencyData>;
    };

    let priceData = $state<BitcoinPrice | null>(null);
    let loading = $state(true);
    let error = $state<string | null>(null);
    let selectedCurrency = $state<'USD' | 'GBP' | 'EUR'>('USD');

    // Demo the SSPR capability
    let renderInfo = $state("Client Rendered");
    if (window?.JSDOM) {
        renderInfo = "Server Pre-Rendered";
    }

    async function fetchBitcoinPrice() {
        try {
            loading = true;
            error = null;
            const response = await fetch('https://api.coindesk.com/v1/bpi/currentprice.json');
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

    <div class="currency-selector">
        <label for="currency">Select Currency:</label>
        <select
                id="currency"
                bind:value={selectedCurrency}
                class="select-input"
        >
            <option value="USD">USD</option>
            <option value="EUR">EUR</option>
            <option value="GBP">GBP</option>
        </select>
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
                    <span class="currency-code">{selectedCurrency}</span>
                    <span class="update-time">Last Updated: {priceData.time.updated}</span>
                </div>
                <div class="current-price">
                    {@html priceData.bpi[selectedCurrency].symbol}
                    {priceData.bpi[selectedCurrency].rate}
                </div>
                <div class="price-description">
                    {priceData.bpi[selectedCurrency].description}
                </div>
            </div>
            <div class="disclaimer">
                {priceData.disclaimer}
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

    .currency-selector {
        margin-bottom: 2rem;
        text-align: center;
    }

    .select-input {
        padding: 0.5rem 1rem;
        border: 2px solid #e2e8f0;
        border-radius: 8px;
        margin-left: 1rem;
        font-size: 1rem;
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

    .price-description {
        text-align: center;
        color: #4a5568;
        margin-bottom: 1rem;
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