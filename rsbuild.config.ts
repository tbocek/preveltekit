import {defineConfig} from '@rsbuild/core';
import {pluginSvelte} from '@rsbuild/plugin-svelte';

export default defineConfig({
    environments: {
        // Configure the web environment for browsers
        web: {
            plugins: [
                pluginSvelte()
            ],
            source: {
                entry: {
                    index: './src/entry-client.ts', //creates a html
                },
            },
            output: {
                //assetPrefix: './',
                target: 'web',
            }
        },
        // Configure the node environment for SSR
        ssr: {
            plugins: [
                pluginSvelte({
                  svelteLoaderOptions: {
                    compilerOptions: {
                      //@ts-ignore -> this is the right option
                      generate: 'server'
                    }
                  }
                })
            ],
            source: {
                entry: {
                    server: './src/entry-server.ts', //creates a js
                }
            },
            output: {
                // Use 'node' target for the Node.js outputs
                target: 'node',
            }
        }
    },
});