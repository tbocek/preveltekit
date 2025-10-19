<script lang="ts">
    // Demo SPAwBR capabilities
    import SSR from "./assets/SSR.svg";
    import SPA from "./assets/SPA.svg";
    import SPAwBR from "./assets/SPAwBR.svg";
    let renderType = "Client Rendered";
    if (window?.__isBuildTime) {
        renderType = "Server Pre-Rendered";
    }
</script>

<div class="docs">
    <header>
        <h1>PrevelteKit Documentation</h1>
        <p class="render-info">({renderType})</p>
    </header>

    <section>
        <h2>Quick Start</h2>
        Make sure you have node/npm installed. Here is a minimalistic example:
        <div class="code-block">
            {@html `<pre><code>mkdir -p preveltekit/src && cd preveltekit
echo '&lcub;"devDependencies": &lcub;"preveltekit": "^1.2.12"&rcub;,"dependencies": &lcub;"svelte": "^5.39.11"&rcub;,"scripts": &lcub;"dev": "preveltekit dev"&rcub;&rcub;' &gt; package.json
npm install
echo '&lt;script&gt;let count = $state(0);&lt;/script&gt;&lt;h1&gt;Count: &lcub;count&rcub;&lt;/h1&gt;&lt;button onclick={() =&gt; count++}&gt;Click me&lt;/button&gt;' &gt; src/Index.svelte
npm run dev
# And open a browser with localhost:3000</code></pre>`}
        </div>
    </section>

    <section>
        <h2>Core Features</h2>

        <h3>
            ⚡ Single Page Application with Built-time Pre-rendering (SPAwBR)
        </h3>
        <p>
            PrevelteKit combines the best of SPA and build-time rendering with
            hydration approaches. Unlike traditional SSR that renders on each
            request, or pure SPA that shows blank pages initially, SPAwBR
            pre-renders your layout and static content at build time while
            maintaining full interactivity through hydration. This provides fast
            initial page loads with visible content, then progressive
            enhancement as JavaScript loads.
        </p>

        <h3>🎯 Simple Architecture</h3>
        <p>
            Built on a clear separation between frontend and backend. Your
            frontend is purely static assets (HTML/CSS/JS) that can be served
            from any CDN or web server, while data comes from dedicated API
            endpoints. No JavaScript runtime required for serving.
        </p>

        <h3>⚡ Lightning Fast Builds</h3>
        <p>
            Built on Rsbuild for builds in the range of hundreds of
            milliseconds. The system automatically handles:
        </p>
        <ul>
            <li>TypeScript compilation and type checking</li>
            <li>Asset optimization and bundling</li>
            <li>CSS processing and minification</li>
            <li>Pre-compression (Brotli, Zstandard, Gzip)</li>
        </ul>

        <h3>🔧 Development Workflow</h3>
        <p>Three modes available to suit your needs:</p>
        <ul>
            <li>
                <strong>Development (npm run dev)</strong>: Express server with
                fast rebuilds and live reloading
            </li>
            <li>
                <strong>Staging (npm run stage)</strong>: Production build with
                local preview server
            </li>
            <li>
                <strong>Production (npm run build)</strong>: Optimized build
                with pre-compression for deployment
            </li>
        </ul>
    </section>

    <section>
        <h2>Rendering Comparison</h2>

        <div class="comparison-table">
            <table>
                <thead>
                    <tr>
                        <th>Rendering Type</th>
                        <th>Initial Load</th>
                        <th>After Script</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td
                            ><strong>SSR</strong> (classic SSR / Next.js / Nuxt)</td
                        >
                        <td>
                            <img
                                src={SSR}
                                alt="SSR Initial"
                                class="comparison-img"
                            />
                            <br />User sees content instantly
                            <br /><small>(rendered on each request)</small>
                        </td>
                        <td
                            ><img
                                src={SSR}
                                alt="SSR After"
                                class="comparison-img"
                            /><br />User sees content instantly<br /><small
                                >(no additional loading)</small
                            ></td
                        >
                    </tr>
                    <tr>
                        <td><strong>SPA</strong> (React App / pure Svelte)</td>
                        <td
                            ><img
                                src={SPA}
                                alt="SPA Initial"
                                class="comparison-img"
                            /><br />User sees white page or spinner<br /><small
                                >(no content until JS loads)</small
                            ></td
                        >
                        <td
                            ><img
                                src={SSR}
                                alt="SPA Loaded"
                                class="comparison-img"
                            /><br />User sees full content<br /><small
                                >(after script execution)</small
                            ></td
                        >
                    </tr>
                    <tr>
                        <td
                            ><strong>SPA + Build-time Pre-Rendering</strong> (this
                            approach)</td
                        >
                        <td
                            ><img
                                src={SPAwBR}
                                alt="SPAwBR Initial"
                                class="comparison-img"
                            /><br />User sees layout and static content<br
                            /><small>(pre-rendered at build time)</small></td
                        >
                        <td
                            ><img
                                src={SSR}
                                alt="SPAwBR Hydrated"
                                class="comparison-img"
                            /><br />User sees interactive content<br /><small
                                >(hydrated with full functionality)</small
                            ></td
                        >
                    </tr>
                </tbody>
            </table>
        </div>
    </section>

    <section>
        <h2>SPAwBR Development</h2>

        <h3>🔍 Detecting Build-time Pre-rendering</h3>
        <p>
            PrevelteKit uses <code>window.__isBuildTime</code> to indicate when code
            is running during build-time pre-rendering. This is crucial for handling
            client-side-only code like API calls and intervals.
        </p>

        <div class="code-block">
            {@html `<pre><code>// Basic detection
let renderInfo = "Client Rendered";
if (window?.__isBuildTime) &lcub;
    renderInfo = "Server Pre-Rendered";
&rcub;</code></pre>`}
        </div>

        <h3>🔄 Handling Client-Side Operations</h3>
        <p>
            PrevelteKit automatically handles fetch requests during build-time
            pre-rendering. Fetch calls made during pre-rendering will timeout
            after 5 seconds, allowing your components to render with loading
            states. You no longer need to wrap fetch calls in <code
                >window.__isBuildTime</code
            > checks.
        </p>

        <div class="code-block">
            {@html `<pre><code>// Fetch automatically handled during pre-rendering
let pricePromise = $state(fetchBitcoinPrice());

// Use Svelte's await block for clean handling
&lcub;#await pricePromise&rcub;
    &lt;p&gt;Loading...&lt;/p&gt;
&lcub;:then data&rcub;
    &lt;p&gt;&lcub;data&rcub;&lt;/p&gt;
&lcub;:catch error&rcub;
    &lt;p&gt;Error: &lcub;error.message&rcub;&lt;/p&gt;
&lcub;/await&rcub;</code></pre>`}
        </div>

        <p>
            <strong>When to still use build-time checks:</strong>
        </p>
        <ul>
            <li>Browser APIs (localStorage, sessionStorage, geolocation)</li>
            <li>DOM manipulation that shouldn't happen during pre-rendering</li>
            <li>Third-party scripts that expect a real browser environment</li>
        </ul>
    </section>

    <section>
        <h2>Configuration</h2>
        <p>
            PrevelteKit uses <code>rsbuild.config.ts</code> for configuration
            with sensible defaults. To customize settings, create an
            <code>rsbuild.config.ts</code> file in your project - it will merge with
            the default configuration.
        </p>

        <p>
            The framework provides fallback files (<code>index.html</code> and
            <code>index.ts</code>) from the default folder when you don't supply
            your own. Once you add your own files, PrevelteKit uses those
            instead, ignoring the defaults.
        </p>
    </section>

    <section>
        <h2>Client-Side Routing</h2>
        <p>
            PrevelteKit includes a built-in routing system that handles
            navigation between different pages in your application. The router
            uses pattern matching to determine which component to render based
            on the current URL path.
        </p>

        <h3>🧭 Route Configuration</h3>
        <p>
            Define your routes as an array of route objects, each specifying a
            path pattern, the component to render, and the static HTML file
            name:
        </p>
        <div class="code-block">
            <pre><code
                    >const routes: Routes = &lcub;
     dynamicRoutes: [
         &lcub;
             path: "*/doc",
             component: Documentation
         &rcub;,
         &lcub;
             path: "*/example",
             component: Example
         &rcub;,
         &lcub;
             path: "*/",
             component: Landing
         &rcub;
     ],
     staticRoutes: [
         &lcub;
             path: "/doc",
             htmlFilename: "doc.html"
         &rcub;,
         &lcub;
             path: "/example",
             htmlFilename: "example.html"
         &rcub;,
         &lcub;
             path: "/",
             htmlFilename: "index.html"
         &rcub;
     ]
 &rcub;;

 &lt;Router routes&gt;</code
                ></pre>
        </div>

        <h3>🔍 Path Patterns</h3>
        <p>PrevelteKit supports flexible path patterns for routing:</p>
        <ul>
            <li>
                <strong>Wildcard prefix (<code>*/path</code>)</strong>: Matches
                any single segment before the path (e.g., <code>*/doc</code>
                matches <code>/doc</code> and <code>/any/doc</code>)
            </li>
            <li>
                <strong>Root wildcard (<code>*/</code>)</strong>: Matches the
                root path and single-segment paths
            </li>
            <li>
                <strong>Exact paths (<code>/about</code>)</strong>: Matches the
                exact path only
            </li>
            <li>
                <strong>Parameters (<code>/user/:id</code>)</strong>: Captures
                URL segments as parameters
            </li>
        </ul>

        <h3>🔗 Navigation</h3>
        <p>
            Use the <code>route</code> action for client-side navigation that updates
            the URL without page reloads:
        </p>
        <div class="code-block">
            <pre><code
                    >import &lcub; route &rcub; from 'preveltekit';

    &lt;a use:link href="doc"&gt;Documentation&lt;/a&gt;
    &lt;a use:link href="example"&gt;Example&lt;/a&gt;</code
                ></pre>
        </div>

        <h3>📄 Static File Mapping & Hybrid Routing</h3>
        <p>
            The <code>staticRoutes</code> array configuration serves a dual purpose
            in PrevelteKit's hybrid routing approach:
        </p>
        <div class="code-block">
            <pre><code
                    >htmlFilename: "doc.html"  // Generates dist/doc.html at build time</code
                ></pre>
        </div>

        <p>
            <strong>Static Generation:</strong> During the build process,
            PrevelteKit generates actual HTML files in your <code>dist/</code> folder
            for each route:
        </p>
        <ul>
            <li><code>dist/index.html</code> - Pre-rendered root route</li>
            <li>
                <code>dist/doc.html</code> - Pre-rendered documentation page
            </li>
            <li><code>dist/example.html</code> - Pre-rendered example page</li>
        </ul>

        <p>
            <strong>Dynamic Routing:</strong> Once the application loads, the same
            route configuration enables client-side navigation between pages without
            full page reloads. This provides:
        </p>
        <ul>
            <li>Fast initial page loads from pre-rendered static HTML</li>
            <li>Instant navigation between routes via client-side routing</li>
            <li>
                SEO benefits from static HTML while maintaining SPA
                functionality
            </li>
        </ul>

        <p>
            This hybrid approach means users get static HTML files for direct
            access (bookmarks, search engines) and dynamic routing for seamless
            navigation within the application.
        </p>

        <h3>⚙️ Route Matching Priority</h3>
        <p>
            Routes are matched based on specificity, with more specific patterns
            taking precedence:
        </p>
        <ol>
            <li>Exact path matches (highest priority)</li>
            <li>Parameter-based routes</li>
            <li>Wildcard patterns (lowest priority)</li>
        </ol>
        <p>
            Always place more specific routes before general wildcard routes in
            your configuration to ensure proper matching behavior.
        </p>
    </section>

    <section>
        <h2>Docker Support</h2>
        <p>Development environment:</p>
        <div class="code-block">
            <pre><code
                    >docker build -f Dockerfile.dev . -t preveltekit-dev
