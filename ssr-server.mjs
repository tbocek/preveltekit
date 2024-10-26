//inspiration: https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs

import express from 'express';
import { createRsbuild, loadConfig } from '@rsbuild/core';

// Implement SSR rendering function
const serverRender = (serverAPI) => async (_req, res) => {
    // Load SSR bundle
    const indexModule = await serverAPI.environments.ssr.loadBundle('server');
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
    app.get('/', async (req, res, next) => {
        try {
            await serverRenderMiddleware(req, res);
        } catch (err) {
            console.error('SSR render error, downgrade to CSR...\n', err);
            next();
        }
    });

    app.use(rsbuildServer.middlewares);

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