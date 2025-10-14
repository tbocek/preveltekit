import type { Component as SvelteComponent } from 'svelte';
export type Component = SvelteComponent<any> | null;

export interface Routes {
  dynamicRoutes?: { path: string; component: Component }[];
  staticRoutes?: { path: string; htmlFilename: string }[];
}

// Extend HTMLScriptElement to include readyState for JSDOM compatibility
declare global {
  // Extend Window interface for JSDOM and Svelte routes
  interface Window {
    __isBuildTime?: boolean;
    __svelteRoutes?: Routes;
  }
}
