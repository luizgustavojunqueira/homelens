import { getAgents } from "./api/agents";
import type { Agent, SnapshotEvent } from "./api/models";

function createAgentStore() {
  let agents: Record<string, Agent> = $state({})

  let ws: WebSocket | null = null;

  async function connect() {
    const initial = await getAgents();
    for (const agent of initial) {
      agents[agent.id] = agent;
    }

    ws = new WebSocket('ws://localhost:6969/api/agents/ws');

    ws.onmessage = (e) => {
      const message: SnapshotEvent = JSON.parse(e.data);
      agents[message.agent_id].latest_snapshot = message.snapshot;
    }

    ws.onclose = () => setTimeout(connect, 2000);
  }

  function disconnect() {
    ws?.close();
  }

  return {
    get agents() { return agents; },
    connect,
    disconnect
  }
}

export const agentStore = createAgentStore();
