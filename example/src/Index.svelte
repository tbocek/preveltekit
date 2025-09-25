<script lang="ts">
    import './main.css';
    import Router from 'preveltekit/Router.svelte';
    import type { Routes } from 'preveltekit/types';
    import Landing from "./Landing.svelte";
    import Documentation from "./Documentation.svelte";
    import Example from "./Example.svelte";

    const routes: Routes = [
        {
            path: "*/doc", // Dynamic path maching, e.g., https://tbocek.github.io/preveltekit/doc matches
            component: Documentation, // Component to display
            static: "doc.html", // For static pre-rendering, the filename
            navigation: "/doc" // For static pre-rendering, we need to have a URL for jsdom. 
                               // Here /doc will be used for http://localhost/doc, which will match path */doc
                               // Depending on your component needs you can use e.g., /preveltekit/doc here.
        },
        {
            path: "*/example",
            component: Example,
            static: "example.html",
            navigation: "/example"
        },
        {
            path: "*/",
            component: Landing,
            static: "index.html",
            navigation: "/"
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

    .server-side-indicator {
        color: #93c5fd;
        font-size: 0.875rem;
        margin-left: 0.5rem;
        align-items: center;
        display: flex;
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
                    <a href="https://github.com/tbocek/preveltekit/" target="_blank" rel="noopener noreferrer" aria-label="Link to project on GitHub">
                      <svg height="24" viewBox="0 0 16 16" width="24" style="fill: #4299e1;">
                        <path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"></path>
                      </svg>
                    </a>
                    <div class="server-side-indicator">(These links here are server-side routing • Page reloads)</div>
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

