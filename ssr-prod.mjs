//inspiration: https://github.com/rspack-contrib/rspack-examples/blob/main/rsbuild/ssr-express/prod-server.mjs

import { Command } from 'commander';
import path from 'path';
import fs from "node:fs";
import { createRequire } from "node:module";

const program = new Command();
const require = createRequire(import.meta.url);

async function generateSSRHtml({ server, client, outputHtml }) {
    try {
        const remotesPath = path.join(process.cwd(), server);
        const importedApp = require(remotesPath);
        const { head, body } = await importedApp.render();
        const template = fs.readFileSync(path.join(process.cwd(), client), "utf-8");

        // Generate final HTML
        const html = template
            .replace('</head>', `${head}</head>`)
            .replace('<div id="root"></div>', `<div id="root">${body}</div>`);

        // Write to output file if specified
        if (outputHtml) {
            const outputPath = path.resolve(process.cwd(), outputHtml);
            fs.writeFileSync(outputPath, html);
            fs.unlinkSync(server);
            if(outputHtml !== client) {
                fs.unlinkSync(client);
            }
            console.log(`HTML generated successfully at: ${outputPath}`);
        } else {
            console.log(html);
        }

    } catch (error) {
        console.error('Error generating SSR HTML:', error);
        throw error;
    }
}

program
    .name('ssr-generator')
    .description('Generate SSR HTML from a Rsbuild entry')
    .version('1.0.0');

program
    .option('-s, --server <name>', 'Name to the Rsbuild server entry file')
    .option('-c, --client <name>', 'Name to the Rsbuild client entry file')
    .option('-o, --output-html <path>', 'Output file path for the generated HTML')
    .option('-d, --output-dir <path>', 'Output file path for the static assets')
    .action(async (options) => {
        try {
            await generateSSRHtml(options);
            process.exit(0);
        } catch (error) {
            process.exit(1);
        }
    });

program.parse();