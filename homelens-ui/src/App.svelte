<script lang="ts">
  import { onMount } from "svelte";
  import { getAgents, type Agent } from "./api/agents";
  import { getSnapshots } from "./api/snapshots";

  let agents: Agent[] = $state([]);

  onMount(() => {
    getAgents().then((response) => {
      console.log("Agents:", response);
      agents = response;
    });

    getSnapshots().then((res) => {
      console.log(res);
    });
  });
</script>

<div>
  <h1>HomeLens UI</h1>
  <p>
    Welcome to the HomeLens UI! This is a placeholder for the main application
    interface.
  </p>

  <ul>
    {#each agents as agent}
      <li>
        {agent.name} - {agent.online ? "Online" : "Offline"} - {agent.last_seen}
      </li>
    {/each}
  </ul>
</div>
