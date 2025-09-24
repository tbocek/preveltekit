import { hydrate } from "svelte";
import Index from "./Index.svelte";

hydrate(Index, {
  target: document.getElementById("root")!,
  props: {},
});
