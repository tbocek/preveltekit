<script lang="ts">
    type BitcoinPrice = {
        RAW: {
            FROMSYMBOL: string;
            TOSYMBOL: string;
            PRICE: number;
            LASTUPDATE: number;
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
        const response = await fetch('https://min-api.cryptocompare.com/data/generateAvg?fsym=BTC&tsym=USD&e=coinbase');
        if (!response.ok) throw new Error('Failed to fetch data');
        return response.json();
    }
    
    let pricePromise = $state<Promise<BitcoinPrice>>(fetchBitcoinPrice());
    
    // Set up refresh interval
    $effect(() => {
        const interval = setInterval(() => {
            pricePromise = fetchBitcoinPrice();
        }, 60000);
        return () => clearInterval(interval);
    });
    
    function retry() {
        pricePromise = fetchBitcoinPrice();
    }
</script>

<div class="container">
    <h2>Bitcoin Price Tracker <small>({renderInfo})</small></h2>
    
    <div class="card">
        {#await pricePromise}
            <p>Loading...</p>
        {:then priceData}
            <div class="price-info">
                <span>{priceData.RAW.FROMSYMBOL}</span>
                <time>Updated: {new Date(priceData.RAW.LASTUPDATE * 1000).toLocaleTimeString()}</time>
            </div>
            <p class="price">{priceData.RAW.TOSYMBOL} {priceData.RAW.PRICE.toFixed(2)}</p>
            <small class="disclaimer">
                Prices are volatile and for reference only. Not financial advice.
            </small>
        {:catch error}
            <p class="error">Error: {error.message}</p>
            <button onclick={retry}>Retry</button>
        {/await}
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