<script lang="ts">
  import Header from "./Header.svelte";
</script>

<main>
  <div class="flex min-h-screen flex-col">
    <Header />
    <div class="flex-1">content</div>
    <div
      class="relative bottom-0 w-full p-4 text-center font-mono text-sm text-stone-400"
    >
      <p>some important info</p>
    </div>
  </div>
</main>
