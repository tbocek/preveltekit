<script lang="ts">
  import type { Route } from "@mateothegreat/svelte5-router";
  import { route, Router } from "@mateothegreat/svelte5-router";
  import Landing from "./Landing.svelte";
  import Dash from "./Dash.svelte";
  import Email from "./Email.svelte";

  const routes: Route[] = [
    {
      path: "/",
      component: Landing,
        pre: (route: Route) => {
            console.info("pre hook fired Landing");
            return route;
        },
        post: () => {
            console.info("post hook fired Landing");
        }
    },
    {
      path: "app",
      component: Dash,
        pre: (route: Route) => {
            console.info("pre hook fired Dash");
            return route;
        },
        post: () => {
            console.info("post hook fired Dash");
        }
    },
      {
          path: "email",
          component: Email,
          pre: (route: Route) => {
              console.info("pre hook fired Email");
              return route;
          },
          post: () => {
              console.info("post hook fired Email");
          }
      }
  ];
</script>
<style>
    nav {
        position: relative;
    }
</style>
<div>
<nav>
    <a use:route href="/">Home</a>
    <a use:route href="/app">Dashboard</a>
    <a use:route href="/email" class="button icon-puzzle">Email</a>
</nav>

<div class="container">
<Router {routes} />
</div>
</div>