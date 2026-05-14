import { getAgents } from "./api/agents";
import type { Agent, SnapshotEntry, SnapshotEvent } from "./api/models";

function createAgentStore() {
  let agents: Record<string, Omit<Agent, "latest_snapshot">> = $state({})
  let snapshots: Record<string, SnapshotEntry> = $state({})

  let ws: WebSocket | null = null;

  async function connect() {
    const initial = await getAgents();
    for (const agent of initial) {
      const { latest_snapshot, ...rest } = agent;
      agents[agent.id] = rest;
      if (latest_snapshot) snapshots[agent.id] = latest_snapshot;
    }

    ws = new WebSocket('ws://localhost:6969/api/agents/ws');

    ws.onmessage = (e) => {
      const message: SnapshotEvent = JSON.parse(e.data);
      snapshots[message.agent_id] = message.snapshot;
    }

    ws.onclose = () => setTimeout(connect, 2000);
  }

  function disconnect() {
    ws?.close();
  }

  return {
    get agents() { return agents; },
    get snapshots() { return snapshots; },
    connect,
    disconnect
  }
}

export const agentStore = createAgentStore();
