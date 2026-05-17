<script lang="ts">
  import { onMount } from "svelte";
  import Dashboard from "./lib/dashboard/dashboard.svelte";
  import Header from "./lib/header.svelte";
  import AgentDetail from "./lib/agent/agentDetail.svelte";
  import { router } from "./router.svelte";
  import { agentStore } from "./ws.svelte";

  onMount(() => {
    agentStore.connect();
    return () => agentStore.disconnect();
  });

  const detailMatch = $derived(router.match("/agents/:id"));
</script>

<main class="relative w-full h-screen flex flex-col overflow-hidden">
  <Header />
  {#if detailMatch}
    <AgentDetail id={detailMatch.id} />
  {:else}
    <Dashboard />
  {/if}
</main>
