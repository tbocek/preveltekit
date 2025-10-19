<script lang="ts">
    import { link } from "preveltekit";
    let message = "PrevelteKit";
    let subtitle = "Minimalistic Build-time Pre-rendering for Svelte";

    if (window.__isBuildTime) {
        message = "PrevelteKit (Pre-rendered)";
        subtitle =
            "This content was pre-rendered at build time and will hydrate on load";
    }
</script>

<div class="hero">
    <div class="hero-content">
        <h1>{message}</h1>
        <p class="subtitle">{subtitle}</p>

        <ul class="key-benefits">
            <li class="benefit">
                <strong>Static-first architecture</strong> - Deploy frontend anywhere as pure HTML/CSS/JS
            </li>
            <li class="benefit">
                <strong>Lightning fast builds</strong> - Hundreds of milliseconds with Rsbuild
            </li>
            <li class="benefit">
                <strong>No JavaScript runtime in production for frontend</strong> - Just static files on a server
            </li>
        </ul>

        <div class="cta-buttons">
            <a use:link href="doc" class="cta-button primary">Documentation</a>
            <a use:link href="example" class="cta-button secondary"
                >Bitcoin Demo</a
            >
        </div>

        <div class="routing-note">
            (These buttons here are client-side routing • No page reloads)
        </div>
    </div>

    <div class="features-section">
        <div class="feature-grid">
            <div class="feature">
                <div class="feature-icon">⚡</div>
                <h3>Lightning Fast Builds</h3>
                <p>Rsbuild bundles in hundreds of milliseconds</p>
            </div>

            <div class="feature">
                <div class="feature-icon">🎯</div>
                <h3>Minimalistic</h3>
                <p>
                    Less than 500 lines of code - just glue for Svelte, Rsbuild,
                    and jsdom
                </p>
            </div>

            <div class="feature">
                <div class="feature-icon">🔄</div>
                <h3>HTML for Initial Page</h3>
                <p>
                    Single Page App with Built-time Pre-Rendering - best of both
                    worlds
                </p>
            </div>

            <div class="feature">
                <div class="feature-icon">🚀</div>
                <h3>Deploy Anywhere</h3>
                <p>GitHub Pages, S3, any web server - just static files</p>
            </div>

            <div class="feature">
                <div class="feature-icon">🛡️</div>
                <h3>Clear Separation</h3>
                <p>Frontend as static assets, backend as dedicated APIs</p>
            </div>

            <div class="feature">
                <div class="feature-icon">📦</div>
                <h3>Zero Config</h3>
                <p>Works out of the box with sensible defaults</p>
            </div>
        </div>
    </div>

    <div class="quick-start">
        <h2>Quick Start</h2>
        <div class="code-block">
            {@html `<pre><code>mkdir -p preveltekit/src && cd preveltekit
echo '&lcub;"devDependencies": &lcub;"preveltekit": "^1.2.9"&rcub;,"dependencies": &lcub;"svelte": "^5.39.11"&rcub;,"scripts": &lcub;"dev": "preveltekit dev"&rcub;&rcub;' &gt; package.json
npm install
echo '&lt;script&gt;let count = $state(0);&lt;/script&gt;&lt;h1&gt;Count: &lcub;count&rcub;&lt;/h1&gt;&lt;button onclick={() =&gt; count++}&gt;Click me&lt;/button&gt;' &gt; src/Index.svelte
npm run dev
# And open a browser with localhost:3000</code></pre>`}
        </div>
    </div>
</div>

<style>
    .hero {
        padding: 2rem 1rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    .hero-content {
        text-align: center;
        margin-bottom: 4rem;
    }

    h1 {
        font-size: 3.5rem;
        margin-bottom: 1rem;
        color: #2d3748;
        font-weight: bold;
    }

    .subtitle {
        font-size: 1.25rem;
        color: #4a5568;
        margin-bottom: 2rem;
        max-width: 600px;
        margin-left: auto;
        margin-right: auto;
    }

    .key-benefits {
        list-style: disc;
        padding-left: 1.5rem;
        margin: 2rem auto;
        max-width: 600px;
        text-align: left;
        display: inline-block;
    }

    .benefit {
        color: #4a5568;
        font-size: 1.1rem;
        margin-bottom: 0.75rem;
    }

    .benefit strong {
        color: #2d3748;
    }

    .cta-buttons {
        margin: 2rem 0;
        display: flex;
        gap: 1rem;
        justify-content: center;
        flex-wrap: wrap;
    }

    .cta-button {
        padding: 0.875rem 2rem;
        border-radius: 8px;
        font-weight: 600;
        font-size: 1.1rem;
        transition: all 0.2s;
        text-decoration: none;
        display: inline-block;
    }

    .cta-button:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
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

    .routing-note {
        color: #718096;
        font-size: 0.875rem;
        margin-top: 1rem;
    }

    .features-section {
        margin: 4rem 0;
    }

    .feature-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
        gap: 2rem;
    }

    .feature {
        background: white;
        padding: 2rem;
        border-radius: 12px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        text-align: center;
        transition: transform 0.2s;
    }

    .feature:hover {
        transform: translateY(-4px);
    }

    .feature-icon {
        font-size: 3rem;
        margin-bottom: 1rem;
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

    .quick-start {
        margin: 4rem 0;
        text-align: center;
    }

    .quick-start h2 {
        color: #2d3748;
        font-size: 2rem;
        margin-bottom: 2rem;
    }

    .code-block {
        text-align: left;
        background: #1a202c;
        border-radius: 8px;
        padding: 1rem;
        margin: 1rem 0;
        overflow-x: auto;
        color: #e2e8f0;
        font-family: monospace;
        font-size: 0.8rem;
    }

    @media (max-width: 768px) {
        h1 {
            font-size: 2.5rem;
        }

        .subtitle {
            font-size: 1.1rem;
        }

        .cta-buttons {
            flex-direction: column;
            align-items: center;
        }

        .cta-button {
            width: 200px;
        }
    }
</style>
