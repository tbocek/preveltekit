// src/types.d.ts - Complete types for jsdom v27 and express
export type RouteParams = Record<string, string>;

export interface Route {
    path: string;
    component: Component;
    static: string;
}

export type Routes = Route[];
export type Component = any;

// JSDOM instance interface
export interface JSDOMInstance {
  window: any;
  serialize(): string;
}

// Svelte module declarations
declare module '*.svelte' {
  import type { Component } from 'svelte';
  const component: Component;
  export default component;
  export const navigate: (path: string) => void;
  export const route: (node: HTMLAnchorElement) => { destroy: () => void };
}

// Complete jsdom types for v27
declare module 'jsdom' {
  export interface AbortablePromise<T> extends Promise<T> {
    abort(): void;
  }
  
  export interface FetchOptions {
    element?: any;
    [key: string]: any;
  }
  
  export class JSDOM {
    constructor(html?: string, options?: any);
    window: any;
    serialize(): string;
  }
  
  export class ResourceLoader {
    constructor();
    fetch(url: string, options?: FetchOptions): AbortablePromise<Buffer> | null;
  }
  
  export class VirtualConsole {
    constructor();
    forwardTo(console: any, options?: { omitJSDOMErrors?: boolean }): void;
    on(event: string, callback: (error: any) => void): void;
  }
}

// Express types
declare module 'express' {
  interface Request {
    protocol: string;
    url: string;
    get(name: string): string;
    params: any;
  }
  
  interface Response {
    writeHead(statusCode: number, headers: any): void;
    end(data: string): void;
    sendFile(path: string, callback?: (err: any) => void): void;
  }
  
  interface Application {
    get(path: string, handler: (req: Request, res: Response, next?: Function) => void): void;
    use(middleware: any): void;
    listen(port: number, callback?: () => void): any;
  }
  
  function express(): Application;
  namespace express {
    function static(path: string): any;
  }
  
  export = express;
}

// Extend HTMLScriptElement to include readyState for JSDOM compatibility
declare global {
  interface HTMLScriptElement {
    readyState?: 'loading' | 'loaded' | 'complete';
  }
  
  // Extend Error interface for jsdom error types
  interface Error {
    type?: string;
  }

  // Extend Window interface for JSDOM and Svelte routes
  interface Window {
    JSDOM?: boolean;
    __svelteRoutes?: any[];
  }
}