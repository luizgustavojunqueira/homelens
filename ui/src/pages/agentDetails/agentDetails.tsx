import { useEffect } from "react";
import { useParams } from "react-router-dom";
import { getSnapshots } from "../../api/snapshots";
import { useAgents } from "../../store/agentsStore";
import { useForm } from "react-hook-form";
import TextInput from "../../components/inputs/TextInput";
import { updateName } from "../../api/agents";
import { toast } from "react-toastify";

import AgentGauges from "./components/AgentGauges";
import AgentHistoryCharts from "./components/AgentHistoryCharts";
import DockerContainersTable from "./components/DockerContainersTable";
import ProcessesTable from "./components/ProcessesTable";

interface AgentForm {
  agentName: string;
}

function formatUptime(seconds?: number): string {
  if (!seconds || seconds <= 0) return "N/A";
  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const mins = Math.floor((seconds % 3600) / 60);

  const parts = [];
  if (days > 0) parts.push(`${days}d`);
  if (hours > 0) parts.push(`${hours}h`);
  if (mins > 0 || parts.length === 0) parts.push(`${mins}m`);

  return parts.join(" ");
}

export default function AgentDetails() {
  const params = useParams();
  const agentGuid = params.guid || "";

  const agent = useAgents((state) => state.agents[agentGuid]);

  const { control, reset } = useForm<AgentForm>({
    defaultValues: { agentName: agent?.name || "" },
  });

  useEffect(() => {
    if (agent?.name) {
      reset({ agentName: agent.name });
    }
  }, [agent?.name, reset]);

  useEffect(() => {
    if (!agentGuid) return;
    getSnapshots(agentGuid)
      .then((res) => {
        useAgents.getState().insertHistory(agentGuid, res.snapshots);
      })
      .catch(() => toast.error("Failed to load historical data"));
  }, [agentGuid]);

  if (!agent || !agent.latest_snapshot?.data) {
    return (
      <div className="flex items-center justify-center h-full w-full">
        <span className="text-(--text-dim)">Loading...</span>
      </div>
    );
  }

  const handleNameUpdate = (newName: string) => {
    if (newName.trim() === "" || newName === agent.name) return;

    updateName({ name: newName, guid: agent.guid })
      .then((res) => {
        if (res) {
          toast.success("Agent name changed");
        } else {
          toast.error("Error changing agent name");
        }
      })
      .catch(() => toast.error("Error changing agent name"));
  };

  const latestData = agent.latest_snapshot.data;
  const agentIp = latestData.agent_ip;
  const host = latestData.host;

  return (
    <section className="px-6 py-6 flex-1 overflow-y-auto max-w-screen">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-6">
        <div className="w-full max-w-sm">
          <TextInput
            name="agentName"
            control={control}
            onDebounce={handleNameUpdate}
            debounceTime={800}
            placeholder="Nome do Agente"
            className="text-xl font-medium w-full"
          />
        </div>

        <div className="flex flex-wrap items-center gap-2 text-xs text-(--text-dim)">
          {host?.hostname && (
            <span className="px-2.5 py-1 rounded-md bg-(--bg-elev) border border-(--border)">
              <strong>Host:</strong> {host.hostname}
            </span>
          )}
          {host?.os && (
            <span className="px-2.5 py-1 rounded-md bg-(--bg-elev) border border-(--border)">
              <strong>OS:</strong> {host.os} {host.platform ? `(${host.platform})` : ""}
            </span>
          )}
          {host?.kernel_version && (
            <span className="px-2.5 py-1 rounded-md bg-(--bg-elev) border border-(--border)">
              <strong>Kernel:</strong> {host.kernel_version}
            </span>
          )}
          {host?.uptime !== undefined && host.uptime > 0 && (
            <span className="px-2.5 py-1 rounded-md bg-(--bg-elev) border border-(--border)">
              <strong>Uptime:</strong> {formatUptime(host.uptime)}
            </span>
          )}
          {agentIp && (
            <span className="px-2.5 py-1 rounded-md bg-(--bg-elev) border border-(--border)">
              <strong>IP:</strong> {agentIp}
            </span>
          )}
        </div>
      </div>

      <div className="space-y-6">
        <AgentGauges latestData={latestData} />

        <div className="border border-(--border) rounded-md bg-(--bg-elev)"></div>

        <AgentHistoryCharts history={agent.history} />

        <DockerContainersTable
          containers={latestData.containers || []}
          agentIp={agentIp}
        />

        <ProcessesTable processes={latestData.processes || []} />
      </div>
    </section>
  );
}
