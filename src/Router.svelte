<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import type { Routes, RouteParams, Component } from './types';

    // Props for the component
    const props = $props();

    // Create a store for the current route
    let currentRoute = $state<string>('/');

    // Extract routes and notFound from props with defaults
    let routes = $derived<Routes>(props.routes as Routes || []);
    let notFound = $derived<Component>(props.notFound as Component ||
        (() => '<h1>404 - Page Not Found</h1>'));

    // Handle browser back/forward navigation
    const handlePopState = () => {
        currentRoute = window.location.pathname;
    }

    // Define event handler function for custom navigation events
    const handleNavigateEvent = (e: Event) => {
        const customEvent = e as CustomEvent<{ path: string }>;
        currentRoute = customEvent.detail.path;
    };

    onMount(() => {
        // Set initial route
        currentRoute = window.location.pathname;

        // Add event listener for back/forward navigation
        window.addEventListener('popstate', handlePopState);
        window.addEventListener('svelteNavigate', handleNavigateEvent);
    });

    onDestroy(() => {
        window.removeEventListener('popstate', handlePopState);
        window.removeEventListener('svelteNavigate', handleNavigateEvent);
    });

    // Helper function to find matching route and extract params
    function findMatchingRoute(path: string): { component: Component, params: RouteParams } {
        // Normalize path for consistent matching
        const normalizedPath = path.endsWith('/') && path.length > 1
            ? path.slice(0, -1)
            : path;

        // For storing matches and their specificity
        let matchedRoutes: Array<{
            component: Component,
            params: RouteParams,
            specificity: number
        }> = [];

        // Check all routes for matches
        for (const route of routes) {
            const routePath = route.path;
            const routeParts = routePath.split('/');
            const pathParts = normalizedPath.split('/');

            // Skip if parts length doesn't match (unless we have a wildcard)
            if (routeParts.length !== pathParts.length && !routePath.includes('*')) {
                continue;
            }

            let isMatch = true;
            const params: RouteParams = {};
            let specificity = 0; // Higher means more specific match

            // Compare each part of the path
            for (let i = 0; i < routeParts.length; i++) {
                const routePart = routeParts[i];
                const pathPart = pathParts[i];

                // Handle parameter segments (starting with ':')
                if (routePart.startsWith(':')) {
                    const paramName = routePart.slice(1);
                    params[paramName] = pathPart;
                    // Parameters are less specific than exact matches
                    specificity += 5;
                }
                // Handle wildcard segments
                else if (routePart === '*') {
                    // Wildcards are the least specific
                    specificity += 1;
                }
                // Handle exact matches
                else if (routePart === pathPart) {
                    // Exact matches are the most specific
                    specificity += 10;
                }
                else {
                    isMatch = false;
                    break;
                }
            }

            if (isMatch) {
                matchedRoutes.push({ component: route.component, params, specificity });
            }
        }

        // Sort by specificity (highest first) and get the best match
        matchedRoutes.sort((a, b) => b.specificity - a.specificity);

        // Return the most specific match, or the notFound component
        return matchedRoutes.length > 0
            ? { component: matchedRoutes[0].component, params: matchedRoutes[0].params }
            : { component: notFound, params: {} };
    }

    // Current matched route and params
    let matchedRoute = $derived(findMatchingRoute(currentRoute));
    let activeComponent = $derived(matchedRoute.component);
    let routeParams = $derived(matchedRoute.params);
</script>

{#snippet renderer()}
    {#if activeComponent}
        {@render activeComponent({ params: routeParams })}
    {:else}
        {@render notFound()}
    {/if}
{/snippet}

{@render renderer()}

<script module lang="ts">
    // Export the navigate function for programmatic navigation
    export function navigate(path: string): void {
        history.pushState(null, "", path);
        window.dispatchEvent(new CustomEvent('svelteNavigate', { detail: { path } }));
    }

    export function route(node: HTMLAnchorElement): { destroy: () => void } {
        const handleClick = (event: MouseEvent) => {
            // Only handle if it's a left-click without modifier keys
            if (event.button === 0 && !event.ctrlKey && !event.metaKey && !event.altKey && !event.shiftKey) {
                event.preventDefault();

                // Get the href from the anchor tag
                const href = node.getAttribute('href') || '';

                // Handle both relative and absolute paths
                const isExternal =
                    href.startsWith('http://') ||
                    href.startsWith('https://') ||
                    href.startsWith('//') ||
                    node.hasAttribute('external');

                if (!isExternal) {
                    // Normalize the path to handle relative links
                    const path = href.startsWith('/') ? href : `/${href}`;
                    navigate(path);
                } else {
                    // For external links, just follow the href
                    window.location.href = href;
                }
            }
        };

        // Add event listener
        node.addEventListener('click', handleClick);

        // Return the destroy method
        return {
            destroy() {
                node.removeEventListener('click', handleClick);
            }
        };
    }
</script>