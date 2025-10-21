#!/usr/bin/env node
import { preprocess } from 'svelte/compiler';
import { sveltePreprocess } from 'svelte-preprocess';
import { readFileSync, writeFileSync } from 'fs';
import { resolve } from 'path';

// Parse command line arguments
const args: string[] = process.argv.slice(2);

const getArg = (flag: string): string | null => {
    const index = args.indexOf(flag);
    return index !== -1 && args[index + 1] ? args[index + 1] : null;
};

const inputFile: string | null = getArg('--input') || getArg('-i');
const outputFile: string | null = getArg('--output') || getArg('-o');
const showHelp: boolean = args.includes('--help') || args.includes('-h');

// Show help message
if (showHelp || !inputFile || !outputFile) {
    console.log(`
Usage: node compileRouter.js --input <input-file> --output <output-file>

Options:
  -i, --input <file>    Input TypeScript Svelte file
  -o, --output <file>   Output JavaScript Svelte file
  -h, --help           Show this help message

Example:
  node compileRouter.js --input src/Router.svelte --output dist/Router.svelte
  node compileRouter.js -i src/Router.svelte -o dist/Router.svelte
`);
    process.exit(showHelp ? 0 : 1);
}

try {
    // Read input file
    const inputPath: string = resolve(inputFile);
    console.log(`Reading: ${inputPath}`);
    const source: string = readFileSync(inputPath, 'utf8');

    // Preprocess TypeScript to JavaScript
    console.log('Compiling TypeScript...');
    const preprocessed = await preprocess(
        source,
        sveltePreprocess({
            typescript: {
                tsconfigFile: './tsconfig.json'
            }
        }),
        { filename: inputFile }
    );

    // Remove lang="ts" from script tag
    let output: string = preprocessed.code.replace(/<script lang="ts">/g, '<script>');

    // Write output file
    const outputPath: string = resolve(outputFile);
    writeFileSync(outputPath, output);
    console.log(`✓ TypeScript compiled successfully`);
    console.log(`Output: ${outputPath}`);
} catch (error: any) {
    console.error(`✗ Error: ${error.message}`);
    process.exit(1);
}