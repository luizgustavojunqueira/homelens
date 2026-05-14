<script lang="ts">
  import Icon from "@iconify/svelte";
  import type { Agent } from "../../api/models";
  import {
    formatBytes,
    formatDate,
    formatTemperature,
    parseDate,
  } from "../../utils/utils";

  let { id, latest_snapshot, name, last_seen, online }: Agent = $props();
</script>

<div
  class="w-full bg-white text-gray-600 shadow-md rounded-lg p-2 flex flex-row gap-4 justify-between items-center"
>
  <h3 class="text-lg font-semibold flex items-center gap-2">
    <Icon
      icon="mdi:circle"
      class={online ? "text-green-500" : "text-gray-400"}
      width="12"
    />
    {name}

    <span class="text-sm text-gray-400"
      >({formatDate(parseDate(last_seen))})</span
    >
  </h3>
  {#if latest_snapshot}
    <p>
      CPU: <span class="whitespace-pre"
        >{latest_snapshot.data.cpu_usage.cpu_avg.toFixed(2).padStart(6)}</span
      >%
    </p>
    <p>
      Mem: {formatBytes(latest_snapshot.data.memory.used)} / {formatBytes(
        latest_snapshot.data.memory.total,
      )}
    </p>

    <p>
      Disk: {formatBytes(latest_snapshot.data.disk_space.used)} / {formatBytes(
        latest_snapshot.data.disk_space.total,
      )}
    </p>

    <p>
      Temp: <span class="whitespace-pre"
        >{formatTemperature(latest_snapshot.data.temperature.temp_avg)}</span
      >
    </p>
  {/if}

  <button
    class="flex items-center justify-center p-2 rounded hover:bg-gray-100"
    onclick={() => console.log(`Viewing details for agent ${id}`)}
    aria-label="View details"
  >
    <Icon icon="mdi:arrow-right" width="20" />
  </button>
</div>
