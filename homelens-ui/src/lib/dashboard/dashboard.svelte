<script lang="ts">
  import { agentStore } from "../../ws.svelte";
  import AgentCard from "./agentCard.svelte";

  const agents = $derived(Object.values(agentStore.agents));
  const onlineCount = $derived(agents.filter((a) => a.online).length);
</script>

<section class="px-6 py-6 flex-1 overflow-y-auto">
  <div class="flex items-baseline gap-3 mb-4">
    <h2 class="text-lg font-medium text-(--text)">Fleet</h2>
    <span class="tnum text-base text-(--text-dim)">
      {onlineCount}/{agents.length} online
    </span>
  </div>

  <div
    class="border border-(--border) rounded-md overflow-hidden bg-(--bg-elev) divide-y divide-(--border)"
  >
    {#each agents as agent (agent.id)}
      <AgentCard {...agent} />
    {/each}
  </div>
</section>
