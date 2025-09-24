# PrevelteKit

PrevelteKit is a minimalistic (>500 LoC) web application framework built on [Svelte 5](https://svelte.dev/), featuring Server-Side Pre Rendering (SSPR) using [Rsbuild](https://rsbuild.dev/) as the build tool and [jsdom](https://github.com/jsdom/jsdom) as the DOM environment for rendering components on the server side.

## Why PrevelteKit?
While there is a go-to solution for SSR for Svelte (SvelteKit), I was missing a minimalistic solution just for pre-rendering. There is the prerender option in SvelteKit, but it's part of SvelteKit that comes with many additional features that might not be necessary for every project. This project can be seen as the minimal setup for server-side pre-rendering and is essentially a glue of  [Svelte](https://svelte.dev/), [Rsbuild](https://rsbuild.dev/), and [jsdom](https://github.com/jsdom/jsdom), with 469 lines of code. It's a starting point for projects that need server-side rendering without the overhead of SvelteKit.

The inspiration for this project comes from the Vue SSR example in the [Rspack examples repository](https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs). This project adapts those concepts for Svelte, providing a minimal setup.

## Key Features
 * ⚡️ Lightning Fast: Pre-rendered pages with hydration with Svelte/Rsbuild
 * 🎯 Simple Routing: Built-in routing system
 * 🔄 SSPR Support: Server-Side Pre-Rendering for faster perceived loading
 * 📦 Zero Config: Works out of the box with sensible defaults
 * 🛠️ Developer Friendly: Hot reload in development, production-ready in minutes
 * 🛡️ Security: Docker-based development environments to protect against supply chain attacks

## Modern Web Rendering Approaches: SSR vs. SSG vs. SSPR
Web applications can be rendered in several ways, each with distinct characteristics and use cases. Server-Side Rendering (SSR) generates HTML dynamically on each request. The server executes the application code, produces HTML with initial state, and sends it to the client.

Static Site Generation (SSG) takes a different approach by generating plain HTML files at build time. These static files are deployed directly to a web server, making them extremely fast to serve. However, SSG remain static without client-side interactivity.

SPA on the other hand, loads a single HTML page initially and dynamically updates content using JavaScript. The client downloads the entire application bundle upfront, then handles routing and rendering in the browser. 

Modern meta-frameworks like Next.js, Nuxt.js, SvelteKit, and Remix enable developers to combine these approaches within a single application. They support hybrid rendering strategies such as SSR with client-side hydration, incremental static regeneration (ISR), and page-level rendering choices where different routes can use SSG, SSR, or SPA based on specific requirements.

Since I did not find a proper term to describe a mix of SPA and build-time pre-rendering, lets call it SPAwBR (Single Page Application with Built-time Pre-rendering). Like SSG, it pre-renders content at build time, but unlike SSG, it includes hydration code. The result is a set of static HTML, JavaScript, and CSS files that can be served by any standard web server (Caddy, Nginx, Apache). This approach provides fast initial page loads like SSG with static content such as layout, while maintaining the ability to become fully interactive like SPA. So, you will see the data provided via REST a bit later, but compared to SPA, you can show the initial layout and structure already. A compromise of simplicity and user experience.

So, why not use Next.js, Nuxt.js, SvelteKit? From an architectural point of view, I prefer the clear separation between view code and server code, where the frontend requests data from the backend via dedicated /api endpoints. This approach treats the frontend as purely static assets (HTML/CSS/JS) that can be served from any CDN or simple web server.

Meta-frameworks blur this separation by requiring a JavaScript runtime (Node.js, Deno, or Bun) to handle server-side rendering, API routes, and build-time generation. While platforms like Vercel and Netlify can help with handling this complex setup (they are great services that I used in the past), serving just static content is even simpler. Static-first architecture offers compelling simplicity: deploy anywhere (GitHub Pages, S3, any web server), predictable performance, easier debugging with clear boundaries, and freedom to choose your backend technology. You avoid the "full-stack JavaScript" complexity for your deployed frontend - it's just files on a server, nothing more. No runtime dependencies, no server processes to monitor.

## Prerequisites

Make sure you have the following installed:
- Node.js (Latest LTS version recommended)
- npm/pnpm or similar

## Quick Start

### Install
```bash
# Create test directory and go into this directory
mkdir -p preveltekit/src && cd preveltekit 
# Declare dependency and the dev script
echo '{"dependencies": {"preveltekit":"^1.0.4"}, "scripts": {"dev": "preveltekit dev"}}' > package.json 
# Download dependencies
npm install 
# A very simple svelte file
echo '<script>let count = $state(0);</script><h1>Count: {count}</h1><button onclick={() => count++}>Click me</button>' > src/Index.svelte 
# And open a browser with localhost:3000
npm run dev 
```

## Slow Start

Another example is the [notary example](https://github.com/tbocek/notary-example). Here you can see, which scripts are supported: dev/stage/prod. 
Lets, look at those in the example folder:

### Start the development server
```bash
npm run dev
```
This starts an Express development server with:
- Live reloading
- No optimization for faster builds
- Ideal for rapid development

### Build for production
```bash
npm run build
```
The production build:
- Uses Caddy as the web server
- Generates pre-compressed static files for optimal serving:
    - Brotli (`.br` files)
    - Zstandard (`.zst` files)
    - Zopfli (`.gz` files)
- Optimizes assets for production

### Staging Environment
```bash
npm stage
```

**Note**: The development server prioritizes fast rebuilds and developer experience, while the production build focuses on optimization and performance. Always test your application with a production build before deploying.

## 🐳 Docker Support

To build with docker in production mode, use

```bash
docker build . -t preveltekit
docker run -p3000:3000 preveltekit
```

To run in development mode, run

```bash
docker build -f Dockerfile.dev . -t preveltekit-dev
docker run -p3000:3000 -v./static:/app/static -v./src:/app/src -v./public:/app/public preveltekit-dev
```

## 📚 Technical Details
PrevelteKit uses SSPR (Server-Side Pre-Rendering) to generate static HTML at build time while maintaining full interactivity through hydration. This approach offers:

 * Better SEO: Search engines see fully rendered content
 * Faster Initial Load: Users see content immediately
 * Full Interactivity: Components hydrate seamlessly
 * Simple Deployment: Deploy to any static hosting

## 🔧 Configuration
PrevelteKit is configured through rsbuild.config.ts and supports multiple deployment targets:

 * Development: Hot reload enabled, unminified for debugging
 * Staging: Production build with local server
 * Production: Optimized build with Caddy server
