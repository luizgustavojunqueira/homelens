import { createContext, use, useEffect, useRef, useState, type ReactNode } from 'react';
import { getAgents } from '../api/agents';
import type { Agent, SnapshotEvent } from '../api/models';

type AgentMap = Record<string, Agent>;

const WebSocketContext = createContext<AgentMap | null>(null);

export function WebSocketProvider({ children }: { children: ReactNode }) {
  const [agents, setAgents] = useState<AgentMap>({});
  const wsRef = useRef<WebSocket | null>(null);
  const closedRef = useRef(false);

  useEffect(() => {
    closedRef.current = false;

    async function connect() {
      const initial = await getAgents();
      setAgents((prev) => {
        const next = { ...prev };
        for (const agent of initial) next[agent.id] = agent;
        return next;
      });

      const ws = new WebSocket('ws://localhost:6969/api/agents/ws');
      wsRef.current = ws;

      ws.onmessage = (e) => {
        const message: SnapshotEvent = JSON.parse(e.data);
        setAgents((prev) => {
          const current = prev[message.agent_id];
          if (!current) return prev;
          return {
            ...prev,
            [message.agent_id]: { ...current, latest_snapshot: message.snapshot },
          };
        });
      };

      ws.onclose = () => {
        if (closedRef.current) return;
        setTimeout(connect, 2000);
      };
    }

    connect();

    return () => {
      closedRef.current = true;
      wsRef.current?.close();
    };
  }, []);

  return <WebSocketContext value={agents}>{children}</WebSocketContext>;
}

export function useAgents(): AgentMap {
  const ctx = use(WebSocketContext);
  if (ctx === null) throw new Error('useAgents must be used within WebSocketProvider');
  return ctx;
}
