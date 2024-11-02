<script lang="ts">
  import type { Route } from "@mateothegreat/svelte5-router";
  import { route, Router } from "@mateothegreat/svelte5-router";
  import Dash from './Dash.svelte';
  import Email from './Email.svelte';
  import Landing from './Landing.svelte';
  const routes: Route[] = [
    {
      path: "/",
      component: Landing,
        pre: (route: Route) => {
            console.log("pre hook fired Landing");
            return route;
        },
        post: () => {
            console.log("post hook fired Landing");
        }
    },
    {
      path: "app",
      component: Dash,
        pre: (route: Route) => {
            console.log("pre hook fired Dash");
            return route;
        },
        post: () => {
            console.log("post hook fired Dash");
        }
    },
      {
          path: "email",
          component: Email,
          pre: (route: Route) => {
              console.log("pre hook fired Email");
              return route;
          },
          post: () => {
              console.log("post hook fired Email");
          }
      }
  ];
</script>
<nav>
    <a use:route href="/">Home</a>
    <a use:route href="/app">Dashboard</a>
    <a use:route href="/email">Email</a>
</nav>

<Router {routes} />