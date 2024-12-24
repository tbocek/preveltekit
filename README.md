# 🚀 PrevelteKit

PrevelteKit is a lightweight, high-performance web application framework built on Svelte 5, featuring Server-Side Pre Rendering (SSPR) using [Rsbuild](https://rsbuild.dev/) as the build tool.

## 🌟 Why PrevelteKit?
This project was created to fill a gap in the Svelte ecosystem. While there is a go-to solution for SSR for Svelte (SvelteKit), there isn't a lightweight example showing how to implement SSPR using Rsbuild.

The inspiration for this project comes from the Vue SSR example in the [Rspack examples repository](https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs). This project adapts those concepts for Svelte, providing a minimal setup for server-side pre-rendering with Svelte and Rsbuild.

## ✨ Key Features
 * ⚡️ Lightning Fast: Pre-rendered pages with hydration for optimal performance
 * 🎯 Simple Routing: Built-in routing system with pre and post hooks
 * 🔄 SSPR Support: Server-Side Pre-Rendering for better SEO and initial load times
 * 📦 Zero Config: Works out of the box with sensible defaults
 * 🛠️ Developer Friendly: Hot reload in development, production-ready in minutes
 * 🔧 Flexible: Choose between development, staging, and production environments

## Modern Web Rendering Approaches: SSR vs. SSG vs. SSPR

Web applications can be rendered in several ways, each with distinct characteristics and use cases. Server-Side Rendering (SSR), as implemented in frameworks like SvelteKit, generates HTML dynamically on each request. The server executes the application code, produces HTML with initial state, and sends it to the client along with JavaScript for hydration, enabling interactivity after the page loads.

Static Site Generation (SSG) takes a different approach by generating plain HTML files at build time. These static files are deployed directly to a web server, making them extremely fast to serve. However, SSG typically doesn't include hydration, meaning the pages remain static without client-side interactivity.

Since I did not find a proper term to describe a mix of both, I call it Server-Side Pre Rendering (SSPR). Like SSG, it pre-renders content at build time, but unlike SSG, it includes hydration code. The result is a set of static HTML, JavaScript, and CSS files that can be served by any standard web server (Caddy, Nginx, Apache). This approach provides fast initial page loads like SSG, while maintaining the ability to become fully interactive like SSR.

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
    - Gzip (`.gz` files)
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
docker run -p3000:3000 -v./src:/app/src -v./public:/app/public preveltekit-dev
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
