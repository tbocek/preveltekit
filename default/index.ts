import { hydrate } from 'svelte';
import Index from './index.svelte';

hydrate(Index, {
    target: document.getElementById('root')!,
    props: {},
});