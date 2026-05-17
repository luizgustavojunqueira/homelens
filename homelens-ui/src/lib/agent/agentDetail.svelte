<script lang="ts">
  import Icon from "@iconify/svelte";
  import { agentStore } from "../../ws.svelte";
  import { router } from "../../router.svelte";

  let { id }: { id: string } = $props();

  const agent = $derived(agentStore.agents[id]);
</script>

<section class="px-6 py-6 flex-1 overflow-y-auto">
  <button
    class="flex items-center gap-2 text-sm text-(--text-dim) hover:text-(--text) transition-colors mb-5 cursor-pointer"
    onclick={() => router.go("/")}
  >
    <Icon icon="lucide:arrow-left" width="16" />
    Fleet
  </button>

  {#if agent}
    <div class="flex items-center gap-3 mb-6">
      <span class="dot {agent.online ? '' : 'off'}"></span>
      <h2 class="text-2xl font-medium text-(--text)">{agent.name}</h2>
    </div>

    <div class="text-(--text-dim) text-sm">
      Detail view — id <span class="tnum text-(--text)">{id}</span>
    </div>
  {:else}
    <div class="text-(--text-faint)">unknown agent</div>
  {/if}
</section>
