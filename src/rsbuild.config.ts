import { defineConfig } from "@rsbuild/core";
import { pluginSvelte } from "@rsbuild/plugin-svelte";
import { existsSync, readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";

// When used as a library, with virtual modules, other users that are using this library can pretend to have 
// ./src/index.ts in the ./src directory without having it. If the file exists, the user provided file will
// be used, otherwise this default will be used. Please note, this feature is marked experimental with rspack.
function getVirtualModules() {
  const virtualModules: Record<string, string> = {};

  const currentFile = fileURLToPath(import.meta.url);
  const currentDir = dirname(currentFile);

  if (!existsSync("./src/index.ts")) {
    const libraryIndexPath = join(currentDir, "default", "index.ts");
    const content = readFileSync(libraryIndexPath, "utf8");
    virtualModules["./src/index.ts"] = content;
  }

  return virtualModules;
}

export const defaultConfig = defineConfig({
  environments: {
    web: {
      plugins: [pluginSvelte()],
      source: {
        entry: {
          index: "./src/index.ts", //default provided see above (virtual modules)
        },
        define: {
          '__SSR_BUILD__': JSON.stringify(true),  
        },
      },
      output: {
        target: "web",
      },
    },
  },
  dev: { hmr: false }, //I had issues with hmr in the past, easiest to disable it
  html: { template: "./src/index.html" }, //default provided, see ssr.ts
  output: { assetPrefix: "./" }, //create relative paths, to run in subdirectories
  tools: { //tools only exist here due to virtual modules, needs restart when changed
    rspack: async (config) => {
      const virtualModules = getVirtualModules();

      // Only create VirtualModulesPlugin if we have virtual modules to add
      if (Object.keys(virtualModules).length > 0) {
        const { rspack } = await import("@rspack/core");
        const { VirtualModulesPlugin } = rspack.experiments;

        config.plugins = config.plugins || [];
        config.plugins.push(new VirtualModulesPlugin(virtualModules));
      }

      return config;
    },
  },
});
