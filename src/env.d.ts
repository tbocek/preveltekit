/// <reference types="@rsbuild/core/types" />
/// <reference types="svelte" />

declare global {
    interface Window {
        //this is set to true when rendering the page on the server side with jsdom
        JSDOM: boolean;

        // Routes data for SSR route management
        __svelteRoutes?: any;
    }
}

export {};