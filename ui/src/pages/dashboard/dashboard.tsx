import { useAgents } from '../../context/websocket';
import { AgentCard } from './components/agentCard';

export function Dashboard() {
  const agentMap = useAgents();
  const agents = Object.values(agentMap);
  const onlineCount = agents.filter((a) => a.online).length;

  return (
    <section className="px-6 py-6 flex-1 overflow-y-auto">
      <div className="flex items-baseline gap-3 mb-4">
        <h2 className="text-lg font-medium text-(--text)">Fleet</h2>
        <span className="tnum text-base text-(--text-dim)">
          {onlineCount}/{agents.length} online
        </span>
      </div>

      <div className="border border-(--border) rounded-md bg-(--bg-elev) divide-y divide-(--border)">
        {agents.map((agent) => (
          <AgentCard key={agent.id} {...agent} />
        ))}
      </div>
    </section>
  );
}
