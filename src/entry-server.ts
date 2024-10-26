import App from './App.svelte';
import { render as renderApp } from 'svelte/server';

export function render() {
    return renderApp(App, {
        props: {}
    });
}