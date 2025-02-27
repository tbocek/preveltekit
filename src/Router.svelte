<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import type { Routes, RouteParams, Component } from './types';

    // Props for the component
    const props = $props();

    // Create a store for the current route
    let currentRoute = $state<string>('/');

    // Extract routes and notFound from props with defaults
    let routes = $derived<Routes>(props.routes as Routes || []);

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

        // Expose routes directly to SSR if running in JSDOM
        if (window?.JSDOM) {
            // Get the current value of the routes
            const currentRoutes = [...routes];
            window.__svelteRoutes = currentRoutes;
        }
    });

    onDestroy(() => {
        window.removeEventListener('popstate', handlePopState);
        window.removeEventListener('svelteNavigate', handleNavigateEvent);
    });

    // Helper function to find matching route and extract params
    function findMatchingRoute(path: string): { component: Component, params: RouteParams } {
        // Normalize the input path (remove trailing slash except for root path)
        const normalizedPath = path === '/' ? '/' : path.endsWith('/') ? path.slice(0, -1) : path;

        // For storing the best match
        let bestMatch: {
            component: Component,
            params: RouteParams,
            specificity: number
        } | null = null;

        // Check all routes for matches
        for (const route of routes) {
            const routePath = route.path;
            let isMatch = false;
            let params: RouteParams = {};
            let specificity = 0;

            // Normalize the route path (remove trailing slash except for root path)
            const normalizedRoutePath = routePath === '/' ? '/' :
                routePath.endsWith('/') ? routePath.slice(0, -1) : routePath;

            // CASE 1: Handle root path special case
            if (normalizedRoutePath === '/') {
                isMatch = normalizedPath === '/';
                specificity = 100; // Highest specificity for root path
            }
            // CASE 2: Handle **/ pattern (matches root path)
            else if (normalizedRoutePath === '**/' || normalizedRoutePath === '**/') {
                isMatch = normalizedPath === '/';
                specificity = 1; // Lowest specificity
            }
            // CASE 3: Handle */ pattern (matches root path and single-segment paths)
            else if (normalizedRoutePath === '*/' || normalizedRoutePath === '*') {
                // Match root path
                if (normalizedPath === '/') {
                    isMatch = true;
                    specificity = 1;
                }
                // Match any single-segment path like /test
                else {
                    const pathSegments = normalizedPath.split('/').filter(Boolean);
                    if (pathSegments.length === 1) {
                        isMatch = true;
                        specificity = 1;
                    }
                }
            }
            // CASE 4: Handle */pattern routes - generic solution
            else if (normalizedRoutePath.startsWith('*/')) {
                // Get the suffix after */
                const suffix = normalizedRoutePath.slice(2);

                // Check if path exactly matches suffix
                if (normalizedPath === '/' + suffix) {
                    isMatch = true;
                    specificity = 2;
                }
                // Check if path ends with /suffix with exactly one segment before
                else {
                    // Build a regex that matches /{segment}/{suffix} exactly
                    const pattern = new RegExp(`^\\/([^\\/]+)\\/${suffix}$`);
                    isMatch = pattern.test(normalizedPath);
                    specificity = 2;
                }
            }
            // CASE 5: Handle other wildcard prefix routes like *path
            else if (normalizedRoutePath.startsWith('*')) {
                const suffix = normalizedRoutePath.slice(1);

                // Simple wildcard matching for other patterns
                if (normalizedPath === suffix ||
                    (suffix.startsWith('/') && normalizedPath.endsWith(suffix))) {
                    isMatch = true;
                    specificity = 2; // Low specificity for wildcard routes
                }
            }
            // CASE 6: Standard path matching with parameters
            else {
                // Split paths into segments for comparison
                const routeSegments = normalizedRoutePath.split('/').filter(Boolean);
                const pathSegments = normalizedPath.split('/').filter(Boolean);

                // For standard routes, segment count must match
                if (routeSegments.length === pathSegments.length) {
                    isMatch = true;
                    specificity = 0;

                    // Compare each segment
                    for (let i = 0; i < routeSegments.length; i++) {
                        const routeSegment = routeSegments[i];
                        const pathSegment = pathSegments[i];

                        // Handle parameter segments
                        if (routeSegment.startsWith(':')) {
                            const paramName = routeSegment.slice(1);
                            params[paramName] = pathSegment;
                            specificity += 5;
                        }
                        // Handle exact matches
                        else if (routeSegment === pathSegment) {
                            specificity += 10;
                        }
                        else {
                            isMatch = false;
                            break;
                        }
                    }
                }
            }

            // Update best match if this route matches and is more specific
            if (isMatch && (!bestMatch || specificity > bestMatch.specificity)) {
                bestMatch = {
                    component: route.component,
                    params,
                    specificity
                };
            }
        }

        // Return the best match or null component
        return bestMatch
            ? { component: bestMatch.component, params: bestMatch.params }
            : { component: null, params: {} };
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
        <h1>404 - Page Not Found for [{currentRoute}]</h1>
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

                // Handle external links
                const isExternal =
                    href.startsWith('http://') ||
                    href.startsWith('https://') ||
                    href.startsWith('//') ||
                    node.hasAttribute('external');

                if (!isExternal) {
                    // Determine the path based on whether it's absolute or relative
                    let path;

                    if (href.startsWith('/')) {
                        // Absolute path - use as is
                        path = href;
                    } else if (href === '' || href === '#') {
                        // Empty href or hash only - stay on current page
                        path = window.location.pathname;
                    } else {
                        // Relative path - combine with current path
                        const currentPath = window.location.pathname;

                        // Ensure the current path ends with a slash if not the root
                        const basePath = currentPath === '/'
                            ? '/'
                            : currentPath.endsWith('/')
                                ? currentPath
                                : currentPath + '/';

                        // Combine base path with relative href
                        path = basePath + href;
                    }

                    // Clean up any double slashes (except after protocol)
                    path = path.replace(/([^:]\/)\/+/g, '$1');

                    // Handle relative paths with ../
                    if (path.includes('../')) {
                        const segments = path.split('/');
                        const cleanSegments = [];

                        for (const segment of segments) {
                            if (segment === '..') {
                                // Go up one level by removing the last segment
                                if (cleanSegments.length > 1) { // Ensure we don't go above root
                                    cleanSegments.pop();
                                }
                            } else if (segment !== '' && segment !== '.') {
                                // Add non-empty segments that aren't current directory
                                cleanSegments.push(segment);
                            }
                        }

                        // Reconstruct the path
                        path = '/' + cleanSegments.join('/');
                    }

                    // Navigate to the computed path
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