docker run -p3000:3000 -v./src:/app/src preveltekit-dev</code
                ></pre>
        </div>

        <p>Production build:</p>
        <div class="code-block">
            <pre><code
                    >docker build . -t preveltekit
docker run -p3000:3000 preveltekit</code
                ></pre>
        </div>
    </section>

    <section>
        <h2>Architecture Philosophy</h2>
        <p>
            PrevelteKit emphasizes <strong>static-first architecture</strong> with
            clear separation between frontend and backend:
        </p>
        <ul>
            <li>
                <strong>Frontend</strong>: Pure static assets (HTML/CSS/JS)
                served from any web server or CDN
            </li>
            <li>
                <strong>Backend</strong>: Dedicated API endpoints for data, can
                be built with any technology
            </li>
            <li>
                <strong>Deployment</strong>: No JavaScript runtime required -
                just static files
            </li>
        </ul>

        <p>
            This approach offers compelling simplicity compared to full-stack
            meta-frameworks:
        </p>
        <ul>
            <li>Deploy anywhere (GitHub Pages, S3, any web server)</li>
            <li>Predictable performance with no server processes to monitor</li>
            <li>Easier debugging with clear boundaries</li>
            <li>Freedom to choose your backend technology</li>
        </ul>
    </section>

    <section>
        <h2>Deployment</h2>
        <p>
            The production build generates static files with pre-compressed
            variants:
        </p>
        <ul>
            <li>Standard files (.js, .css, .html)</li>
            <li>Brotli compressed (.br)</li>
            <li>Gzip compressed (.gz)</li>
            <li>Zstandard compressed (.zst)</li>
        </ul>
        <p>
            Deploy to any static hosting or web server. The pre-compressed files
            enable optimal performance when served with appropriate web server
            configuration.
        </p>
    </section>

    <section>
        <h2>Why PrevelteKit?</h2>
        <p>
            While SvelteKit provides comprehensive capabilities, PrevelteKit
            focuses on a minimalistic solution for build-time pre-rendering.
            With less than 500 lines of code, it's essentially glue code for
            Svelte, Rsbuild, and jsdom - perfect for projects that need fast
            initial loads without the complexity of full JavaScript
            infrastructure for the frontend deployment.
        </p>

        <p>
            PrevelteKit serves as a starting point for projects that need
            pre-rendered content without the overhead of a full meta-framework,
            following a "convention over configuration" approach.
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
        font-family: monospace;
        font-size: 0.8rem;
    }

    .code-block {
        background: #1a202c;
        border-radius: 8px;
        padding: 1rem;
        margin: 1rem 0;
        overflow-x: auto;
        color: #e2e8f0;
        font-family: monospace;
        font-size: 0.8rem;
    }

    .code-block pre {
        margin: 0;
        color: #e2e8f0 !important;
        font-family: monospace;
        font-size: 0.8rem;
    }

    .code-block code {
        color: #e2e8f0 !important;
        background: transparent !important;
        padding: 0;
        font-family: monospace;
        font-size: 0.8rem;
        line-height: 1.5;
    }

    .comparison-table {
        margin: 1rem 0;
        overflow-x: auto;
    }

    table {
        width: 100%;
        border-collapse: collapse;
        background: white;
        border-radius: 8px;
        overflow: hidden;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }

    th,
    td {
        padding: 1rem;
        text-align: left;
        border-bottom: 1px solid #e2e8f0;
    }

    th {
        background: #f7fafc;
        font-weight: 600;
        color: #2d3748;
    }

    td {
        color: #4a5568;
        vertical-align: top;
    }

    small {
        color: #718096;
        display: block;
        margin-top: 0.25rem;
    }

    .comparison-img {
        width: 120px;
        height: auto;
        margin-bottom: 0.5rem;
        border-radius: 4px;
    }
</style>
