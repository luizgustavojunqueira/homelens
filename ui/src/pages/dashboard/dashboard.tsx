import { useAgents } from "../../store/agentsStore";
import { AgentCard } from "./components/agentCard";
import { useMemo } from "react";

export function Dashboard() {
  // Subscribe only to the keys so Dashboard doesn't re-render on every agent update
  const agentGuids = useAgents((state) => Object.keys(state.agents));
  
  // To show online count without passing down data, we can just fetch it inside
  const agents = useAgents((state) => state.agents);
  const onlineCount = useMemo(() => Object.values(agents).filter((a) => a.online).length, [agents]);

  return (
    <section className="px-6 py-6 flex-1 overflow-y-auto">
      <div className="flex items-baseline gap-3 mb-4">
        <h2 className="text-lg font-medium text-(--text)">Fleet</h2>
        <span className="tnum text-base text-(--text-dim)">
          {onlineCount}/{agentGuids.length} online
        </span>
      </div>

      <div className="border border-(--border) rounded-md bg-(--bg-elev) divide-y divide-(--border)">
        {agentGuids.map((guid) => (
          <AgentCard key={guid} guid={guid} />
        ))}
      </div>
    </section>
  );
}
