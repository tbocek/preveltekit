# Svelte 5 SSR with Rsbuild

A Svelte 5 application template featuring Server-Side Rendering (SSR) using [Rsbuild](https://rsbuild.dev/) as the build tool.

This project was created to fill a gap in the Svelte ecosystem. While there is a go-to solution for SSR for Svelte (SvelteKit), there isn't a lightweight example showing how to implement SSR using Rsbuild - a fast build tool.

The inspiration for this project comes from the Vue SSR example in the [Rspack examples repository](https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs). This project adapts those concepts for Svelte, providing a minimal yet production-ready setup for server-side rendering with Svelte and Rsbuild.

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
