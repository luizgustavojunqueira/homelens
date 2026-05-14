<script lang="ts">
  import { onMount } from "svelte";
  import { agentStore } from "../../ws.svelte";
  import AgentCard from "./agentCard.svelte";

  onMount(() => {
    agentStore.connect();
    return () => {
      agentStore.disconnect();
    };
  });
</script>

<div class="p-4 flex flex-col gap-4">
  {#each Object.values(agentStore.agents) as agent (agent.id)}
    <AgentCard {...agent} />
  {/each}
</div>
