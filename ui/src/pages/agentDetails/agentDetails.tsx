import { useEffect } from "react";
import { useParams } from "react-router-dom";
import { getSnapshots } from "../../api/snapshots";
import Gauge from "../../components/charts/gauge";
import { convertByteToMetric, formatByteStr } from "../../utils";
import Line from "../../components/charts/line";
import { useAgents } from "../../store/agentsStore";
import { getMultiSeries, getSeries } from "./agentDetailsUtils";

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

  if (agent === undefined) return;

  const timestamps = agent.history.map((snap) => snap.timestamp) ?? [];

  const currentDiskUsage =
    agent.latest_snapshot.data.disk_space.usage_percent ?? 0;
  const currentDiskUsed = agent.latest_snapshot.data.disk_space.used ?? 0;
  const currentDiskTotal = agent.latest_snapshot.data.disk_space.total ?? 0;

  const diskTotalIoHistory = getMultiSeries(agent.history, (snap) =>
    snap.data.disk_io_usage.map((io) => io.read_mbps + io.write_mbps),
  );

  const diskNames = agent.latest_snapshot.data.disk_io_usage.map(
    (disk) => disk.name,
  );

  const currentCpuAvgUsage = agent.latest_snapshot.data.cpu_usage.cpu_avg ?? 0;
  const cpusHistory = getMultiSeries(agent.history, (snap) =>
    snap.data.cpu_usage.cpu_info.map((cpu) => cpu.usage_percent),
  );

  const currentMemUsed = agent.latest_snapshot.data.memory.used ?? 0;
  const currentMemTotal = agent.latest_snapshot.data.memory.total ?? 0;
  const currentMemUsage = (currentMemUsed / currentMemTotal) * 100;
  const memUsedHistory = agent.history.map((snap) =>
    convertByteToMetric(snap.data.memory.used, "GB", "KB"),
  );

  const netRxHistory = getSeries(agent.history, (snap) =>
    convertByteToMetric(
      snap.data.net_usage.reduce((sum, net) => sum + net.rx_bps, 0),
      "MB",
    ),
  );

  const netTxHistory = getSeries(agent.history, (snap) =>
    convertByteToMetric(
      snap.data.net_usage.reduce((sum, net) => sum + net.tx_bps, 0),
      "MB",
    ),
  );

  return (
    <section className="px-6 py-6 flex-1 overflow-y-auto">
      <div className="flex items-baseline gap-3 mb-6">
        <h2 className="text-xl font-medium text-(--text)">{agent.name}</h2>
      </div>

      <div className="space-y-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
          <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
            <Gauge
              value={currentDiskUsage}
              label="Disk Usage"
              total={formatByteStr(currentDiskTotal)}
              used={formatByteStr(currentDiskUsed)}
            />
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
            <Gauge
              value={currentCpuAvgUsage}
              label="CPU Usage"
              total="100 %"
              used={`${currentCpuAvgUsage.toFixed(2)} %`}
            />
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
            <Gauge
              value={currentMemUsage}
              label="RAM Usage"
              total={formatByteStr(currentMemTotal, "KB")}
              used={formatByteStr(currentMemUsed, "KB")}
            />
          </div>
        </div>

        <div className="border border-(--border) rounded-md bg-(--bg-elev)">
          <div className="px-4 py-3 border-b border-(--border)">
            <h3 className="text-sm font-medium text-(--text)">
              CPU Usage History
            </h3>
          </div>

          <div className="h-96 p-2">
            <Line
              isTotalAverage={true}
              timestamps={timestamps}
              label="CPU Usage (%)"
              valueFormatter={(v) => `${v.toFixed(2)} %`}
              series={[
                ...cpusHistory.map((cpu, i) => ({
                  name: `CPU ${i}`,
                  values: cpu,
                  subtle: true,
                })),
              ]}
            />
          </div>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">
                Memory Usage History
              </h3>
            </div>

            <div className="h-80 p-2">
              <Line
                valueFormatter={(value) => formatByteStr(value, "GB")}
                timestamps={timestamps}
                label="RAM Used (GB)"
                series={[
                  {
                    values: memUsedHistory,
                    name: "RAM",
                  },
                ]}
              />
            </div>
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">
                Disk Usage History
              </h3>
            </div>

            <div className="h-80 p-2">
              <Line
                timestamps={timestamps}
                label="Disk IO (MB/s)"
                valueFormatter={(v) => `${v.toFixed(2)} MB/s`}
                series={diskTotalIoHistory.map((disk, i) => ({
                  name: diskNames[i],
                  values: disk,
                }))}
              />
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">
                Network Throughput
              </h3>
            </div>

            <div className="h-80 p-2">
              <Line
                timestamps={timestamps}
                label="Network Throughput"
                valueFormatter={(v) => `${v.toFixed(2)} MB/s`}
                series={[
                  {
                    name: "RX",
                    values: netRxHistory,
                  },
                  {
                    name: "TX",
                    values: netTxHistory,
                  },
                ]}
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
