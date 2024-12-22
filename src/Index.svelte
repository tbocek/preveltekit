<script lang="ts">
    import type { Route } from "@mateothegreat/svelte5-router";
    import { route, Router } from "@mateothegreat/svelte5-router";
    import Landing from "./Landing.svelte";
    import Dash from "./Dash.svelte";
    import Email from "./Email.svelte";

    let isMenuOpen = $state(false);

    const routes: Route[] = [
        {
            path: "/",
            component: Landing,
            pre: (route: Route) => {
                console.info("Pre-hook: Landing page");
                return route;
            },
            post: () => {
                console.info("Post-hook: Landing page loaded");
            }
        },
        {
            path: "app",
            component: Dash,
            pre: (route: Route) => {
                console.info("Pre-hook: Dashboard");
                return route;
            },
            post: () => {
                console.info("Post-hook: Dashboard loaded");
            }
        },
        {
            path: "email",
            component: Email,
            pre: (route: Route) => {
                console.info("Pre-hook: Email builder");
                return route;
            },
            post: () => {
                console.info("Post-hook: Email builder loaded");
            }
        }
    ];

    function toggleMenu() {
        isMenuOpen = !isMenuOpen;
    }

    $effect(() => {
        console.log("LightKit router initialized");
    });
</script>

<div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4">
            <div class="flex justify-between h-16">
                <div class="flex">
                    <div class="flex-shrink-0 flex items-center">
                        <span class="text-xl font-bold text-blue-600">LightKit</span>
                    </div>

                    <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
                        <a use:route href="/" class="inline-flex items-center px-1 pt-1 text-gray-700 hover:text-blue-600">
                            Home
                        </a>
                        <a use:route href="/app" class="inline-flex items-center px-1 pt-1 text-gray-700 hover:text-blue-600">
                            Dashboard
                        </a>
                        <a use:route href="/email" class="inline-flex items-center px-1 pt-1 text-gray-700 hover:text-blue-600">
                            Email Builder
                        </a>
                    </div>
                </div>

                <div class="flex items-center sm:hidden">
                    <button
                            onclick={toggleMenu}
                            class="inline-flex items-center justify-center p-2 rounded-md text-gray-700 hover:text-blue-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
                    >
                        <span class="sr-only">Open main menu</span>
                        {#if isMenuOpen}
                            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        {:else}
                            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                            </svg>
                        {/if}
                    </button>
                </div>
            </div>
        </div>

        {#if isMenuOpen}
            <div class="sm:hidden">
                <div class="pt-2 pb-3 space-y-1">
                    <a
                            use:route
                            href="/"
                            class="block pl-3 pr-4 py-2 text-base font-medium text-gray-700 hover:text-blue-600 hover:bg-gray-50"
                            onclick={() => isMenuOpen = false}
                    >
                        Home
                    </a>
                    <a
                            use:route
                            href="/app"
                            class="block pl-3 pr-4 py-2 text-base font-medium text-gray-700 hover:text-blue-600 hover:bg-gray-50"
                            onclick={() => isMenuOpen = false}
                    >
                        Dashboard
                    </a>
                    <a
                            use:route
                            href="/email"
                            class="block pl-3 pr-4 py-2 text-base font-medium text-gray-700 hover:text-blue-600 hover:bg-gray-50"
                            onclick={() => isMenuOpen = false}
                    >
                        Email Builder
                    </a>
                </div>
            </div>
        {/if}
    </nav>

    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="px-4 py-4 sm:px-0">
            <Router {routes} />
        </div>
    </main>

    <footer class="bg-white mt-12">
        <div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
            <p class="text-center text-gray-500 text-sm">
                Built with LightKit - Powered by Svelte 5 and Rsbuild
            </p>
        </div>
    </footer>
</div>