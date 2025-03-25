import {defineConfig} from '@rsbuild/core';
import {pluginSvelte} from '@rsbuild/plugin-svelte';
import * as process from "node:process";
import { convert } from 'tsconfig-to-swcconfig';

const swcConfig = convert("tsconfig.json"); // This will look for tsconfig.json in the current directory
swcConfig.env = null;

export default defineConfig({
    environments: {
        // Configure the web environment for browsers
        web: {
            plugins: [
                pluginSvelte()
            ],
            source: {
                entry: {
                    index: './src/index.ts'
                }
            },
            output: {
                target: 'web',
                minify:  process.env.NODE_ENV === 'production',
            }
        }
    },
    dev: {
        hmr: false
    },
    tools: {
        swc: swcConfig,
    },
    html: {
        template: './static/index.html',
    },
    output: {
        assetPrefix: './'
    },
});
