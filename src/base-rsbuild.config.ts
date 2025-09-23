import { defineConfig } from "@rsbuild/core";
import { pluginSvelte } from "@rsbuild/plugin-svelte";
import fs from "node:fs";
import { fileURLToPath } from "node:url";
import path from "node:path";

// Helper function to collect virtual modules
function getVirtualModules() {
  const virtualModules: Record<string, string> = {};

  // Re-calculate __dirname
  const currentFile = fileURLToPath(import.meta.url);
  const currentDir = path.dirname(currentFile);

  // Check for index.ts
  if (!fs.existsSync("./src/index.ts")) {
    const libraryIndexPath = path.join(currentDir, "default", "index.ts");
    const content = fs.readFileSync(libraryIndexPath, "utf8");
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
          index: "./src/index.ts",
        },
      },
      output: {
        target: "web",
        minify: process.env.NODE_ENV === "production",
      },
    },
  },
  dev: { hmr: false },
  html: { template: "./src/index.html" },
  output: { assetPrefix: "./" },
  tools: {
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
