<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { setupStaticRoutes } from "./router";
    import type { Routes, Component } from "./types";

    interface Props {
        routes: Routes;
    }
    
    let { routes }: Props = $props();

    let currentRoute = $state<string>("/");

    const handlePopState = () => {
        currentRoute = window.location.pathname;
    };

    const handleNavigateEvent = (e: Event) => {
        const customEvent = e as CustomEvent<{ path: string }>;
        currentRoute = customEvent.detail.path;
    };

    onMount(() => {
        currentRoute = new URL(window.location.href).pathname;

        // Add event listener for back/forward navigation
        window.addEventListener("popstate", handlePopState);
        window.addEventListener("svelteNavigate", handleNavigateEvent);

        // Expose routes directly to SSR if running in JSDOM
        setupStaticRoutes(routes);
    });

    onDestroy(() => {
        window.removeEventListener("popstate", handlePopState);
        window.removeEventListener("svelteNavigate", handleNavigateEvent);
    });

    // Helper function to find matching route and extract params
    function findMatchingRoute(path: string): {
        component: Component;
        params: Record<string, string>;
    } {
        // Normalize the input path (remove trailing slash except for root path)
        const normalizedPath =
            path === "/" ? "/" : path.endsWith("/") ? path.slice(0, -1) : path;

        // For storing the best match
        let bestMatch: {
            component: Component;
            params: Record<string, string>;
            specificity: number;
        } | null = null;

        if (routes.dynamicRoutes) {
            // Check all routes for matches
            for (const route of routes.dynamicRoutes) {
                const routePath = route.path;
                let isMatch = false;
                let params: Record<string, string> = {};
                let specificity = 0;

                // Normalize the route path (remove trailing slash except for root path)
                const normalizedRoutePath =
                    routePath === "/"
                        ? "/"
                        : routePath.endsWith("/")
                          ? routePath.slice(0, -1)
                          : routePath;

                // CASE 1: Handle root path special case
                if (normalizedRoutePath === "/") {
                    isMatch = normalizedPath === "/";
                    specificity = 100; // Highest specificity for root path
                }
                // CASE 2: Handle **/ pattern (matches root path)
                else if (
                    normalizedRoutePath === "**/" ||
                    normalizedRoutePath === "**/"
                ) {
                    isMatch = normalizedPath === "/";
                    specificity = 1; // Lowest specificity
                }
                // CASE 3: Handle */ pattern (matches root path and single-segment paths)
                else if (
                    normalizedRoutePath === "*/" ||
                    normalizedRoutePath === "*"
                ) {
                    // Match root path
                    if (normalizedPath === "/") {
                        isMatch = true;
                        specificity = 1;
                    }
                    // Match any single-segment path like /test
                    else {
                        const pathSegments = normalizedPath
                            .split("/")
                            .filter(Boolean);
                        if (pathSegments.length === 1) {
                            isMatch = true;
                            specificity = 1;
                        }
                    }
                }
                // CASE 4: Handle */pattern routes - generic solution
                else if (normalizedRoutePath.startsWith("*/")) {
                    // Get the suffix after */
                    const suffix = normalizedRoutePath.slice(2);

                    // Check if path exactly matches suffix
                    if (normalizedPath === "/" + suffix) {
                        isMatch = true;
                        specificity = 2;
                    }
                    // Check if path ends with /suffix with exactly one segment before
                    else {
                        // Build a regex that matches /{segment}/{suffix} exactly
                        const pattern = new RegExp(
                            `^\\/([^\\/]+)\\/${suffix}$`,
                        );
                        isMatch = pattern.test(normalizedPath);
                        specificity = 2;
                    }
                }
                // CASE 5: Handle other wildcard prefix routes like *path
                else if (normalizedRoutePath.startsWith("*")) {
                    const suffix = normalizedRoutePath.slice(1);

                    // Simple wildcard matching for other patterns
                    if (
                        normalizedPath === suffix ||
                        (suffix.startsWith("/") &&
                            normalizedPath.endsWith(suffix))
                    ) {
                        isMatch = true;
                        specificity = 2; // Low specificity for wildcard routes
                    }
                }
                // CASE 6: Standard path matching with parameters
                else {
                    // Split paths into segments for comparison
                    const routeSegments = normalizedRoutePath
                        .split("/")
                        .filter(Boolean);
                    const pathSegments = normalizedPath
                        .split("/")
                        .filter(Boolean);

                    // For standard routes, segment count must match
                    if (routeSegments.length === pathSegments.length) {
                        isMatch = true;
                        specificity = 0;

                        // Compare each segment
                        for (let i = 0; i < routeSegments.length; i++) {
                            const routeSegment = routeSegments[i];
                            const pathSegment = pathSegments[i];

                            // Handle parameter segments
                            if (routeSegment.startsWith(":")) {
                                const paramName = routeSegment.slice(1);
                                params[paramName] = pathSegment;
                                specificity += 5;
                            }
                            // Handle exact matches
                            else if (routeSegment === pathSegment) {
                                specificity += 10;
                            } else {
                                isMatch = false;
                                break;
                            }
                        }
                    }
                }

                // Update best match if this route matches and is more specific
                if (
                    isMatch &&
                    (!bestMatch || specificity > bestMatch.specificity)
                ) {
                    bestMatch = {
                        component: route.component,
                        params,
                        specificity,
                    };
                }
            }
        }

        // Return the best match or null component
        return bestMatch
            ? { component: bestMatch.component, params: bestMatch.params }
            : { component: null, params: {} };
    }

    // Current matched route and params
    let matchedRoute = $derived(findMatchingRoute(currentRoute));
    let ActiveComponent = $derived(matchedRoute.component);
    let routeParams = $derived(matchedRoute.params);
</script>

{#if ActiveComponent}
    <ActiveComponent params={routeParams}></ActiveComponent>
{:else}
    <h1>404 - Page Not Found for [{currentRoute}]</h1>
{/if}
