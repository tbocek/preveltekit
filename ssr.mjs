//inspiration: https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs
import {Command} from 'commander';
import path from 'node:path';
import fs from "node:fs";
import express from 'express';
import {JSDOM, ResourceLoader, VirtualConsole} from 'jsdom';
import {createRsbuild, loadConfig} from '@rsbuild/core';

const program = new Command();

//https://github.com/jsdom/jsdom/issues/2112
//TB: I want to ignore reloading with jsdom, but this seems not possible
//so we ignore the error
const virtualConsole = new VirtualConsole();
virtualConsole.sendTo(console, { omitJSDOMErrors: true });
virtualConsole.on("jsdomError", (e) => {
    if (e.type === "not implemented" && e.message.match("navigation")) {
        // handle navigation logic
    } else {
        console.error(e);
    }
});

const {content} = await loadConfig();
const rsbuild = await createRsbuild({rsbuildConfig: content});
const config = rsbuild.getRsbuildConfig();

//Load from disk instead as in the prod release, we do not run a server
class LocalResourceLoader extends ResourceLoader {
    constructor(resourceFolder) {
        super();
        this.resourceFolder = resourceFolder;
    }

    async fetch(url, options) {
        //we are using dev server, so we have nothing on the disk
        if(!this.resourceFolder) {
            const origResource = super.fetch(url, options)
            return Promise.resolve(origResource);
        }

        const urlPath = new URL(url).pathname;
        // Remove leading slash and join with dist directory
        const localPath = path.join(process.cwd(), this.resourceFolder, urlPath.replace(/^\//, ''));

        try {
            await fs.promises.access(localPath);
            const content = await fs.promises.readFile(localPath);
            return Promise.resolve(Buffer.from(content));
        } catch {
            const origResource = super.fetch(url, options)
            return Promise.resolve(origResource);
        }
    }
}

async function fakeBrowser(ssrUrl, html, resourceFolder) {
    const dom = new JSDOM(html, {
        url: ssrUrl,
        pretendToBeVisual: true,
        runScripts: 'dangerously',
        resources: new LocalResourceLoader(resourceFolder),
        virtualConsole
    });

    return new Promise((resolve, reject) => {
        let isResolved = false;
        const timeout = setTimeout(() => {
            if (!isResolved) {
                isResolved = true;
                reject(new Error('Timeout waiting for resources to load'));
            }
        }, 5000);

        try {
            const allScripts = Array.from(dom.window.document.querySelectorAll('script'));
            let loadedScripts = 0;

            function cleanup() {
                clearTimeout(timeout);
            }

            function handleLoadComplete() {
                if (loadedScripts === allScripts.length) {
                    const marker = 'SCRIPTS_EXECUTED_' + Date.now();
                    const markComplete = dom.window.document.createElement('script');
                    markComplete.setAttribute('data-marker', 'true');

                    markComplete.textContent = `
                    Promise.resolve().then(() => {
                        return new Promise(resolve => setTimeout(resolve, 0));
                    }).then(() => {
                        window['${marker}'] = true;
                    });
                `;

                    dom.window.document.body.appendChild(markComplete);

                    let checkCount = 0;
                    const maxChecks = 500;

                    const checkExecution = () => {
                        if (dom.window[marker]) {
                            if (!isResolved) {
                                isResolved = true;
                                cleanup();
                                const markerScript = dom.window.document.querySelector('script[data-marker="true"]');
                                if (markerScript) {
                                    markerScript.remove();
                                }
                                resolve(dom);
                            }
                        } else if (checkCount++ < maxChecks) {
                            setTimeout(checkExecution, 10);
                        } else {
                            if (!isResolved) {
                                isResolved = true;
                                cleanup();
                                reject(new Error('Script execution check timed out'));
                            }
                        }
                    };

                    checkExecution();
                }
            }

            function handleLoad() {
                loadedScripts++;
                handleLoadComplete();
            }

            function handleError(error) {
                if (!isResolved) {
                    isResolved = true;
                    cleanup();
                    reject(error);
                }
            }

            allScripts.forEach(script => {
                if (script.readyState === 'complete' || script.readyState === 'loaded') {
                    handleLoad();
                } else {
                    script.addEventListener('load', handleLoad);
                    script.addEventListener('error', handleError);
                }
            });

            if (allScripts.length === 0) {
                dom.window.addEventListener('load', () => {
                    if (!isResolved) {
                        isResolved = true;
                        cleanup();
                        resolve(dom);
                    }
                });
            }

        } catch (error) {
            if (!isResolved) {
                isResolved = true;
                clearTimeout(timeout);
                reject(error);
            }
        }
    });
}

async function generateSSRHtml() {
    const config = await runRsbuildBuild();

    const startTime = process.hrtime.bigint();
    console.log('Starting SSR HTML generation...');
    try {
        const promises = Object.keys(config.environments.web.source.entry).map(async entryName => {
            const fullUrl = `http://localhost/`; //here it does not matter, as we can get everything from disk
            const fileName = `${config.output.distPath.root}/${entryName}.html`;
            const html = await fs.promises.readFile(path.join(process.cwd(), fileName), "utf-8");
            const dom = await fakeBrowser(fullUrl, html, config.output.distPath.root);
            const finalHtml = dom.serialize();
            await fs.promises.writeFile(fileName, finalHtml);
        });

        try {
            await Promise.all(promises);
        } catch (error) {
            console.error('Error processing entries:', error);
        }

    } catch (error) {
        console.error('Error generating SSR HTML:', error);
        throw error;
    } finally {
        const endTime = process.hrtime.bigint();
        console.log(`Total SSR execution time: ${(endTime - startTime) / BigInt(1000000)}ms`);
    }
}

const mapUrlToEntry = (url, entrySourceMap, basePath) => {
    const urlParts = url === '/' ? [] : url.split('/').filter(Boolean);

    while (urlParts.length >= 0) {
        const currentPath = '/' + (urlParts.length ? urlParts.join('/') : '');

        const entryName = urlParts.length ? `${urlParts.join('_')}_index` : 'index';
        const value = entrySourceMap.get(entryName);
        if (value !== undefined) {
            return { entry: entryName, value, currentPath };
        }

        const sourcePathSlash = basePath + urlParts.length ? `${urlParts.join('/')}/index.ts` : 'index.ts';
        const sourcePathUnderscore = basePath + urlParts.length ? `${urlParts.join('_')}_index.ts` : 'index.ts';

        for (const [entry, source] of entrySourceMap.entries()) {
            if (source === sourcePathSlash || source === sourcePathUnderscore) {
                return { entry, value: source, currentPath };
            }
        }

        if (urlParts.length === 0) break;
        urlParts.pop();
    }

    throw new Error(`No valid entry found for URL: ${url}`);
};

// Implement SSR rendering function
const serverRender = (serverAPI) => async (req, res) => {
    const entrySourceMap = new Map();
    Object.entries(config.environments.web.source.entry).forEach(([entry, source]) => {
        entrySourceMap.set(entry, source);
    });
    const entry = await mapUrlToEntry(req.url, entrySourceMap, './src');

    const template = await serverAPI.environments.web.getTransformedHtml(entry.entry);

    const dom = await fakeBrowser(`${req.protocol}://${req.get('host')}/${entry.currentPath}`, template, null);
    const finalHtml = dom.serialize();
    res.writeHead(200, {
        'Content-Type': 'text/html',
    });
    res.end(finalHtml);
};

async function runRsbuildBuild() {
    try {
        // Load the Rsbuild configuration
        const { content } = await loadConfig();
        const rsbuild = await createRsbuild({
            rsbuildConfig: content
        });

        // Start the build process
        await rsbuild.build();
        return rsbuild.getRsbuildConfig();
    } catch (error) {
        console.error('Build failed:', error);
        throw error;
    }
}

async function devServer() {
    const { content } = await loadConfig();
    const rsbuild = await createRsbuild({rsbuildConfig: content});
    const rsbuildServer = await rsbuild.createDevServer();
    const serverRenderMiddleware = serverRender(rsbuildServer);

    // SSR rendering when accessing /index.html
    app.get('*', async (req, res, next) => {
        // Skip SSR for static files and HMR endpoints
        if (req.url.includes('/static/') || req.url.includes('/__rsbuild_hmr') || req.url.includes('.hot-update.')) {
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

    const connections = new Set();
    httpServer.on('connection', (socket) => {
        connections.add(socket);
        socket.on('close', () => {
            connections.delete(socket);
        });
    });

    // Connect WebSocket for hot reloading
    rsbuildServer.connectWebSocket({
        server: httpServer,
        // Enable HMR
        hot: false,
        // Enable live reload
        liveReload: true
    });

    // Handle graceful shutdown
    function shutdown() {
        console.log('\nShutting down dev server...');
        for (const socket of connections) {
            socket.destroy();
        }
        httpServer.close(async () => {
            await rsbuildServer.close();
            console.log('Rsbuild server closed');
            process.exit(0);
        });

        // Force close after timeout
        setTimeout(() => {
            console.log('Forcing shutdown after timeout');
            process.exit(1);
        }, 5000);
    }

    // Handle different shutdown signals
    process.on('SIGTERM', shutdown);
    process.on('SIGINT', shutdown);
    process.on('SIGHUP', shutdown);

    return httpServer;
}

const app = express();
program
    .name('ssr-generator')
    .description('Generate SSR HTML from a Rsbuild entry')
    .version('1.0.0');

program
    .option('-p, --prod-build', 'Run the prod build, then exit')
    .option('-d, --dev-server', 'Run the dev server')
    .action(async (options) => {
        try {
            if (options.prodBuild) {
                await generateSSRHtml();
                process.exit(0);
            } else if (options.devServer) {
                await devServer();
            } else {
                console.log('Please specify either --prod-build or --dev-server');
                program.help();
                process.exit(1);
            }
        } catch (error) {
            process.exit(1);
        }
    });

program.parse();