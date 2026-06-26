import { toast } from "react-toastify";
import { getAgents } from "../../api/agents";
import {
  AlertType,
  SnapshotType,
  StatusChangeType,
  type AlertPayload,
  type BroadcastMessage,
  type SnapshotEvent,
  type StatusChangeEvent,
} from "../../api/models";
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
    const message: BroadcastMessage = JSON.parse(e.data);

    switch (message.type) {
      case SnapshotType:
        const dataSnapshot: SnapshotEvent = message.payload as SnapshotEvent;
        useAgents
          .getState()
          .appendSnapshot(
            dataSnapshot.agent_guid,
            dataSnapshot.snapshot,
            dataSnapshot.agent_name,
          );
        break;
      case StatusChangeType:
        const dataStatus: StatusChangeEvent =
          message.payload as StatusChangeEvent;
        useAgents
          .getState()
          .changeOnline(dataStatus.agent_guid, dataStatus.online);
        break;
      case AlertType:
        const alert: AlertPayload = message.payload as AlertPayload;
        if (alert.active) {
          toast.error(
            `Alert for ${alert.agent_name}. ${alert.metric} is at ${alert.value}`,
          );
        } else {
          toast.success(`Alert resolved for ${alert.agent_name}.`);
        }
        break;
    }
  };

  ws.onclose = () => {
    setTimeout(connectWS, 2000);
  };
}
