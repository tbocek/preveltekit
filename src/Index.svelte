<script lang="ts">
    import type { Route } from "@mateothegreat/svelte5-router";
    import { Router } from "@mateothegreat/svelte5-router";
    import Landing from "./Landing.svelte";
    import Documentation from "./Documentation.svelte";
    import Example from "./Example.svelte";

    const routes: Route[] = [
        {
            path: "/",
            component: Landing
        },
        {
            path: "doc",
            component: Documentation
        },
        {
            path: "example",
            component: Example
        }
    ];
</script>

<style>
    .app {
        min-height: 100vh;
        background-color: #f7fafc;
    }

    nav {
        background-color: white;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .nav-content {
        max-width: 80rem;
        margin: 0 auto;
        padding: 0 1rem;
    }

    .nav-items {
        display: flex;
        align-items: center;
        height: 4rem;
        gap: 3rem;
    }

    .logo {
        font-size: 1.5rem;
        font-weight: bold;
        color: #4299e1;
    }

    .nav-links {
        display: flex;
        gap: 3rem;
    }

    .nav-links a {
        color: #4299e1;
        font-size: 1.1rem;
        font-weight: 500;
        text-decoration: none;
        transition: color 0.2s;
        padding: 0.5rem 0;
    }

    .nav-links a:hover {
        color: #2b6cb0;
    }

    main {
        max-width: 80rem;
        margin: 0 auto;
        padding: 1.5rem 1rem;
    }

    .content {
        padding: 1rem 0;
    }

    footer {
        background-color: white;
        margin-top: 3rem;
    }

    .footer-content {
        max-width: 80rem;
        margin: 0 auto;
        padding: 1.5rem 1rem;
        text-align: center;
        color: #718096;
        font-size: 0.875rem;
    }
</style>

<div class="app">
    <nav>
        <div class="nav-content">
            <div class="nav-items">
                <span class="logo">PrevelteKit</span>
                <div class="nav-links">
                    <a href="./">Home</a>
                    <a href="./doc">Documentation</a>
                    <a href="./example">Example</a>
                </div>
            </div>
        </div>
    </nav>

    <main>
        <div class="content">
            <Router {routes} />
        </div>
    </main>

    <footer>
        <div class="footer-content">
            <p><span class="highlight">PrevelteKit</span> • Lightning fast builds • Pre-rendered SPA • Pure HTML/CSS/JS output</p>
        </div>
    </footer>
</div>

