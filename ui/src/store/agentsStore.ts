import { create } from "zustand";
import type { Agent, SnapshotEntry } from "../api/models";

type AgentState = Agent & {
  history: SnapshotEntry[];
}

interface AgentsStore {
  agents: Record<string, AgentState>;
  appendSnapshot: (agent_id: string, snapshot: SnapshotEntry) => void;
  getAgentState: (agent_id: string) => AgentState | undefined;
  insertHistory: (agent_id: string, snapshots: SnapshotEntry[]) => void;
}

export const useAgents = create<AgentsStore>((set) => ({
  agents: {},
  appendSnapshot: (agent_id: string, snapshot: SnapshotEntry) => {
    set((state) => {
      const agentState: AgentState = state.agents[agent_id];
      if (!agentState) {
        return {
          agents: {
            ...state.agents,
            [agent_id]: {
              id: agent_id,
              name: agent_id,
              last_seen: String(snapshot.timestamp),
              online: true,
              latest_snapshot: snapshot,
              history: [snapshot],
            },
          },
        };
      }

      const updatedHistory = [...agentState.history, snapshot];
      return {
        agents: {
          ...state.agents,
          [agent_id]: {
            ...agentState,
            latest_snapshot: snapshot,
            history: updatedHistory,
          },
        },
      };

    })
  },
  getAgentState: (agent_id: string) => {
    const agentState: AgentState = useAgents.getState().agents[agent_id];
    return agentState;
  },
  insertHistory: (agent_id: string, snapshots: SnapshotEntry[]) => {
    set((state) => {
      const currentOldest = state.agents[agent_id]?.history[0];
      const filtered = snapshots.filter(snap => {
        if (!currentOldest) return true;
        return snap.timestamp < currentOldest.timestamp;
      })
      const agentState: AgentState = state.agents[agent_id];
      if (!agentState) return state;
      return {
        agents: {
          ...state.agents,
          [agent_id]: {
            ...agentState,
            history: [...filtered, ...agentState.history],
          }
        }
      }
    })
  }
}))

