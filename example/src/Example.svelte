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
    if (window?.__isBuildTime) {
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
        if (!window?.__isBuildTime) {
            fetchBitcoinPrice();
            // Set up refresh interval
            const interval = setInterval(fetchBitcoinPrice, 60000); // Update every minute
            return () => clearInterval(interval);
        }
    });
</script>

<div class="container">
    <h2>Bitcoin Price Tracker <small>({renderInfo})</small></h2>
    
    <div class="card">
        {#if loading && !priceData}
            <p>Loading...</p>
        {:else if error}
            <p class="error">Error: {error}</p>
            <button onclick={fetchBitcoinPrice}>Retry</button>
        {:else if priceData}
            <div class="price-info">
                <span>{priceData.RAW.FROMSYMBOL}</span>
                <time>Updated: {new Date(priceData.RAW.LASTUPDATE * 1000).toLocaleTimeString()}</time>
            </div>
            <p class="price">{priceData.RAW.TOSYMBOL} {priceData.RAW.PRICE.toFixed(2)}</p>
            <small class="disclaimer">
                Prices are volatile and for reference only. Not financial advice.
            </small>
        {/if}
    </div>
</div>

<style>
    .container {
        max-width: 600px;
        margin: 2rem auto;
        padding: 1rem;
    }
    
    h2 {
        text-align: center;
        margin-bottom: 1.5rem;
    }
    
    small {
        color: #666;
        font-size: 0.875rem;
    }
    
    .card {
        background: white;
        padding: 2rem;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        text-align: center;
    }
    
    .price-info {
        display: flex;
        justify-content: space-between;
        margin-bottom: 1rem;
    }
    
    .price {
        font-size: 2.5rem;
        font-weight: bold;
        margin: 1rem 0;
    }
    
    .disclaimer {
        display: block;
        color: #666;
        margin-top: 1rem;
        padding-top: 1rem;
        border-top: 1px solid #eee;
    }
    
    .error {
        color: #e53e3e;
    }
    
    button {
        margin-top: 1rem;
        padding: 0.5rem 1rem;
        background: #e53e3e;
        color: white;
        border: none;
        border-radius: 4px;
        cursor: pointer;
    }
</style>