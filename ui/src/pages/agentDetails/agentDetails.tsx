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

  return (
    <section className="px-6 py-6 flex-1 overflow-y-auto max-w-screen">
      <div className="flex items-baseline gap-3 mb-6 w-full max-w-sm">
        <TextInput
          name="agentName"
          control={control}
          onDebounce={handleNameUpdate}
          debounceTime={800}
          placeholder="Nome do Agente"
          className="text-xl font-medium w-full"
        />
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
