import { useEffect } from "react";
import { useParams } from "react-router-dom";
import { getSnapshots } from "../../api/snapshots";
import DiskGauge from "../../components/charts/gauge";
import { convertByteToMetric, formatByteStr } from "../../utils";
import Line from "../../components/charts/line";
import { useAgents } from "../../store/agentsStore";

export default function AgentDetails() {
  const params = useParams();
  const agentId = params.agentId;

  const agent = useAgents((state) => state.getAgentState(agentId!));

  useEffect(() => {
    if (!agentId) return;
    getSnapshots(agentId).then((res) => {
      useAgents.getState().insertHistory(agentId, res.snapshots);
    });
  }, [agentId]);

  const currentDiskUsage =
    agent?.latest_snapshot?.data.disk_space.usage_percent ?? 0;
  const currentDiskUsed = agent?.latest_snapshot?.data.disk_space.used ?? 0;
  const currentDiskTotal = agent?.latest_snapshot?.data.disk_space.total ?? 0;
  const diskUsedHistory = agent?.history
    .map((snap) => convertByteToMetric(snap.data.disk_space.used, "GB"))
    .filter((du) => du !== undefined);
  const diskUsedTimestamps = agent?.history.map((snap) => snap.timestamp);

  const cpuAvgHistory = agent?.history
    .map((snap) => snap.data.cpu_usage.cpu_avg)
    .filter((cu) => cu !== undefined);
  const cpuAvgTimestamps = agent?.history.map((snap) => snap.timestamp);

  const cpusHistory = agent?.history[0].data.cpu_usage.cpu_info.map(
    (_, index) =>
      agent.history
        .map((snap) => snap.data.cpu_usage.cpu_info[index].usage_percent)
        .filter((cu) => cu != undefined),
  );
  return (
    <div className="px-6 py-6 flex-1 overflow-y-auto">
      <h2 className="text-lg font-medium text-(--text) mb-4">{agent?.name}</h2>

      <div className="mb-6 flex flex-row justify-start gap-6">
        <div className="w-64">
          <DiskGauge
            value={currentDiskUsage}
            label="Disk Usage"
            total={formatByteStr(currentDiskTotal)}
            used={formatByteStr(currentDiskUsed)}
          />
        </div>

        <div className="w-120">
          <Line
            values={diskUsedHistory ?? []}
            valueFormatter={(value: number) => `${value.toFixed(2)} GB`}
            timestamps={diskUsedTimestamps ?? []}
            label="Disk Used (GB)"
          />
        </div>

        <div className="w-120">
          <Line
            values={cpuAvgHistory ?? []}
            valueFormatter={(value: number) => `${value.toFixed(2)} %`}
            secondaryValues={cpusHistory}
            timestamps={cpuAvgTimestamps ?? []}
            label="CPU Average Usage (%)"
            tooltipItemPrefix="CPU"
          />
        </div>
      </div>
    </div>
  );
}
