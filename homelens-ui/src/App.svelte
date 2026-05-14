<script lang="ts">
  import { onMount } from "svelte";
  import { getAgents } from "./api/agents";
  import type { Agent } from "./api/agents";

  let agents: Agent[] = $state([]);

  onMount(() => {
    getAgents().then((res) => {
      agents = res;
    });
  });
</script>

<div>
  <h1>HomeLens UI</h1>

  <ul>
    {#each agents as agent}
      <li>
        {agent.name} - {agent.online ? "Online" : "Offline"} - {agent.last_seen}
        - {agent.latest_snapshot.data.cpu_usage.cpu_avg}
      </li>
    {/each}
  </ul>
</div>
