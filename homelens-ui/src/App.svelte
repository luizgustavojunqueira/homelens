<script lang="ts">
  import { onMount } from "svelte";
  import { getAgents } from "./api/agents";
  import type { Agent } from "./api/agents";
  import { getSnapshots } from "./api/snapshots";
  import type { AgentSnapshots } from "./api/snapshots";

  let agents: Agent[] = $state([]);
  let agentSnapshots: AgentSnapshots[] = $state([]);

  onMount(() => {
    getAgents().then((res) => {
      agents = res;
    });

    getSnapshots().then((res) => {
      agentSnapshots = res.agents;
    });
  });

  function getLatestSnapshot(agentId: string) {
    const entry = agentSnapshots.find((s) => s.agent_id === agentId);
    if (!entry || entry.snapshots.length === 0) return null;
    return entry.snapshots[entry.snapshots.length - 1];
  }
</script>

<div>
  <h1>HomeLens UI</h1>

  <ul>
    {#each agents as agent}
      {@const latest = getLatestSnapshot(agent.id)}
      <li>
        {agent.name} - {agent.online ? "Online" : "Offline"} - {agent.last_seen}
        {latest ? `CPU: ${latest.data.cpu_usage.cpu_avg.toFixed(1)}%` : "No snapshot"}
      </li>
    {/each}
  </ul>
</div>
