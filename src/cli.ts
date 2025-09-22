#!/usr/bin/env node
import { Command } from 'commander';
import { PrevelteSSR } from './index.js';

const program = new Command();
const ssr = new PrevelteSSR();

program
  .name('preveltekit')
  .description('PrevelteKit SSR utilities')
  .version('1.0.0');

program
  .option('-p, --prod-build', 'Run production build and generate SSR HTML')
  .option('-d, --dev-server', 'Run development server')
  .option('-s, --stage-server', 'Run staging server')
  .option('--port <port>', 'Port number', '3000')
  .action(async (options) => {
    try {
      if (options.prodBuild) {
        process.env.NODE_ENV = 'production';
        await ssr.generateSSRHtml();
        process.exit(0);
      } else if (options.devServer) {
        process.env.NODE_ENV = 'development';
        const createServer = ssr.createDevServer();
        await createServer(parseInt(options.port));
      } else if (options.stageServer) {
        process.env.NODE_ENV = 'production';
        const createServer = ssr.createStageServer();
        createServer(parseInt(options.port));
      } else {
        program.help();
      }
    } catch (error) {
      console.error(error);
      process.exit(1);
    }
  });

program.parse();