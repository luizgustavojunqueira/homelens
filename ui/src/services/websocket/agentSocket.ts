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

let ws: WebSocket | null = null;
let retryCount = 0;
const MAX_RETRIES = 20;
const BASE_DELAY = 1000;
const MAX_DELAY = 30000;
let reconnectTimer: number | null = null;

export function connectWS() {
  if (ws?.readyState === WebSocket.CONNECTING || ws?.readyState === WebSocket.OPEN) {
    return;
  }

  const isFirstConnect = retryCount === 0;

  const connect = async () => {
    if (isFirstConnect) {
      try {
        const initial = await getAgents();
        for (const agent of initial) {
          useAgents
            .getState()
            .appendSnapshot(agent.guid, agent.latest_snapshot!, agent.name);
        }
      } catch {
        scheduleReconnect();
        return;
      }
    }

    ws = new WebSocket(`/api/agents/ws`);

    ws.onopen = () => {
      retryCount = 0;
    };

    ws.onmessage = (e) => {
      try {
        const message: BroadcastMessage = JSON.parse(e.data);

        switch (message.type) {
          case SnapshotType: {
            const dataSnapshot: SnapshotEvent = message.payload as SnapshotEvent;
            useAgents
              .getState()
              .appendSnapshot(
                dataSnapshot.agent_guid,
                dataSnapshot.snapshot,
                dataSnapshot.agent_name,
              );
            break;
          }
          case StatusChangeType: {
            const dataStatus: StatusChangeEvent =
              message.payload as StatusChangeEvent;
            useAgents
              .getState()
              .changeOnline(dataStatus.agent_guid, dataStatus.online);
            break;
          }
          case AlertType: {
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
        }
      } catch (err) {
        console.error("Failed to parse websocket message", err);
      }
    };

    ws.onerror = () => {
      console.error("WebSocket error");
    };

    ws.onclose = () => {
      ws = null;
      scheduleReconnect();
    };
  };

  connect();
}

function scheduleReconnect() {
  if (retryCount >= MAX_RETRIES) {
    console.error("WebSocket: max retries reached");
    return;
  }
  const delay = Math.min(BASE_DELAY * 2 ** retryCount, MAX_DELAY);
  retryCount++;
  if (reconnectTimer) window.clearTimeout(reconnectTimer);
  reconnectTimer = window.setTimeout(connectWS, delay);
}

export function disconnectWS() {
  if (reconnectTimer) {
    window.clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
  if (ws) {
    ws.onclose = null; // Prevent reconnect loop on intentional disconnect
    ws.close();
    ws = null;
  }
}
