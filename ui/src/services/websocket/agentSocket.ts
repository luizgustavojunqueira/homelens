import { getAgents } from "../../api/agents";
import type { SnapshotEvent } from "../../api/models";
import { useAgents } from "../../store/agentsStore";

export async function connectWS() {
  const initial = await getAgents();

  for (const agent of initial) {
    useAgents.getState().appendSnapshot(agent.id, agent.latest_snapshot!);
  }

  const ws = new WebSocket('ws://localhost:6969/api/agents/ws');

  ws.onmessage = (e) => {
    const message: SnapshotEvent = JSON.parse(e.data);
    useAgents.getState().appendSnapshot(message.agent_id, message.snapshot);
  };

  ws.onclose = () => {
    setTimeout(connectWS, 2000);
  };
}

