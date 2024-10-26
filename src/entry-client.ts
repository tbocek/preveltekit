import { hydrate } from 'svelte';
import App from './App.svelte';
import './index.css';

hydrate(App, {
    target: document.getElementById('root')!,
    props: {},
});