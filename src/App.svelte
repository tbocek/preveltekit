<script lang="ts">
  import type { Route } from "@mateothegreat/svelte5-router";
  import { Router } from "@mateothegreat/svelte5-router";
  import Dash from './Dash.svelte';
  import Landing from './Landing.svelte';
  const routes: Route[] = [
    {
      path: "/",
      component: Landing
    },
    {
      path: "app",
      component: Dash
    }
  ];
</script>

<Router {routes} />