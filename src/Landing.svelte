<script lang="ts">
    let message = "Welcome to LightKit";
    if (window?.JSDOM) {
        message = "Server-Side Pre-Rendered with LightKit";
    }
</script>

<div class="hero">
    <h1>{message}</h1>
    <p class="subtitle">Lightweight, Fast, and SEO-Friendly Web Applications</p>

    <div class="features">
        <div class="feature">
            <h3>⚡️ Lightning Fast</h3>
            <p>Pre-rendered pages with hydration for optimal performance</p>
        </div>

        <div class="feature">
            <h3>🎯 Simple Routing</h3>
            <p>Built-in routing system with pre and post hooks</p>
        </div>

        <div class="feature">
            <h3>🔄 SSPR Support</h3>
            <p>Server-Side Pre-Rendering for better SEO</p>
        </div>
    </div>
</div>

<style>
    .hero {
        text-align: center;
        padding: 4rem 2rem;
    }

    h1 {
        font-size: 2.5rem;
        margin-bottom: 1rem;
        color: #2d3748;
    }

    .subtitle {
        font-size: 1.25rem;
        color: #4a5568;
        margin-bottom: 3rem;
    }

    .features {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .feature {
        padding: 1.5rem;
        background: #fff;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }

    .feature h3 {
        color: #2d3748;
        margin-bottom: 0.5rem;
    }

    .feature p {
        color: #4a5568;
        line-height: 1.5;
    }
</style>