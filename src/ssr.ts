/// <reference path="./types.d.ts" />
import { JSDOM, ResourceLoader, VirtualConsole } from 'jsdom';
import path from 'node:path';
import fs from 'node:fs';
import { fileURLToPath } from 'node:url';
import express from 'express';
import { createRsbuild, loadConfig } from '@rsbuild/core';
import type { Routes } from './types.js';
import { mergeRsbuildConfig } from '@rsbuild/core';
import { defaultConfig } from './rsbuild.config.js';

class LocalResourceLoader extends ResourceLoader {
  constructor(private resourceFolder?: string) {
    super();
  }

  fetch(url: string, options: any = {}) {
    if (!this.resourceFolder) {
      return super.fetch(url, options);
    }

    const urlPath = new URL(url).pathname;
    const localPath = path.join(process.cwd(), this.resourceFolder, urlPath.replace(/^\//, ''));

    const promise = fs.promises.access(localPath)
      .then(() => fs.promises.readFile(localPath))
      .then(content => Buffer.from(content))
      .catch(() => super.fetch(url, options));

    // Add abort method to match AbortablePromise interface
    (promise as any).abort = () => {};

    return promise as any;
  }
}

class FetchWrapper {
  async fetch(_url: string, _init?: RequestInit): Promise<Response> {
    return new Promise(() => {
      // Promise stays open forever
    });
  }
}

class RequestWrapper {
  constructor(input: RequestInfo | URL, init?: RequestInit) {
    // Resolve relative URLs to absolute URLs using the document's base URL
    if (typeof input === 'string' && !input.startsWith('http://') && !input.startsWith('https://')) {
      const baseURL = globalThis.location?.href || 'http://localhost/';
      input = new URL(input, baseURL).href;
    }
    
    // If there's a signal from JSDOM, remove it before creating Node's Request
    if (init?.signal) {
      const { signal, ...restInit } = init;
      return new Request(input, restInit);
    }
    return new Request(input, init);
  }
}

async function fakeBrowser(ssrUrl: string, html: string, resourceFolder?: string, timeout = 5000): Promise<JSDOM> {
  const virtualConsole = new VirtualConsole();
  // ** for debugging **
  virtualConsole.forwardTo(console);
  virtualConsole.on("jsdomError", (e:any) => {
    if (e.type === "not-implemented" && e.message.includes("navigation")) {
      } else {
      console.error("jsdomError", e);
    }
  });

  const dom = new JSDOM(html, {
    url: ssrUrl,
    pretendToBeVisual: true,
    runScripts: 'dangerously',
    resources: new LocalResourceLoader(resourceFolder),
    virtualConsole,
  });
  
  // Inject fetch and TextEncoder into the JSDOM window
  const fetchWrapper = new FetchWrapper();
  dom.window.fetch = fetchWrapper.fetch.bind(fetchWrapper);
  dom.window.Request = RequestWrapper;
  dom.window.Response = Response;
  dom.window.Headers = Headers;
  dom.window.FormData = FormData;
  dom.window.Blob = Blob;
  dom.window.TextEncoder = TextEncoder;
  dom.window.TextDecoder = TextDecoder;

  dom.window.__isBuildTime = true;

  return new Promise((resolve, reject) => {
    let isResolved = false;
    const timeoutHandle = setTimeout(() => {
      if (!isResolved) {
        isResolved = true;
        reject(new Error('Timeout waiting for resources to load'));
      }
    }, timeout);

    try {
      const allScripts = Array.from(dom.window.document.querySelectorAll('script')) as HTMLScriptElement[];
      let loadedScripts = 0;

      function cleanup() {
        clearTimeout(timeoutHandle);
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

      function handleError(error: any) {
        if (!isResolved) {
          isResolved = true;
          cleanup();
          reject(error);
        }
      }

      allScripts.forEach(script => {
        script.addEventListener('load', handleLoad);
        script.addEventListener('error', handleError);
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
        clearTimeout(timeoutHandle);
        reject(error);
      }
    }
  });
}

export class PrevelteSSR {
  private async createCustomRsbuild() {
    const { content } = await loadConfig();
    const finalConfig = mergeRsbuildConfig(defaultConfig, content);
    const currentDir = path.dirname(fileURLToPath(import.meta.url));
    const libraryDir = path.join(currentDir, 'default');
    if (!fs.existsSync('./src/index.html')) {
      finalConfig.html!.template = path.join(libraryDir, 'index.html');
    }
    const rsbuild = await createRsbuild({ rsbuildConfig: finalConfig });
    return rsbuild;
  }

  private async compressFiles(distPath:string) {
    const { exec } = await import('child_process');
    const { promisify } = await import('util');
    const execAsync = promisify(exec);

    try {
      await execAsync(
        `find ${distPath} -regex '.*\\.\\(js\\|css\\|html\\|svg\\)$' -exec sh -c 'zopfli {} & brotli -f {} & zstd -19f {} > /dev/null 2>&1 & wait' \\;`
      );
      console.log('Files compressed with brotli, zopfli, and zstd');
    } catch (error) {
      console.warn('Compression failed:', error);
    }
  }

  async generateSSRHtml(noZip: boolean) {
    const rsbuild = await this.createCustomRsbuild();
    await rsbuild.build();
    const config = rsbuild.getRsbuildConfig();

    const indexFileName = `${config?.output?.distPath?.root}/index.html`;
    const indexHtml = await fs.promises.readFile(path.join(process.cwd(), indexFileName), "utf-8");
    const dom = await fakeBrowser('http://localhost/', indexHtml, config?.output?.distPath?.root);

    const processedDoms = new Map();
    processedDoms.set('index.html',  dom);

    const svelteRoutes = dom.window.__svelteRoutes as Routes;
    if (svelteRoutes?.staticRoutes) { //we may not have svelteRoutes or staticRoutes
      const promises: Promise<void>[] = [];

      for (const route of svelteRoutes.staticRoutes) {
        if (processedDoms.has(route.htmlFilename)) continue;

        const promise = (async () => {
          try {
            const dom = await fakeBrowser(`http://localhost${route.path}`, indexHtml, config?.output?.distPath?.root);
            processedDoms.set(route.htmlFilename, dom);
          } catch (error) {
            console.error(`Error processing route ${route.path}:`, error);
          }
        })();

        promises.push(promise);
      }

      await Promise.all(promises);
    }

    for (const [htmlFilename, dom] of processedDoms.entries()) {
      const fileName = `${config?.output?.distPath?.root}/${htmlFilename}`;
      const finalHtml = dom.serialize();
      await fs.promises.writeFile(fileName, finalHtml);
      console.log(`Generated ${fileName}`);
      dom.window.close();
    }

    if (!noZip && process.env.NODE_ENV === 'production') {
      const distPath = config?.output?.distPath?.root || 'dist';
      await this.compressFiles(distPath as string);
    }
  }

  async createDevServer() {
    const rsbuild = await this.createCustomRsbuild();
    const rsbuildServer = await rsbuild.createDevServer();
    const template = await rsbuildServer.environments.web.getTransformedHtml("index");
    return async (port = 3000) => {
      const app = express();
      app.use(async (req, res, next) => {
        if (req.url.includes("static/") || req.url.includes("rsbuild-hmr?token=")) {
          return next();
        }
        try {
          const dom = await fakeBrowser(`${req.protocol}://${req.get('host')}${req.url}`, template);
          try {
            const svelteRoutes = dom.window.__svelteRoutes as Routes;
            if (svelteRoutes?.staticRoutes) { //we may not have svelteRoutes or staticRoutes
              for (const route of svelteRoutes.staticRoutes) {
                if (req.url.startsWith(route.path)) {
                  const html = dom.serialize();
                  res.writeHead(200, { 'Content-Type': 'text/html' });
                  res.end(html);
                  return; // stop here, do not continue to next middleware
                }
              }
            }
          } finally {
            dom.window.close();
          }
        } catch (err) {
          console.error(`SSR render error, downgrade to CSR for [${req.url}]`, err);
        }
        return next();
      })
      app.use(rsbuildServer.middlewares);

      const httpServer = app.listen(port, async () => {
        await rsbuildServer.afterListen();
        rsbuildServer.connectWebSocket({ server: httpServer });
        console.log(`Dev server running on port ${port}`);
      });

      return { server: httpServer, rsbuildServer };
    };
  }

  private listHtmlFiles(folder: string): string[] {
    return fs.readdirSync(folder)
      .filter(file => file.endsWith('.html'))
      .map(file => path.basename(file, '.html'));
  }

  createStageServer() {
    return (port = 3000) => {
      const app = express();

      app.get('/', (req, _, next) => {
        req.url = '/index.html';
         return next();
      });

      const htmlFolder = path.join(process.cwd(), 'dist');
      const htmlFiles = this.listHtmlFiles(htmlFolder);
      htmlFiles.forEach((filename) => {
        app.get(`/${filename}`, (req, _, next) => {
          req.url = `/${filename}.html`;
          return next();
        });
      });

      app.use(express.static('dist'));
      return app.listen(port, () => {
        console.log(`Stage server running on port ${port}`);
      });
    };
  }
}
