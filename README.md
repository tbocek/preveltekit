# Svelte 5 SSR with Rsbuild

A Svelte 5 application template featuring Server-Side Pre Rendering (SSPR) using [Rsbuild](https://rsbuild.dev/) as the build tool.

This project was created to fill a gap in the Svelte ecosystem. While there is a go-to solution for SSR for Svelte (SvelteKit), there isn't a lightweight example showing how to implement SSPR using Rsbuild - a fast build tool.

The inspiration for this project comes from the Vue SSR example in the [Rspack examples repository](https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs). This project adapts those concepts for Svelte, providing a minimal setup for server-side pre-rendering with Svelte and Rsbuild.

## Modern Web Rendering Approaches: SSR vs. SSG vs. SSPR

Web applications can be rendered in several ways, each with distinct characteristics and use cases. Server-Side Rendering (SSR), as implemented in frameworks like SvelteKit, generates HTML dynamically on each request. The server executes the application code, produces HTML with initial state, and sends it to the client along with JavaScript for hydration, enabling interactivity after the page loads.

Static Site Generation (SSG) takes a different approach by generating plain HTML files at build time. These static files are deployed directly to a web server, making them extremely fast to serve. However, SSG typically doesn't include hydration, meaning the pages remain static without client-side interactivity.

Since I did not find a proper term to describe a mix of both, I call it Server-Side Pre Rendering (SSPR). Like SSG, it pre-renders content at build time, but unlike SSG, it includes hydration code. The result is a set of static HTML, JavaScript, and CSS files that can be served by any standard web server (Caddy, Nginx, Apache). This approach provides fast initial page loads like SSG, while maintaining the ability to become fully interactive like SSR.

## Why Rsbuild + Svelte SSR?

- **Lightweight**: No complex framework overhead, just the essentials
- **Flexible**: Full control over your SSR implementation
- **Fast Builds**: Leverages Rsbuild's performance optimizations
- **Modern Stack**: Uses latest versions of Svelte and TypeScript

## Features

- ⚡️ **Svelte 5** - Latest version of the Svelte framework
- 🔥 **TypeScript** - Full type safety and modern JavaScript features
- 📦 **Rsbuild** - Fast and flexible build tool with dual environment support
- 🎯 **Pre-rendered SSR Support** - Pre-rendered Server-side rendering for improved performance and SEO
- 🛠️ **Development Server** - Live reload and fast refresh
- 🎨 **CSS Support** - Built-in CSS processing with PostCSS

## Prerequisites

Make sure you have the following installed:
- Node.js (Latest LTS version recommended)
- pnpm (Recommended package manager)

## Setup

1. Install dependencies:
```bash
pnpm install
```

2. Start the development server:
```bash
pnpm dev
```
This starts an Express development server with:
- Live reloading
- No optimization for faster builds
- Ideal for rapid development

3. Build for production:
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

**Note**: The development server prioritizes fast rebuilds and developer experience, while the production build focuses on optimization and performance. Always test your application with a production build before deploying.

4. Using Docker

To build with docker in production mode, use

```bash
docker build . -t tag
docker run -p3000:3000 tag
```

To run in dev mode, run

```bash
docker build -f Dockerfile.dev . -t tag
docker run -p3000:3000 -v./src:/app/src -v./public:/app/public tag
```