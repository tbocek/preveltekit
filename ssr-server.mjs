//inspiration: https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs

import express from 'express';
import { createRsbuild, loadConfig } from '@rsbuild/core';

function setupMockLocation(url = 'http://localhost:3000') {
    const parsedUrl = new URL(url);

    global.location = {
        pathname: parsedUrl.pathname,
        search: parsedUrl.search,
        hash: parsedUrl.hash,
        href: parsedUrl.href,
        origin: parsedUrl.origin,
        protocol: parsedUrl.protocol,
        host: parsedUrl.host,
        hostname: parsedUrl.hostname,
        port: parsedUrl.port
    };

    // Mock window
    global.window = {
        location: global.location,
        addEventListener: () => {},
        history: {
            pushState: () => {},
            replaceState: () => {},
            back: () => {},
            forward: () => {},
            go: () => {}
        },
    };

}

// Implement SSR rendering function
const serverRender = (serverAPI) => async (req, res) => {
    // Load SSR bundle
    const indexModule = await serverAPI.environments.ssr.loadBundle('server');

    const fullUrl1 = `http://${req.headers.host}`;
    setupMockLocation(fullUrl1);

    const {head, body} = await indexModule.render();
    const template = await serverAPI.environments.web.getTransformedHtml('index');

    // Insert SSR rendering content into HTML template
    const html = template
        .replace('</head>', `${head}</head>`)
        .replace('<div id="root"></div>', `<div id="root">${body}</div>`);

    res.writeHead(200, {
        'Content-Type': 'text/html',
    });
    res.end(html);
};

// Custom server
async function startDevServer() {
    const { content } = await loadConfig({});

    const rsbuild = await createRsbuild({
        rsbuildConfig: content,
    });

    const app = express();

    const rsbuildServer = await rsbuild.createDevServer();

    const serverRenderMiddleware = serverRender(rsbuildServer);

    // SSR rendering when accessing /index.html
    app.get('*', async (req, res, next) => {

        // Skip SSR for static files and HMR endpoints
        if (req.url.includes('/static/') || req.url.includes('/__rsbuild_hmr')) {
            return next();
        }

        try {
            await serverRenderMiddleware(req, res);
        } catch (err) {
            console.error('SSR render error, downgrade to CSR...\n', err);
            next();
        }
    });

    // Fallback middleware
    app.use((req, res, next) => {
        if (!res.headersSent) {
            rsbuildServer.middlewares(req, res, next);
        }
    });

    const httpServer = app.listen(rsbuildServer.port, async () => {
        await rsbuildServer.afterListen();
    });

    // Connect WebSocket for hot reloading
    rsbuildServer.connectWebSocket({
        server: httpServer,
        // Enable HMR
        hot: true,
        // Enable live reload
        liveReload: true
    });
}

startDevServer();