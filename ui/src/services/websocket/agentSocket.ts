import { getAgents } from "../../api/agents";
import type { SnapshotEvent } from "../../api/models";
import { useAgents } from "../../store/agentsStore";

export async function connectWS() {
  const initial = await getAgents();

  for (const agent of initial) {
    useAgents
      .getState()
      .appendSnapshot(agent.guid, agent.latest_snapshot!, agent.name);
  }

  const ws = new WebSocket(`/api/agents/ws`);

  ws.onmessage = (e) => {
    const message: SnapshotEvent = JSON.parse(e.data);
    useAgents
      .getState()
      .appendSnapshot(message.agent_guid, message.snapshot, message.agent_name);
  };

  ws.onclose = () => {
    setTimeout(connectWS, 2000);
  };
}
