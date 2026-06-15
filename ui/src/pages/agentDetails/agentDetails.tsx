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
  const agentGuid = params.guid;

  const agent = useAgents((state) => state.getAgentState(agentGuid!));

  useEffect(() => {
    if (!agentGuid) return;
    getSnapshots(agentGuid).then((res) => {
      useAgents.getState().insertHistory(agentGuid, res.snapshots);
    });
  }, [agentGuid]);

  if (agent === undefined) return;

  const timestamps = agent.history.map((snap) => snap.timestamp) ?? [];

  const currentDiskUsage =
    agent.latest_snapshot.data.disk.disk_space.usage_percent;
  const currentDiskUsed = agent.latest_snapshot.data.disk.disk_space.used;
  const currentDiskTotal = agent.latest_snapshot.data.disk.disk_space.total;

  const diskTotalIoHistory = getMultiSeries(agent.history, (snap) =>
    snap.data.disk.disk_io_usage.map((io) => io.read_mbps + io.write_mbps),
  );

  const diskNames = agent.latest_snapshot.data.disk.disk_io_usage.map(
    (disk) => disk.name,
  );

  const currentCpuAvgUsage =
    agent.latest_snapshot.data.cpu.reduce(
      (cum, curr) => cum + curr.usage_percent,
      0,
    ) / agent.latest_snapshot.data.cpu.length;
  const cpusHistory = getMultiSeries(agent.history, (snap) =>
    snap.data.cpu.map((cpu) => cpu.usage_percent),
  );
  const cpuAvgHistory = getSeries(
    agent.history,
    (snap) =>
      snap.data.cpu.reduce((cum, curr) => cum + curr.usage_percent, 0) /
      snap.data.cpu.length,
  );

  const currentMemUsed = agent.latest_snapshot.data.memory.used ?? 0;
  const currentMemTotal = agent.latest_snapshot.data.memory.total ?? 0;
  const currentMemUsage = (currentMemUsed / currentMemTotal) * 100;
  const currentTemp =
    agent.latest_snapshot.data.temperature.reduce(
      (cum, curr) => cum + curr.temp_c,
      0,
    ) / agent.latest_snapshot.data.temperature.length;
  const memUsedHistory = agent.history.map((snap) =>
    convertByteToMetric(snap.data.memory.used, "GB", "KB"),
  );

  const netRxHistory = getSeries(agent.history, (snap) =>
    convertByteToMetric(
      snap.data.network.reduce((sum, net) => sum + net.rx_bps, 0),
      "MB",
    ),
  );

  const netTxHistory = getSeries(agent.history, (snap) =>
    convertByteToMetric(
      snap.data.network.reduce((sum, net) => sum + net.tx_bps, 0),
      "MB",
    ),
  );

  return (
    <section className="px-6 py-6 flex-1 overflow-y-auto">
      <div className="flex items-baseline gap-3 mb-6">
        <h2 className="text-xl font-medium text-(--text)">{agent.name}</h2>
      </div>

      <div className="space-y-6">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-4">
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

          <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
            <Gauge
              value={currentTemp}
              label="Temperature"
              symbol="C°"
              total={"100 C°"}
              used={`${currentTemp.toFixed(2)} C°`}
            />
          </div>
        </div>

        <div className="border border-(--border) rounded-md bg-(--bg-elev)"></div>

        <div className="grid grid-cols-1 xl:grid-cols-2 grid-rows-2 gap-4">
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
                  {
                    name: "Average",
                    values: cpuAvgHistory,
                  },
                  ...cpusHistory.map((cpu, i) => ({
                    name: `CPU ${i}`,
                    values: cpu,
                    visible: false,
                  })),
                ]}
              />
            </div>
          </div>

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

        <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">
                Docker Containers
              </h3>
            </div>
          </div>

          <div className="border border-(--border) rounded-md bg-(--bg-elev)">
            <div className="px-4 py-3 border-b border-(--border)">
              <h3 className="text-sm font-medium text-(--text)">Processess</h3>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
