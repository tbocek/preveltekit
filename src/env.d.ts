/// <reference types="@rsbuild/core/types" />
/// <reference types="svelte" />
declare global {
    interface Window {
        JSDOM: boolean;
    }
}