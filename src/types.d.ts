export type RouteParams = Record<string, string>;

export interface Route {
    path: string;
    component: Component;
    static: string;
}

export type Routes = Route[];
export type Component = any;

//export declare function route(node: HTMLElement, href: string): void;

// JSDOM instance interface
export interface JSDOMInstance {
  window: any;
  serialize(): string;
}

// Types in npmjs do not cover the latest version 27.0.0, so make it compile:
declare module 'jsdom' {
  export class VirtualConsole {
    constructor();
    forwardTo(console: any, options?: { omitJSDOMErrors?: boolean }): void;
    on(event: string, callback: (error: any) => void): void;
  }
}

// Extend HTMLScriptElement to include readyState for JSDOM compatibility
declare global {
  // Extend Window interface for JSDOM and Svelte routes
  interface Window {
    __isBuildTime?: boolean;
    __svelteRoutes?: Route[];
  }
}