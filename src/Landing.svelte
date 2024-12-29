<script lang="ts">
    import { route } from "@mateothegreat/svelte5-router";

    let message = "Welcome to PrevelteKit";
    if (window?.JSDOM) {
        message = "Server-Side Pre-Rendered with PrevelteKit, you see this in the source code, but not in the DOM after loading";
    }
</script>

<div class="hero">
    <h1>{message}</h1>
    <p class="subtitle">A Modern Framework for Building Fast, SEO-Friendly Web Applications</p>

    <div class="cta-buttons">
        <a use:route href="./doc" class="cta-button primary">View Documentation</a>
        <a use:route href="./example" class="cta-button secondary">Try Bitcoin Demo</a>
    </div>

    <div class="features">
        <div class="feature">
            <h3>⚡️ Lightning Fast</h3>
            <p>Pre-rendered pages with hydration for optimal performance</p>
        </div>

        <div class="feature">
            <h3>📈 Real-time Data</h3>
            <p>Seamless integration with external APIs and live updates</p>
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
        background: linear-gradient(to bottom, #f7fafc, #edf2f7);
    }

    h1 {
        font-size: 3rem;
        margin-bottom: 1rem;
        color: #2d3748;
        font-weight: bold;
    }

    .subtitle {
        font-size: 1.5rem;
        color: #4a5568;
        margin-bottom: 3rem;
        max-width: 800px;
        margin-left: auto;
        margin-right: auto;
    }

    .cta-buttons {
        margin-bottom: 4rem;
        display: flex;
        gap: 1rem;
        justify-content: center;
    }

    .cta-button {
        padding: 0.75rem 1.5rem;
        border-radius: 8px;
        font-weight: 600;
        transition: transform 0.2s;
        text-decoration: none;
    }

    .cta-button:hover {
        transform: translateY(-2px);
    }

    .primary {
        background: #4299e1;
        color: white;
    }

    .secondary {
        background: white;
        color: #4299e1;
        border: 2px solid #4299e1;
    }

    .features {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
        gap: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .feature {
        padding: 2rem;
        background: white;
        border-radius: 12px;
        box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        transition: transform 0.2s;
    }

    .feature:hover {
        transform: translateY(-5px);
    }

    .feature h3 {
        color: #2d3748;
        margin-bottom: 1rem;
        font-size: 1.25rem;
    }

    .feature p {
        color: #4a5568;
        line-height: 1.6;
    }

    @media (max-width: 640px) {
        .cta-buttons {
            flex-direction: column;
            padding: 0 2rem;
        }

        h1 {
            font-size: 2.5rem;
        }

        .subtitle {
            font-size: 1.25rem;
        }
    }
</style>