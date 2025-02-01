import { hydrate } from 'svelte';
import Index from './Index.svelte';
import './main.css';

hydrate(Index, {
    target: document.getElementById('root')!,
    props: {},
});