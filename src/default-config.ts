import {defineConfig} from '@rsbuild/core';
import {pluginSvelte} from '@rsbuild/plugin-svelte';

export const defaultConfig = defineConfig({
    environments: {
        web: {
            plugins: [pluginSvelte()],
            source: {
                entry: {
                    index: './src/index.ts'
                }
            },
            output: {
                target: 'web',
                minify: process.env.NODE_ENV === 'production',
            }
        }
    },
    dev: { hmr: false },
    html: { template: './src/index.html' },
    output: { assetPrefix: './' },
});