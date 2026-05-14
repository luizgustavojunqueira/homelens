<script lang="ts">
  import { onMount } from "svelte";
  import { agentStore } from "./ws.svelte";

  onMount(() => {
    agentStore.connect();
    return () => {
      agentStore.disconnect();
    };
  });
</script>

<div>
  <h1>HomeLens UI</h1>

  <ul>
    {#each Object.values(agentStore.agents) as agent (agent.id)}
      <li>
        {agent.name} - {agent.online ? "Online" : "Offline"} - {agent.last_seen}
        - {agentStore.snapshots[agent.id]?.data?.cpu_usage?.cpu_avg}
      </li>
    {/each}
  </ul>
</div>
