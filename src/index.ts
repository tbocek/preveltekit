import { hydrate } from 'svelte';
import App from './App.svelte';
import 'picnic';

hydrate(App, {
    target: document.getElementById('root')!,
    props: {},
});