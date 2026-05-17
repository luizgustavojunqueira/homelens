<script lang="ts">
  import Icon from "@iconify/svelte";
  import type { Agent } from "../../api/models";
  import { router } from "../../router.svelte";

  let { id, latest_snapshot, name, online }: Agent = $props();

  const snap = $derived(latest_snapshot?.data);
  const cpuPct = $derived(snap?.cpu_usage?.cpu_avg ?? 0);
  const memUsed = $derived((snap?.memory?.used ?? 0) * 1024);
  const memTotal = $derived((snap?.memory?.total ?? 0) * 1024);
  const memPct = $derived(memTotal > 0 ? (memUsed / memTotal) * 100 : 0);
  const diskPct = $derived(snap?.disk_space?.usage_percent ?? 0);
  const temp = $derived(snap?.temperature?.temp_avg ?? 0);

  function level(pct: number, warn = 70, crit = 90): "" | "warn" | "crit" {
    if (pct >= crit) return "crit";
    if (pct >= warn) return "warn";
    return "";
  }
  function tempLevel(t: number): "" | "warn" | "crit" {
    if (t >= 80) return "crit";
    if (t >= 65) return "warn";
    return "";
  }
</script>

<button
  class="w-full grid grid-cols-[minmax(180px,1.4fr)_repeat(4,minmax(0,1fr))_24px] items-center gap-6 px-5 py-4 text-left hover:bg-(--bg-hover) transition-colors cursor-pointer"
  onclick={() => router.go(`/agents/${id}`)}
>
  <div class="flex items-center gap-3 min-w-0">
    <span class="dot {online ? '' : 'off'}"></span>
    <span class="text-base font-medium text-(--text) truncate">{name}</span>
  </div>

  {#if snap}
    <div class="flex items-center gap-2 min-w-0">
      <span class="text-sm text-(--text-dim) w-10">CPU</span>
      <div class="bar {level(cpuPct)} flex-1">
        <span style="width: {Math.min(100, cpuPct)}%"></span>
      </div>
      <span class="tnum text-base text-(--text) w-12 text-right"
        >{cpuPct.toFixed(0)}%</span
      >
    </div>

    <div class="flex items-center gap-2 min-w-0">
      <span class="text-sm text-(--text-dim) w-10">MEM</span>
      <div class="bar {level(memPct)} flex-1">
        <span style="width: {Math.min(100, memPct)}%"></span>
      </div>
      <span class="tnum text-base text-(--text) w-12 text-right"
        >{memPct.toFixed(0)}%</span
      >
    </div>

    <div class="flex items-center gap-2 min-w-0">
      <span class="text-sm text-(--text-dim) w-10">DSK</span>
      <div class="bar {level(diskPct, 80, 95)} flex-1">
        <span style="width: {Math.min(100, diskPct)}%"></span>
      </div>
      <span class="tnum text-base text-(--text) w-12 text-right"
        >{diskPct.toFixed(0)}%</span
      >
    </div>

    <div class="flex items-center gap-2 min-w-0">
      <span class="text-sm text-(--text-dim) w-10">TMP</span>
      <div class="bar {tempLevel(temp)} flex-1">
        <span style="width: {Math.min(100, temp)}%"></span>
      </div>
      <span class="tnum text-base text-(--text) w-12 text-right"
        >{temp.toFixed(0)}°</span
      >
    </div>
  {:else}
    <div class="col-span-4 text-sm text-(--text-faint)">no data</div>
  {/if}

  <Icon icon="lucide:chevron-right" width="16" class="text-(--text-faint)" />
</button>
