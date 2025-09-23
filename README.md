# 🚀 PrevelteKit

PrevelteKit is a lightweight, high-performance web application framework built on [Svelte 5](https://svelte.dev/), featuring Server-Side Pre Rendering (SSPR) using [Rsbuild](https://rsbuild.dev/) as the build tool and [jsdom](https://github.com/jsdom/jsdom) as the DOM environment for rendering components on the server side.

## 🌟 Why PrevelteKit?
This project was created to fill a gap in the Svelte ecosystem. While there is a go-to solution for SSR for Svelte (SvelteKit), there isn't a lightweight solution to implement SSPR using Rsbuild. There is the prerender option in SvelteKit, but it's part of SvelteKit that comes with many additional features that might not be necessary for every project. This project can be seen as the minimal setup for server-side pre-rendering and is essentially a glue of  [Svelte](https://svelte.dev/), [Rsbuild](https://rsbuild.dev/), [jsdom](https://github.com/jsdom/jsdom), and [svelte5-router](https://github.com/mateothegreat/svelte5-router), with ~400 lines of code. It's a great starting point for projects that need server-side rendering without the overhead of SvelteKit.

The inspiration for this project comes from the Vue SSR example in the [Rspack examples repository](https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs). This project adapts those concepts for Svelte, providing a minimal setup for server-side pre-rendering with [Svelte](https://svelte.dev/), [Rsbuild](https://rsbuild.dev/), [jsdom](https://github.com/jsdom/jsdom), and [svelte5-router](https://github.com/mateothegreat/svelte5-router).

## ✨ Key Features
 * ⚡️ Lightning Fast: Pre-rendered pages with hydration for optimal performance
 * 🎯 Simple Routing: Built-in routing system with pre and post hooks
 * 🔄 SSPR Support: Server-Side Pre-Rendering for better SEO and initial load times
 * 📦 Zero Config: Works out of the box with sensible defaults
 * 🛠️ Developer Friendly: Hot reload in development, production-ready in minutes
 * 🔧 Flexible: Choose between development, staging, and production environments
 * 🛡️ Security: Docker-based development environments to protect against supply chain attacks

## Modern Web Rendering Approaches: SSR vs. SSG vs. SSPR

Web applications can be rendered in several ways, each with distinct characteristics and use cases. Server-Side Rendering (SSR), as implemented in frameworks like SvelteKit, generates HTML dynamically on each request. The server executes the application code, produces HTML with initial state, and sends it to the client along with JavaScript for hydration, enabling interactivity after the page loads.

Static Site Generation (SSG) takes a different approach by generating plain HTML files at build time. These static files are deployed directly to a web server, making them extremely fast to serve. However, SSG typically doesn't include hydration, meaning the pages remain static without client-side interactivity.

Since I did not find a proper term to describe a mix of both, I call it Server-Side Pre Rendering (SSPR). Like SSG, it pre-renders content at build time, but unlike SSG, it includes hydration code. The result is a set of static HTML, JavaScript, and CSS files that can be served by any standard web server (Caddy, Nginx, Apache). This approach provides fast initial page loads like SSG, while maintaining the ability to become fully interactive like SSR. So, you will see the data provided via REST a bit later, but compared to SPA, you can show the initial layout and structure already. A compromise of simplicity and user experience.

## What about Isomorphic Rendering or Dynamic Rendering?

Isomorphic rendering (IR) is a technique where the same JavaScript code runs on both the server and client sides. When a user requests a page, the server performs the initial render for fast page load, after which the client side takes over to handle interactivity. This approach uses frameworks like React.js with Node.js to maintain a single codebase that works across environments. Companies like Airbnb have implemented this using frameworks such as Rendr, achieving both optimal performance and maintainability.

Dynamic rendering (DR) takes a different approach by serving different versions of content based on who's requesting it. When a search engine bot visits the site, it receives pre-rendered static HTML, while human users get the normal client-side rendered version. This method is particularly valuable for JavaScript-heavy websites that need to maintain good SEO. It's generally simpler to implement than isomorphic rendering and is actively recommended by major search engines like Google and Bing. [source](https://prerender.io/blog/isomorphic-rendering/)

The fundamental difference between these approaches (IR, DR, SSPR) lies in when and how the rendering occurs. Isomorphic rendering maintains consistent code execution across server and client, dynamic rendering serves different versions based on the user agent, and SSPR pre-generates static content with hydration capabilities at build time. Each approach has its own sweet spot depending on the specific needs of the project, balancing factors like performance, SEO requirements, and development complexity.

## Why Rsbuild + Svelte SSR?

- **Lightweight**: No complex framework overhead, just the essentials
- **Flexible**: Full control over your SSR implementation
- **Fast Builds**: Leverages Rsbuild's performance optimizations
- **Modern Stack**: Uses latest versions of Svelte and TypeScript

## Prerequisites

Make sure you have the following installed:
- Node.js (Latest LTS version recommended)
- pnpm (Recommended package manager)

## 🚦 Quick Start

### Install
```bash
git clone https://github.com/tbocek/preveltekit.git
cd preveltekit
pnpm install
```

### Start the development server
```bash
pnpm dev
```
This starts an Express development server with:
- Live reloading
- No optimization for faster builds
- Ideal for rapid development

### Build for production
```bash
pnpm build
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
pnpm stage
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
