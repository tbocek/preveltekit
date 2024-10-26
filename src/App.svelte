<script lang="ts">
  let count = $state(0);
  function increment() {
    count++;
  }
</script>

<main>
  <div class="content">
    <h1>Rsbuild with Svelte5</h1>
    <p>Start building amazing things with Rsbuild.</p>
    <p>Current count: {count}</p>
    <button class="button" onclick={increment}>
      Click me!
    </button>
  </div>
</main>

<style>
.content {
  display: flex;
  min-height: 100vh;
  line-height: 1.1;
  text-align: center;
  flex-direction: column;
  justify-content: center;
}

.content h1 {
  font-size: 3.6rem;
  font-weight: 700;
}

.content p {
  font-size: 1.2rem;
  font-weight: 400;
  opacity: 0.5;
}
</style>
