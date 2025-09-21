// src/ssr.ts - SSR functionality (same as previous artifact)
import { JSDOM, ResourceLoader, VirtualConsole } from 'jsdom';
import path from 'node:path';
import fs from 'node:fs';
import express from 'express';
import { createRsbuild, loadConfig } from '@rsbuild/core';
import type { JSDOMInstance } from './types.js';
import { mergeRsbuildConfig } from '@rsbuild/core';
import { defaultConfig } from './default-config.js';

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

async function fakeBrowser(ssrUrl: string, html: string, resourceFolder?: string, timeout = 5000): Promise<JSDOMInstance> {  const virtualConsole = new VirtualConsole();
  virtualConsole.forwardTo(console, { omitJSDOMErrors: true });
  virtualConsole.on("jsdomError", (e) => {
    if (e.type === "not implemented" && e.message.match("navigation")) {
      // handle navigation logic
    } else {
      console.error(e);
    }
  });

  const dom = new JSDOM(html, {
    url: ssrUrl,
    pretendToBeVisual: true,
    runScripts: 'dangerously',
    resources: new LocalResourceLoader(resourceFolder),
    virtualConsole,
  });

  dom.window.JSDOM = true;

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
        clearTimeout(timeoutHandle);
        reject(error);
      }
    }
  });
}

export class PrevelteSSR {
  private config: any;

  async build() {
    const { content } = await loadConfig();
    const finalConfig = mergeRsbuildConfig(defaultConfig, content);
    console.log("AOEUAOEUAOE",finalConfig);
    const rsbuild = await createRsbuild({ rsbuildConfig: finalConfig });
    await rsbuild.build();
    this.config = rsbuild.getRsbuildConfig();
    return this.config;
  }

  async generateSSRHtml() {
    const config = this.config || await this.build();
    const processedDoms = new Map();

    const indexFileName = `${config.output.distPath.root}/index.html`;
    const indexHtml = await fs.promises.readFile(path.join(process.cwd(), indexFileName), "utf-8");
    const indexDom = await fakeBrowser('http://localhost/', indexHtml, config.output.distPath.root);
    processedDoms.set('index.html', indexDom);

    const svelteRoutes = indexDom.window.__svelteRoutes;
    if (svelteRoutes && Array.isArray(svelteRoutes)) {
      const promises = svelteRoutes.map(async route => {
        if (!route?.static || !route?.path) return;
        if (processedDoms.has(route.static)) return;

        const cleanPath = route.path.replace(/\*/g, '').replace(/^\//, '');
        try {
          const dom = await fakeBrowser(`http://localhost/${cleanPath}`, indexHtml, config.output.distPath.root);
          processedDoms.set(route.static, dom);
        } catch (error) {
          console.error(`Error processing route ${cleanPath}:`, error);
        }
      });

      await Promise.all(promises);
    }

    for (const [staticName, dom] of processedDoms.entries()) {
      const fileName = `${config.output.distPath.root}/${staticName}`;
      const finalHtml = dom.serialize();
      await fs.promises.writeFile(fileName, finalHtml);
      console.log(`Generated ${fileName}`);
    }
    
    if (process.env.NODE_ENV === 'production') {
      await this.compressFiles();
    }
  }
  
  private async compressFiles() {
    const { exec } = await import('child_process');
    const { promisify } = await import('util');
    const execAsync = promisify(exec);
    
    const distPath = this.config?.output?.distPath?.root || 'dist';
    
    try {
      await execAsync(
        `find ${distPath} -regex '.*\\.\\(js\\|css\\|html\\)$' -exec sh -c 'zopfli {} & brotli -f {} & zstd -19f {} > /dev/null 2>&1 & wait' \\;`
      );
      console.log('Files compressed with brotli, zopfli, and zstd');
    } catch (error) {
      console.warn('Compression failed:', error);
    }
  }

  createDevServer() {
    return async (port = 3000) => {
      const { content } = await loadConfig();
      const finalConfig = mergeRsbuildConfig(defaultConfig, content);
      console.log("AOEUAOEUAOE2",finalConfig);
      const rsbuild = await createRsbuild({ rsbuildConfig: finalConfig });
      const rsbuildServer = await rsbuild.createDevServer();
      
      const app = express();
      
      app.get('/', async (req, res, next) => {
        if (req.url.includes('/static/') || req.url.includes('/__rsbuild_hmr') || req.url.includes('.hot-update.')) {
          return next();
        }

        try {
          const template = await rsbuildServer.environments.web.getTransformedHtml("index");
          const fullUrl = `${req.protocol}://${req.get('host')}${req.url}`;
          const dom = await fakeBrowser(fullUrl, template);
          res.writeHead(200, {'Content-Type': 'text/html'});
          res.end(dom.serialize());
        } catch (err) {
          console.error('SSR render error, downgrade to CSR...\n', err);
          next();
        }
      });

      app.use(rsbuildServer.middlewares);

      const httpServer = app.listen(port, async () => {
        await rsbuildServer.afterListen();
        console.log(`Dev server running on port ${port}`);
      });

      return { server: httpServer, rsbuildServer };
    };
  }

  createStageServer() {
    return (port = 3000) => {
      const app = express();
      app.use(express.static('dist'));
      app.get('*', (req, res) => {
        const htmlFile = path.join(process.cwd(), 'dist', `${req.url.slice(1) || 'index'}.html`);
        const normalizedPath = path.normalize(htmlFile).replace(/^(\.\.[\/\\])+/, '');
        
        res.sendFile(normalizedPath, err => {
          if (err) {
            res.sendFile(path.join(process.cwd(), 'dist', 'index.html'));
          }
        });
      });
      
      return app.listen(port, () => {
        console.log(`Stage server running on port ${port}`);
      });
    };
  }
}