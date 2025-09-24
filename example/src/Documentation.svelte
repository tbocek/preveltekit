<script lang="ts">
    // Demo SSPR capabilities
    let renderType = "Client Rendered";
    if (window?.JSDOM) {
        renderType = "Server Pre-Rendered";
    }
</script>

<div class="docs">
    <header>
        <h1>PrevelteKit Documentation</h1>
        <p class="render-info">({renderType})</p>
    </header>

    <section>
        <h2>Getting Started</h2>
        <div class="code-block">
            <pre><code>git clone https://github.com/tbocek/preveltekit.git
cd preveltekit
pnpm install
pnpm dev    # Development mode
pnpm build  # Production build</code></pre>
        </div>
    </section>

    <section>
        <h2>Core Features</h2>

        <h3>📦 Server-Side Pre-Rendering (SSPR)</h3>
        <p>
            Unlike traditional SSR or SSG, SSPR pre-renders your app at build time while maintaining full interactivity through hydration.
            This gives you the best of both worlds: fast initial loads and dynamic functionality.
        </p>

        <h3>⚡ Build System</h3>
        <p>
            Built on Rsbuild for lightning-fast builds. The system automatically handles:
        </p>
        <ul>
            <li>TypeScript compilation and type checking</li>
            <li>Asset optimization and bundling</li>
            <li>CSS processing and minification</li>
            <li>Production optimizations like code splitting</li>
        </ul>

        <h3>🔄 Development Workflow</h3>
        <p>
            Three modes available to suit your needs:
        </p>
        <ul>
            <li><strong>Development (pnpm dev)</strong>: Express server with fast rebuilds and hot module replacement</li>
            <li><strong>Staging (pnpm stage)</strong>: Production build with local preview server</li>
            <li><strong>Production (pnpm build)</strong>: Optimized build with pre-compression (Brotli, Gzip, Zstandard)</li>
        </ul>
    </section>

    <section>
        <h2>SSPR Development</h2>

        <h3>🔍 Detecting Server Pre-Rendering</h3>
        <p>
            PrevelteKit uses <code>window.JSDOM</code> to indicate when code is running during server pre-rendering.
            This is crucial for handling client-side-only code like API calls and intervals.
        </p>

        <div class="code-block">
            <pre><code>{"// Basic detection\nlet renderInfo = \"Client Rendered\";\nif (window?.JSDOM) {\n    renderInfo = \"Server Pre-Rendered\";\n}"}</code></pre>
        </div>

        <h3>🔄 Handling Client-Side Operations</h3>
        <p>
            When working with APIs, timers, or browser-specific features, wrap them in a JSDOM check to prevent execution during pre-rendering:
        </p>

        <div class="code-block">
            <pre><code>{"$effect(() => {\n    if (!window?.JSDOM) {\n        fetchBitcoinPrice();\n        // Set up refresh interval\n        const interval = setInterval(fetchBitcoinPrice, 60000);\n        return () => clearInterval(interval);\n    }\n});"}</code></pre>
        </div>

        <p>
            <strong>Common use cases for JSDOM checks:</strong>
        </p>
        <ul>
            <li>API calls and data fetching</li>
            <li>Browser APIs (localStorage, sessionStorage)</li>
            <li>Timers and intervals</li>
            <li>DOM manipulation</li>
            <li>Browser-specific features (geolocation, notifications)</li>
        </ul>
    </section>

    <section>
        <h2>Static Site Configuration</h2>
        <p>
            To add static routes to your application, configure them in <code>rsbuild.config.ts</code>.
            Each route needs an entry in the configuration. For example, to add routes for
            <code>/doc</code> and <code>/example</code>:
        </p>
        <div class="code-block">
            <pre><code>{"export default defineConfig({\n    environments: {\n        web: {\n            plugins: [\n                pluginSvelte(),\n                pluginCssMinimizer()\n            ],\n            source: {\n                entry: {\n                    // Each entry corresponds to a static route\n                    index: './src/index.ts',    // https://example.com/\n                    doc: './src/index.ts',      // https://example.com/doc\n                    example: './src/index.ts',  // https://example.com/example\n                }\n            },\n            output: {\n                target: 'web',\n                minify: process.env.NODE_ENV === 'production',\n            }\n        }\n    }\n});"}</code></pre>
        </div>
        <p>
            <strong>Important Notes:</strong>
        </p>
        <ul>
            <li>Each entry key becomes a URL path in your application</li>
            <li>Subdirectories are not currently supported</li>
            <li>All entries point to the same <code>index.ts</code> file</li>
            <li>The router in your application will handle the actual component rendering</li>
        </ul>
    </section>

    <section>
        <h2>Docker Support</h2>
        <p>Development environment:</p>
        <div class="code-block">
            <pre><code>docker build -f Dockerfile.dev . -t preveltekit-dev
