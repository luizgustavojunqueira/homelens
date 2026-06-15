import { create } from "zustand";
import type { Agent, SnapshotEntry } from "../api/models";

type AgentState = Agent & {
  history: SnapshotEntry[];
};

interface AgentsStore {
  agents: Record<string, AgentState>;
  appendSnapshot: (guid: string, snapshot: SnapshotEntry, name: string) => void;
  getAgentState: (guid: string) => AgentState | undefined;
  insertHistory: (guid: string, snapshots: SnapshotEntry[]) => void;
}

export const useAgents = create<AgentsStore>((set) => ({
  agents: {},
  appendSnapshot: (
    guid: string,
    snapshot: SnapshotEntry,
    name: string = "",
  ) => {
    set((state) => {
      const agentState: AgentState = state.agents[guid];
      if (!agentState) {
        return {
          agents: {
            ...state.agents,
            [guid]: {
              guid: guid,
              name: name,
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
          [guid]: {
            ...agentState,
            latest_snapshot: snapshot,
            history: updatedHistory,
          },
        },
      };
    });
  },
  getAgentState: (guid: string) => {
    const agentState: AgentState = useAgents.getState().agents[guid];
    return agentState;
  },
  insertHistory: (guid: string, snapshots: SnapshotEntry[]) => {
    set((state) => {
      const currentOldest = state.agents[guid]?.history[0];
      const filtered = snapshots.filter((snap) => {
        if (!currentOldest) return true;
        return snap.timestamp < currentOldest.timestamp;
      });
      const agentState: AgentState = state.agents[guid];
      if (!agentState) return state;
      return {
        agents: {
          ...state.agents,
          [guid]: {
            ...agentState,
            history: [...filtered, ...agentState.history],
          },
        },
      };
    });
  },
}));
