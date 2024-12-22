<script lang="ts">
    let count = $state(0);
    let loading = $state(true);

    // Simulate data loading
    setTimeout(() => {
        loading = false;
    }, 1000);

    function increment() {
        count++;
    }

    // Demo SSPR capabilities
    let renderType = "client-side";
    if (window?.JSDOM) {
        renderType = "server-side";
    }
</script>

<div class="dashboard">
    <h1>LightKit Dashboard</h1>

    <div class="stats-grid">
        <div class="stat-card">
            <h3>Render Type</h3>
            <p>{renderType}</p>
        </div>

        <div class="stat-card">
            <h3>Counter Demo</h3>
            <p>{count}</p>
            <button onclick={increment}>Increment</button>
        </div>

        <div class="stat-card">
            <h3>Loading State</h3>
            {#if loading}
                <p>Loading...</p>
            {:else}
                <p>Data Loaded!</p>
            {/if}
        </div>
    </div>

    <div class="info-section">
        <h2>Framework Features</h2>
        <ul>
            <li>Server-Side Pre-Rendering (SSPR)</li>
            <li>Reactive State Management</li>
            <li>Component Hydration</li>
            <li>Development Hot Reload</li>
        </ul>
    </div>
</div>

<style>
    .dashboard {
        padding: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    h1 {
        color: #2d3748;
        margin-bottom: 2rem;
    }

    .stats-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 1.5rem;
        margin-bottom: 2rem;
    }

    .stat-card {
        background: white;
        padding: 1.5rem;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    button {
        background: #4299e1;
        color: white;
        border: none;
        padding: 0.5rem 1rem;
        border-radius: 4px;
        cursor: pointer;
    }

    button:hover {
        background: #3182ce;
    }

    .info-section {
        background: white;
        padding: 2rem;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    ul {
        list-style-type: none;
        padding: 0;
    }

    li {
        padding: 0.5rem 0;
        color: #4a5568;
    }

    h2 {
        color: #2d3748;
        margin-bottom: 1rem;
    }
</style>