docker run -p3000:3000 -v./src:/app/src -v./public:/app/public preveltekit-dev</code></pre>
        </div>

        <p>Production build:</p>
        <div class="code-block">
            <pre><code>docker build . -t preveltekit
docker run -p3000:3000 preveltekit</code></pre>
        </div>
    </section>

    <section>
        <h2>Project Structure</h2>
        <ul>
            <li><strong>/src</strong>: Application source code</li>
            <li><strong>/public</strong>: Static assets</li>
            <li><strong>/dist</strong>: Production build output</li>
            <li><strong>rsbuild.config.ts</strong>: Build configuration</li>
            <li><strong>ssr.mjs</strong>: SSPR implementation</li>
        </ul>
    </section>

    <section>
        <h2>Deployment</h2>
        <p>
            The production build generates static files with pre-compressed variants:
        </p>
        <ul>
            <li>Standard files (.js, .css, .html)</li>
            <li>Brotli compressed (.br)</li>
            <li>Gzip compressed (.gz)</li>
            <li>Zstandard compressed (.zst)</li>
        </ul>
        <p>
            Deploy to any static hosting or web server that supports serving compressed assets.
        </p>
    </section>
</div>

<style>
    .docs {
        max-width: 800px;
        margin: 0 auto;
        padding: 2rem 1rem;
    }

    header {
        margin-bottom: 3rem;
        text-align: center;
    }

    h1 {
        color: #2d3748;
        font-size: 2.5rem;
        margin-bottom: 0.5rem;
    }

    .render-info {
        color: #718096;
        font-size: 0.875rem;
    }

    section {
        margin-bottom: 3rem;
    }

    h2 {
        color: #2d3748;
        font-size: 1.8rem;
        margin-bottom: 1.5rem;
        padding-bottom: 0.5rem;
        border-bottom: 2px solid #e2e8f0;
    }

    h3 {
        color: #2d3748;
        font-size: 1.3rem;
        margin: 2rem 0 1rem;
    }

    p {
        color: #4a5568;
        line-height: 1.6;
        margin-bottom: 1rem;
    }

    ul {
        list-style-type: disc;
        padding-left: 1.5rem;
        margin: 1rem 0;
        color: #4a5568;
    }

    li {
        margin-bottom: 0.5rem;
        line-height: 1.6;
    }

    strong {
        color: #2d3748;
        font-weight: 600;
    }

    code {
        background: #edf2f7;
        padding: 0.2rem 0.4rem;
        border-radius: 4px;
        font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
        font-size: 0.9rem;
    }

    .code-block {
        background: #1a202c;
        border-radius: 8px;
        padding: 1.5rem;
        margin: 1rem 0;
        overflow-x: auto;
    }

    .code-block pre {
        margin: 0;
    }

    .code-block code {
        color: #e2e8f0;
        background: transparent;
        padding: 0;
        font-family: 'Menlo', 'Monaco', 'Courier New', monospace;
        font-size: 0.9rem;
        line-height: 1.5;
    }
</style>