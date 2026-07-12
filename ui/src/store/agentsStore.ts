import { create } from "zustand";
import type { Agent, SnapshotEntry } from "../api/models";

const MAX_HISTORY = 500;

type AgentState = Agent & {
  history: SnapshotEntry[];
};

interface AgentsStore {
  agents: Record<string, AgentState>;
  appendSnapshot: (guid: string, snapshot: SnapshotEntry, name: string) => void;
  insertHistory: (guid: string, snapshots: SnapshotEntry[]) => void;
  changeOnline: (guid: string, online: boolean) => void;
}

export const useAgents = create<AgentsStore>((set) => ({
  agents: {},
  changeOnline: (guid: string, online: boolean) => {
    set((state) => {
      const agentState: AgentState = state.agents[guid];

      if (!agentState) {
        return state;
      }
      return {
        agents: {
          ...state.agents,
          [guid]: {
            ...agentState,
            online: online,
          },
        },
      };
    });
  },
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

      const updatedHistory = [...agentState.history, snapshot].slice(-MAX_HISTORY);
      return {
        agents: {
          ...state.agents,
          [guid]: {
            ...agentState,
            name: name,
            online: true,
            last_seen: String(snapshot.timestamp),
            latest_snapshot: snapshot,
            history: updatedHistory,
          },
        },
      };
    });
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
      
      const newHistory = [...filtered, ...agentState.history].slice(-MAX_HISTORY);
      
      return {
        agents: {
          ...state.agents,
          [guid]: {
            ...agentState,
            history: newHistory,
          },
        },
      };
    });
  },
}));
