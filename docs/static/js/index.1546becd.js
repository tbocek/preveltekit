(() => { // webpackBootstrap
"use strict";
var __webpack_modules__ = ({
529: (function (__unused_webpack_module, __unused_webpack___webpack_exports__, __webpack_require__) {

// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/index-client.js + 1 modules
var index_client = __webpack_require__(732);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/disclose-version.js + 1 modules
var disclose_version = __webpack_require__(999);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/flags/legacy.js
var legacy = __webpack_require__(306);
// EXTERNAL MODULE: ./node_modules/.pnpm/svelte@5.39.3/node_modules/svelte/src/internal/client/index.js + 50 modules
var client = __webpack_require__(750);
// EXTERNAL MODULE: ./node_modules/.pnpm/preveltekit@file+.._@rsbuild+core@1.5.10_@rsbuild+plugin-svelte@1.0.11_@rsbuild+core@1._94c2fb735f7985d6da51a252e56b6ee3/node_modules/preveltekit/dist/Router.svelte
var Router_svelte = __webpack_require__(966);
// EXTERNAL MODULE: ./node_modules/.pnpm/preveltekit@file+.._@rsbuild+core@1.5.10_@rsbuild+plugin-svelte@1.0.11_@rsbuild+core@1._94c2fb735f7985d6da51a252e56b6ee3/node_modules/preveltekit/dist/router-utils.js
var router_utils = __webpack_require__(521);
;// CONCATENATED MODULE: ./src/Landing.svelte.2.css!=!./node_modules/.pnpm/svelte-loader@3.2.4_svelte@5.39.3/node_modules/svelte-loader/index.js?cssPath=/home/draft/git/pv2/example/src/Landing.svelte.2.css!./src/Landing.svelte
// extracted by css-extract-rspack-plugin

;// CONCATENATED MODULE: ./src/Landing.svelte




var root = client/* .from_html */.vUu(`<div class="hero svelte-1nllfry"><h1 class="svelte-1nllfry"> </h1> <p class="subtitle svelte-1nllfry">A Modern Framework for Building Fast, SEO-Friendly Web Applications</p> <div class="cta-buttons svelte-1nllfry"><a href="doc" class="cta-button primary svelte-1nllfry">View Documentation</a> <a href="example" class="cta-button secondary svelte-1nllfry">Try Bitcoin Demo</a> <div class="server-side-indicator svelte-1nllfry">(these links are client-side routed links)</div></div> <div class="features svelte-1nllfry"><div class="feature svelte-1nllfry"><h3 class="svelte-1nllfry">⚡️ Lightning Fast</h3> <p class="svelte-1nllfry">Pre-rendered pages with hydration for optimal performance</p></div> <div class="feature svelte-1nllfry"><h3 class="svelte-1nllfry">📈 Real-time Data</h3> <p class="svelte-1nllfry">Seamless integration with external APIs and live updates</p></div> <div class="feature svelte-1nllfry"><h3 class="svelte-1nllfry">🔄 SSPR Support</h3> <p class="svelte-1nllfry">Server-Side Pre-Rendering for better SEO</p></div></div></div>`);
function Landing($$anchor) {
    var _window;
    let message = client/* .mutable_source */.zgK("Welcome to PrevelteKit");
    if ((_window = window) === null || _window === void 0 ? void 0 : _window.JSDOM) {
        client/* .set */.hZp(message, "Server-Side Pre-Rendered with PrevelteKit, you see this in the source code, you may see it flashing briefly, but you will not see this in the DOM after loading");
    }
    var div = root();
    var h1 = client/* .child */.jfp(div);
    var text = client/* .child */.jfp(h1, true);
    client/* .reset */.cLc(h1);
    var div_1 = client/* .sibling */.hg4(h1, 4);
    var a = client/* .child */.jfp(div_1);
    client/* .action */.XId(a, ($$node)=>router_utils/* .route */.w === null || router_utils/* .route */.w === void 0 ? void 0 : (0,router_utils/* .route */.w)($$node));
    var a_1 = client/* .sibling */.hg4(a, 2);
    client/* .action */.XId(a_1, ($$node)=>router_utils/* .route */.w === null || router_utils/* .route */.w === void 0 ? void 0 : (0,router_utils/* .route */.w)($$node));
    client/* .next */.K2T(2);
    client/* .reset */.cLc(div_1);
    client/* .next */.K2T(2);
    client/* .reset */.cLc(div);
    client/* .template_effect */.vNg(()=>client/* .set_text */.jax(text, client/* .get */.JtY(message)));
    client/* .append */.BCw($$anchor, div);
}


;// CONCATENATED MODULE: ./src/Documentation.svelte.1.css!=!./node_modules/.pnpm/svelte-loader@3.2.4_svelte@5.39.3/node_modules/svelte-loader/index.js?cssPath=/home/draft/git/pv2/example/src/Documentation.svelte.1.css!./src/Documentation.svelte
// extracted by css-extract-rspack-plugin

;// CONCATENATED MODULE: ./src/Documentation.svelte



var Documentation_svelte_root = client/* .from_html */.vUu(`<div class="docs svelte-an6hzx"><header class="svelte-an6hzx"><h1 class="svelte-an6hzx">PrevelteKit Documentation</h1> <p class="render-info svelte-an6hzx"> </p></header> <section class="svelte-an6hzx"><h2 class="svelte-an6hzx">Getting Started</h2> <div class="code-block svelte-an6hzx"><pre class="svelte-an6hzx"><code class="svelte-an6hzx">git clone https://github.com/tbocek/preveltekit.git
cd preveltekit
pnpm install
pnpm dev    # Development mode
pnpm build  # Production build</code></pre></div></section> <section class="svelte-an6hzx"><h2 class="svelte-an6hzx">Core Features</h2> <h3 class="svelte-an6hzx">📦 Server-Side Pre-Rendering (SSPR)</h3> <p class="svelte-an6hzx">Unlike traditional SSR or SSG, SSPR pre-renders your app at build time while maintaining full interactivity through hydration.
            This gives you the best of both worlds: fast initial loads and dynamic functionality.</p> <h3 class="svelte-an6hzx">⚡ Build System</h3> <p class="svelte-an6hzx">Built on Rsbuild for lightning-fast builds. The system automatically handles:</p> <ul class="svelte-an6hzx"><li class="svelte-an6hzx">TypeScript compilation and type checking</li> <li class="svelte-an6hzx">Asset optimization and bundling</li> <li class="svelte-an6hzx">CSS processing and minification</li> <li class="svelte-an6hzx">Production optimizations like code splitting</li></ul> <h3 class="svelte-an6hzx">🔄 Development Workflow</h3> <p class="svelte-an6hzx">Three modes available to suit your needs:</p> <ul class="svelte-an6hzx"><li class="svelte-an6hzx"><strong class="svelte-an6hzx">Development (pnpm dev)</strong>: Express server with fast rebuilds and hot module replacement</li> <li class="svelte-an6hzx"><strong class="svelte-an6hzx">Staging (pnpm stage)</strong>: Production build with local preview server</li> <li class="svelte-an6hzx"><strong class="svelte-an6hzx">Production (pnpm build)</strong>: Optimized build with pre-compression (Brotli, Gzip, Zstandard)</li></ul></section> <section class="svelte-an6hzx"><h2 class="svelte-an6hzx">SSPR Development</h2> <h3 class="svelte-an6hzx">🔍 Detecting Server Pre-Rendering</h3> <p class="svelte-an6hzx">PrevelteKit uses <code class="svelte-an6hzx">window.JSDOM</code> to indicate when code is running during server pre-rendering.
            This is crucial for handling client-side-only code like API calls and intervals.</p> <div class="code-block svelte-an6hzx"><pre class="svelte-an6hzx"><code class="svelte-an6hzx"></code></pre></div> <h3 class="svelte-an6hzx">🔄 Handling Client-Side Operations</h3> <p class="svelte-an6hzx">When working with APIs, timers, or browser-specific features, wrap them in a JSDOM check to prevent execution during pre-rendering:</p> <div class="code-block svelte-an6hzx"><pre class="svelte-an6hzx"><code class="svelte-an6hzx"></code></pre></div> <p class="svelte-an6hzx"><strong class="svelte-an6hzx">Common use cases for JSDOM checks:</strong></p> <ul class="svelte-an6hzx"><li class="svelte-an6hzx">API calls and data fetching</li> <li class="svelte-an6hzx">Browser APIs (localStorage, sessionStorage)</li> <li class="svelte-an6hzx">Timers and intervals</li> <li class="svelte-an6hzx">DOM manipulation</li> <li class="svelte-an6hzx">Browser-specific features (geolocation, notifications)</li></ul></section> <section class="svelte-an6hzx"><h2 class="svelte-an6hzx">Static Site Configuration</h2> <p class="svelte-an6hzx">To add static routes to your application, configure them in <code class="svelte-an6hzx">rsbuild.config.ts</code>.
            Each route needs an entry in the configuration. For example, to add routes for <code class="svelte-an6hzx">/doc</code> and <code class="svelte-an6hzx">/example</code>:</p> <div class="code-block svelte-an6hzx"><pre class="svelte-an6hzx"><code class="svelte-an6hzx"></code></pre></div> <p class="svelte-an6hzx"><strong class="svelte-an6hzx">Important Notes:</strong></p> <ul class="svelte-an6hzx"><li class="svelte-an6hzx">Each entry key becomes a URL path in your application</li> <li class="svelte-an6hzx">Subdirectories are not currently supported</li> <li class="svelte-an6hzx">All entries point to the same <code class="svelte-an6hzx">index.ts</code> file</li> <li class="svelte-an6hzx">The router in your application will handle the actual component rendering</li></ul></section> <section class="svelte-an6hzx"><h2 class="svelte-an6hzx">Docker Support</h2> <p class="svelte-an6hzx">Development environment:</p> <div class="code-block svelte-an6hzx"><pre class="svelte-an6hzx"><code class="svelte-an6hzx">docker build -f Dockerfile.dev . -t preveltekit-dev
docker run -p3000:3000 -v./src:/app/src -v./public:/app/public preveltekit-dev</code></pre></div> <p class="svelte-an6hzx">Production build:</p> <div class="code-block svelte-an6hzx"><pre class="svelte-an6hzx"><code class="svelte-an6hzx">docker build . -t preveltekit
docker run -p3000:3000 preveltekit</code></pre></div></section> <section class="svelte-an6hzx"><h2 class="svelte-an6hzx">Project Structure</h2> <ul class="svelte-an6hzx"><li class="svelte-an6hzx"><strong class="svelte-an6hzx">/src</strong>: Application source code</li> <li class="svelte-an6hzx"><strong class="svelte-an6hzx">/public</strong>: Static assets</li> <li class="svelte-an6hzx"><strong class="svelte-an6hzx">/dist</strong>: Production build output</li> <li class="svelte-an6hzx"><strong class="svelte-an6hzx">rsbuild.config.ts</strong>: Build configuration</li> <li class="svelte-an6hzx"><strong class="svelte-an6hzx">ssr.mjs</strong>: SSPR implementation</li></ul></section> <section class="svelte-an6hzx"><h2 class="svelte-an6hzx">Deployment</h2> <p class="svelte-an6hzx">The production build generates static files with pre-compressed variants:</p> <ul class="svelte-an6hzx"><li class="svelte-an6hzx">Standard files (.js, .css, .html)</li> <li class="svelte-an6hzx">Brotli compressed (.br)</li> <li class="svelte-an6hzx">Gzip compressed (.gz)</li> <li class="svelte-an6hzx">Zstandard compressed (.zst)</li></ul> <p class="svelte-an6hzx">Deploy to any static hosting or web server that supports serving compressed assets.</p></section></div>`);
function Documentation($$anchor) {
    "use strict";
    var _window;
    // Demo SSPR capabilities
    let renderType = client/* .mutable_source */.zgK("Client Rendered");
    if ((_window = window) === null || _window === void 0 ? void 0 : _window.JSDOM) {
        client/* .set */.hZp(renderType, "Server Pre-Rendered");
    }
    var div = Documentation_svelte_root();
    var header = client/* .child */.jfp(div);
    var p = client/* .sibling */.hg4(client/* .child */.jfp(header), 2);
    var text = client/* .child */.jfp(p);
    client/* .reset */.cLc(p);
    client/* .reset */.cLc(header);
    var section = client/* .sibling */.hg4(header, 6);
    var div_1 = client/* .sibling */.hg4(client/* .child */.jfp(section), 6);
    var pre = client/* .child */.jfp(div_1);
    var code = client/* .child */.jfp(pre);
    code.textContent = '// Basic detection\nlet renderInfo = "Client Rendered";\nif (window?.JSDOM) {\n    renderInfo = "Server Pre-Rendered";\n}';
    client/* .reset */.cLc(pre);
    client/* .reset */.cLc(div_1);
    var div_2 = client/* .sibling */.hg4(div_1, 6);
    var pre_1 = client/* .child */.jfp(div_2);
    var code_1 = client/* .child */.jfp(pre_1);
    code_1.textContent = '$effect(() => {\n    if (!window?.JSDOM) {\n        fetchBitcoinPrice();\n        // Set up refresh interval\n        const interval = setInterval(fetchBitcoinPrice, 60000);\n        return () => clearInterval(interval);\n    }\n});';
    client/* .reset */.cLc(pre_1);
    client/* .reset */.cLc(div_2);
    client/* .next */.K2T(4);
    client/* .reset */.cLc(section);
    var section_1 = client/* .sibling */.hg4(section, 2);
    var div_3 = client/* .sibling */.hg4(client/* .child */.jfp(section_1), 4);
    var pre_2 = client/* .child */.jfp(div_3);
    var code_2 = client/* .child */.jfp(pre_2);
    code_2.textContent = 'export default defineConfig({\n    environments: {\n        web: {\n            plugins: [\n                pluginSvelte(),\n                pluginCssMinimizer()\n            ],\n            source: {\n                entry: {\n                    // Each entry corresponds to a static route\n                    index: \'./src/index.ts\',    // https://example.com/\n                    doc: \'./src/index.ts\',      // https://example.com/doc\n                    example: \'./src/index.ts\',  // https://example.com/example\n                }\n            },\n            output: {\n                target: \'web\',\n                minify: process.env.NODE_ENV === \'production\',\n            }\n        }\n    }\n});';
    client/* .reset */.cLc(pre_2);
    client/* .reset */.cLc(div_3);
    client/* .next */.K2T(4);
    client/* .reset */.cLc(section_1);
    client/* .next */.K2T(6);
    client/* .reset */.cLc(div);
    client/* .template_effect */.vNg(()=>client/* .set_text */.jax(text, `(${client/* .get */.JtY(renderType) ?? ''})`));
    client/* .append */.BCw($$anchor, div);
}


;// CONCATENATED MODULE: ./src/Example.svelte.3.css!=!./node_modules/.pnpm/svelte-loader@3.2.4_svelte@5.39.3/node_modules/svelte-loader/index.js?cssPath=/home/draft/git/pv2/example/src/Example.svelte.3.css!./src/Example.svelte
// extracted by css-extract-rspack-plugin

;// CONCATENATED MODULE: ./src/Example.svelte


var root_1 = client/* .from_html */.vUu(`<div class="loading svelte-18byggp">Loading Bitcoin prices...</div>`);
var root_3 = client/* .from_html */.vUu(`<div class="error svelte-18byggp"> <button class="svelte-18byggp">Retry</button></div>`);
var root_5 = client/* .from_html */.vUu(`<div class="price-card svelte-18byggp"><div class="price-header svelte-18byggp"><span class="currency-code svelte-18byggp"><!></span> <span class="update-time svelte-18byggp"> </span></div> <div class="current-price svelte-18byggp"><!> </div></div> <div class="disclaimer svelte-18byggp">Cryptocurrency prices are highly volatile and subject to market risks. The displayed price information is for reference only and may not reflect real-time market conditions. Past performance is not indicative of future results. Please conduct your own research and consider your financial situation before making any investment decisions.</div>`, 1);
var Example_svelte_root = client/* .from_html */.vUu(`<div class="bitcoin-dashboard svelte-18byggp"><div class="header svelte-18byggp"><h2 class="svelte-18byggp">Bitcoin Price Tracker</h2> <p class="render-info svelte-18byggp"> </p></div> <div class="price-display svelte-18byggp"><!></div></div>`);
function Example($$anchor, $$props) {
    var _window;
    client/* .push */.VCO($$props, true);
    "use strict";
    let priceData = client/* .state */.wk1(null);
    let loading = client/* .state */.wk1(true);
    let error = client/* .state */.wk1(null);
    // Demo the SSPR capability
    let renderInfo = client/* .state */.wk1("Client Rendered");
    if ((_window = window) === null || _window === void 0 ? void 0 : _window.JSDOM) {
        client/* .set */.hZp(renderInfo, "Server Pre-Rendered");
    }
    async function fetchBitcoinPrice() {
        try {
            client/* .set */.hZp(loading, true);
            client/* .set */.hZp(error, null);
            const response = await fetch('https://min-api.cryptocompare.com/data/generateAvg?fsym=BTC&tsym=USD&e=coinbase');
            if (!response.ok) throw new Error('Failed to fetch data');
            client/* .set */.hZp(priceData, await response.json(), true);
        } catch (e) {
            client/* .set */.hZp(error, e instanceof Error ? e.message : 'An error occurred', true);
        } finally{
            client/* .set */.hZp(loading, false);
        }
    }
    // Fetch initial data
    client/* .user_effect */.MWq(()=>{
        var _window;
        if (!((_window = window) === null || _window === void 0 ? void 0 : _window.JSDOM)) {
            fetchBitcoinPrice();
            // Set up refresh interval
            const interval = setInterval(fetchBitcoinPrice, 60000); // Update every minute
            return ()=>clearInterval(interval);
        }
    });
    var div = Example_svelte_root();
    var div_1 = client/* .child */.jfp(div);
    var p = client/* .sibling */.hg4(client/* .child */.jfp(div_1), 2);
    var text = client/* .child */.jfp(p);
    client/* .reset */.cLc(p);
    client/* .reset */.cLc(div_1);
    var div_2 = client/* .sibling */.hg4(div_1, 2);
    var node = client/* .child */.jfp(div_2);
    {
        var consequent = ($$anchor)=>{
            var div_3 = root_1();
            client/* .append */.BCw($$anchor, div_3);
        };
        var alternate_1 = ($$anchor)=>{
            var fragment = client/* .comment */.Imx();
            var node_1 = client/* .first_child */.esp(fragment);
            {
                var consequent_1 = ($$anchor)=>{
                    var div_4 = root_3();
                    var text_1 = client/* .child */.jfp(div_4);
                    var button = client/* .sibling */.hg4(text_1);
                    button.__click = fetchBitcoinPrice;
                    client/* .reset */.cLc(div_4);
                    client/* .template_effect */.vNg(()=>client/* .set_text */.jax(text_1, `Error: ${client/* .get */.JtY(error) ?? ''} `));
                    client/* .append */.BCw($$anchor, div_4);
                };
                var alternate = ($$anchor)=>{
                    var fragment_1 = client/* .comment */.Imx();
                    var node_2 = client/* .first_child */.esp(fragment_1);
                    {
                        var consequent_2 = ($$anchor)=>{
                            var fragment_2 = root_5();
                            var div_5 = client/* .first_child */.esp(fragment_2);
                            var div_6 = client/* .child */.jfp(div_5);
                            var span = client/* .child */.jfp(div_6);
                            var node_3 = client/* .child */.jfp(span);
                            client/* .html */.qyt(node_3, ()=>client/* .get */.JtY(priceData).RAW.FROMSYMBOL);
                            client/* .reset */.cLc(span);
                            var span_1 = client/* .sibling */.hg4(span, 2);
                            var text_2 = client/* .child */.jfp(span_1);
                            client/* .reset */.cLc(span_1);
                            client/* .reset */.cLc(div_6);
                            var div_7 = client/* .sibling */.hg4(div_6, 2);
                            var node_4 = client/* .child */.jfp(div_7);
                            client/* .html */.qyt(node_4, ()=>client/* .get */.JtY(priceData).RAW.TOSYMBOL);
                            var text_3 = client/* .sibling */.hg4(node_4);
                            client/* .reset */.cLc(div_7);
                            client/* .reset */.cLc(div_5);
                            client/* .next */.K2T(2);
                            client/* .template_effect */.vNg(()=>{
                                client/* .set_text */.jax(text_2, `Last Updated: ${client/* .get */.JtY(priceData).RAW.LASTUPDATE ?? ''}`);
                                client/* .set_text */.jax(text_3, ` ${client/* .get */.JtY(priceData).RAW.PRICE ?? ''}`);
                            });
                            client/* .append */.BCw($$anchor, fragment_2);
                        };
                        client["if"](node_2, ($$render)=>{
                            if (client/* .get */.JtY(priceData)) $$render(consequent_2);
                        }, true);
                    }
                    client/* .append */.BCw($$anchor, fragment_1);
                };
                client["if"](node_1, ($$render)=>{
                    if (client/* .get */.JtY(error)) $$render(consequent_1);
                    else $$render(alternate, false);
                }, true);
            }
            client/* .append */.BCw($$anchor, fragment);
        };
        client["if"](node, ($$render)=>{
            if (client/* .get */.JtY(loading) && !client/* .get */.JtY(priceData)) $$render(consequent);
            else $$render(alternate_1, false);
        });
    }
    client/* .reset */.cLc(div_2);
    client/* .reset */.cLc(div);
    client/* .template_effect */.vNg(()=>client/* .set_text */.jax(text, `(${client/* .get */.JtY(renderInfo) ?? ''})`));
    client/* .append */.BCw($$anchor, div);
    client/* .pop */.uYY();
}
client/* .delegate */.MmH([
    'click'
]);


;// CONCATENATED MODULE: ./src/Index.svelte.0.css!=!./node_modules/.pnpm/svelte-loader@3.2.4_svelte@5.39.3/node_modules/svelte-loader/index.js?cssPath=/home/draft/git/pv2/example/src/Index.svelte.0.css!./src/Index.svelte
// extracted by css-extract-rspack-plugin

;// CONCATENATED MODULE: ./src/Index.svelte







var Index_svelte_root = client/* .from_html */.vUu(`<div class="app svelte-15izthb"><nav class="svelte-15izthb"><div class="nav-content svelte-15izthb"><div class="nav-items svelte-15izthb"><span class="logo svelte-15izthb">PrevelteKit</span> <div class="nav-links svelte-15izthb"><a href="./" class="svelte-15izthb">Home</a> <a href="./doc" class="svelte-15izthb">Documentation</a> <a href="./example" class="svelte-15izthb">Example</a> <a href="https://github.com/tbocek/preveltekit/" target="_blank" rel="noopener noreferrer" aria-label="Link to project on GitHub" class="svelte-15izthb"><svg height="24" viewBox="0 0 16 16" width="24" style="fill: #4299e1;"><path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"></path></svg></a> <div class="server-side-indicator svelte-15izthb">(these links are server-side routed links)</div></div></div></div></nav> <main class="svelte-15izthb"><div class="content svelte-15izthb"><!></div></main> <footer class="svelte-15izthb"><div class="footer-content svelte-15izthb"><p><span class="highlight">PrevelteKit</span> • Lightning fast builds • Pre-rendered SPA • Pure HTML/CSS/JS output</p></div></footer></div>`);
function Index($$anchor) {
    const routes = [
        {
            path: "*/doc",
            component: Documentation,
            static: "doc.html"
        },
        {
            path: "*/example",
            component: Example,
            static: "example.html"
        },
        {
            path: "*/",
            component: Landing,
            static: "index.html"
        }
    ];
    var div = Index_svelte_root();
    var main = client/* .sibling */.hg4(client/* .child */.jfp(div), 2);
    var div_1 = client/* .child */.jfp(main);
    var node = client/* .child */.jfp(div_1);
    (0,Router_svelte/* ["default"] */.A)(node, {
        get routes () {
            return routes;
        }
    });
    client/* .reset */.cLc(div_1);
    client/* .reset */.cLc(main);
    client/* .next */.K2T(2);
    client/* .reset */.cLc(div);
    client/* .append */.BCw($$anchor, div);
}


;// CONCATENATED MODULE: ./src/index.ts


(0,index_client/* .hydrate */.Qv)(Index, {
    target: document.getElementById('root'),
    props: {}
});


}),

});
/************************************************************************/
// The module cache
var __webpack_module_cache__ = {};

// The require function
function __webpack_require__(moduleId) {

// Check if module is in cache
var cachedModule = __webpack_module_cache__[moduleId];
if (cachedModule !== undefined) {
return cachedModule.exports;
}
// Create a new module (and put it into the cache)
var module = (__webpack_module_cache__[moduleId] = {
exports: {}
});
// Execute the module function
__webpack_modules__[moduleId](module, module.exports, __webpack_require__);

// Return the exports of the module
return module.exports;

}

// expose the modules object (__webpack_modules__)
__webpack_require__.m = __webpack_modules__;

/************************************************************************/
// webpack/runtime/define_property_getters
(() => {
__webpack_require__.d = (exports, definition) => {
	for(var key in definition) {
        if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
            Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
        }
    }
};
})();
// webpack/runtime/has_own_property
(() => {
__webpack_require__.o = (obj, prop) => (Object.prototype.hasOwnProperty.call(obj, prop))
})();
// webpack/runtime/on_chunk_loaded
(() => {
var deferred = [];
__webpack_require__.O = (result, chunkIds, fn, priority) => {
	if (chunkIds) {
		priority = priority || 0;
		for (var i = deferred.length; i > 0 && deferred[i - 1][2] > priority; i--)
			deferred[i] = deferred[i - 1];
		deferred[i] = [chunkIds, fn, priority];
		return;
	}
	var notFulfilled = Infinity;
	for (var i = 0; i < deferred.length; i++) {
		var [chunkIds, fn, priority] = deferred[i];
		var fulfilled = true;
		for (var j = 0; j < chunkIds.length; j++) {
			if (
				(priority & (1 === 0) || notFulfilled >= priority) &&
				Object.keys(__webpack_require__.O).every((key) => (__webpack_require__.O[key](chunkIds[j])))
			) {
				chunkIds.splice(j--, 1);
			} else {
				fulfilled = false;
				if (priority < notFulfilled) notFulfilled = priority;
			}
		}
		if (fulfilled) {
			deferred.splice(i--, 1);
			var r = fn();
			if (r !== undefined) result = r;
		}
	}
	return result;
};

})();
// webpack/runtime/jsonp_chunk_loading
(() => {

      // object to store loaded and loading chunks
      // undefined = chunk not loaded, null = chunk preloaded/prefetched
      // [resolve, reject, Promise] = chunk loading, 0 = chunk loaded
      var installedChunks = {"410": 0,};
      __webpack_require__.O.j = (chunkId) => (installedChunks[chunkId] === 0);
// install a JSONP callback for chunk loading
var webpackJsonpCallback = (parentChunkLoadingFunction, data) => {
	var [chunkIds, moreModules, runtime] = data;
	// add "moreModules" to the modules object,
	// then flag all "chunkIds" as loaded and fire callback
	var moduleId, chunkId, i = 0;
	if (chunkIds.some((id) => (installedChunks[id] !== 0))) {
		for (moduleId in moreModules) {
			if (__webpack_require__.o(moreModules, moduleId)) {
				__webpack_require__.m[moduleId] = moreModules[moduleId];
			}
		}
		if (runtime) var result = runtime(__webpack_require__);
	}
	if (parentChunkLoadingFunction) parentChunkLoadingFunction(data);
	for (; i < chunkIds.length; i++) {
		chunkId = chunkIds[i];
		if (
			__webpack_require__.o(installedChunks, chunkId) &&
			installedChunks[chunkId]
		) {
			installedChunks[chunkId][0]();
		}
		installedChunks[chunkId] = 0;
	}
	return __webpack_require__.O(result);
};

var chunkLoadingGlobal = self["webpackChunkpreveltekit_example"] = self["webpackChunkpreveltekit_example"] || [];
chunkLoadingGlobal.forEach(webpackJsonpCallback.bind(null, 0));
chunkLoadingGlobal.push = webpackJsonpCallback.bind(null, chunkLoadingGlobal.push.bind(chunkLoadingGlobal));

})();
/************************************************************************/
// startup
// Load entry module and return exports
// This entry module depends on other loaded chunks and execution need to be delayed
var __webpack_exports__ = __webpack_require__.O(undefined, ["298"], function() { return __webpack_require__(529) });
__webpack_exports__ = __webpack_require__.O(__webpack_exports__);
})()
;