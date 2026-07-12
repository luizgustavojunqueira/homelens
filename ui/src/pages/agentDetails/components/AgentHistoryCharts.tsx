import { useMemo } from "react";
import Line from "../../../components/charts/line";
import type { SnapshotEntry } from "../../../api/models";
import { getMultiSeries, getSeries } from "../agentDetailsUtils";
import { convertByteToMetric, formatByteStr } from "../../../utils";

export default function AgentHistoryCharts({ history }: { history: SnapshotEntry[] }) {
  const timestamps = useMemo(() => history.map((snap) => snap.timestamp), [history]);

  const cpusHistory = useMemo(() => getMultiSeries(history, (snap) =>
    snap.data.cpu?.map((cpu) => cpu.usage_percent) || []
  ), [history]);
  
  const cpuAvgHistory = useMemo(() => getSeries(
    history,
    (snap) => snap.data.cpu && snap.data.cpu.length > 0
      ? snap.data.cpu.reduce((cum, curr) => cum + curr.usage_percent, 0) / snap.data.cpu.length
      : 0
  ), [history]);

  const memUsedHistory = useMemo(() => history.map((snap) =>
    convertByteToMetric(snap.data.memory?.used || 0, "GB", "KB")
  ), [history]);

  const diskTotalIoHistory = useMemo(() => getMultiSeries(history, (snap) =>
    snap.data.disk?.disk_io_usage?.map((io) => io.read_mbps + io.write_mbps) || []
  ), [history]);

  const diskNames = history[0]?.data.disk?.disk_io_usage?.map((disk) => disk.name) || [];

  const netRxHistory = useMemo(() => getSeries(history, (snap) =>
    convertByteToMetric(
      snap.data.network?.reduce((sum, net) => sum + net.rx_bps, 0) || 0,
      "MB"
    )
  ), [history]);

  const netTxHistory = useMemo(() => getSeries(history, (snap) =>
    convertByteToMetric(
      snap.data.network?.reduce((sum, net) => sum + net.tx_bps, 0) || 0,
      "MB"
    )
  ), [history]);

  return (
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
              name: diskNames[i] || `Disk ${i}`,
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
  );
}
