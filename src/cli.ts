#!/usr/bin/env node
import { PrevelteSSR } from './ssr.js';
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { dirname, join } from 'node:path';

const args = process.argv.slice(2);
const command = args[0];
const ssr = new PrevelteSSR();

// Get version from package.json
function getVersion():string {
  try {
    const __filename = fileURLToPath(import.meta.url);
    const __dirname = dirname(__filename);
    const packageJsonPath = join(__dirname, '..', 'package.json');
    const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf8'));
    return packageJson.version;
  } catch (error) {
    return '1.0.0'; // fallback
  }
}

const version = getVersion();

function getPort():number {
  const portIndex = args.indexOf('-p') !== -1 ? args.indexOf('-p') : args.indexOf('--port');
  return portIndex !== -1 && args[portIndex + 1] ? parseInt(args[portIndex + 1]) : 3000;
}

function getNoZip(): boolean {
  const noZipIndex = args.indexOf('--no-zip');
  return noZipIndex !== -1;
}

function showHelp() {
  console.log(`
PrevelteKit SSR utilities v${version}

Usage:
  preveltekit <command> [options]

Commands:
  prod                    Run production build and generate SSR HTML
  dev                     Run development server
  stage                   Run staging server

Options:
  -p, --port <port>       Port number (default: 3000)
  --no-zip                Do not compress zip/br/zstd the output files
  -h, --help              Show help
  -v, --version           Show version

Examples:
  preveltekit prod
  preveltekit dev -p 3000
  preveltekit stage --port 8080
`);
}

async function main() {
  try {
    if (!command || args.includes('-h') || args.includes('--help')) {
      showHelp();
      process.exit(0);
    }
    
    if (args.includes('-v') || args.includes('--version')) {
      console.log(version);
      process.exit(0);
    }

    switch (command) {
      case 'prod':
        await ssr.generateSSRHtml(getNoZip());
        break;
        
      case 'dev':
        const devPort = getPort();
        const createDevServer = await ssr.createDevServer();
        await createDevServer(devPort);
        break;
        
      case 'stage':
        const stagePort = getPort();
        await ssr.generateSSRHtml(getNoZip());
        const createStageServer = ssr.createStageServer();
        createStageServer(stagePort);
        break;
        
      default:
        console.error(`Unknown command: ${command}`);
        showHelp();
        process.exit(1);
    }
  } catch (error) {
    console.error(error);
    process.exit(1);
  }
}

main();