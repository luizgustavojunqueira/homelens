import Gauge from "../../../components/charts/gauge";
import { formatByteStr } from "../../../utils";
import type { SystemInfo } from "../../../api/models";

export default function AgentGauges({ latestData }: { latestData: SystemInfo }) {
  const currentDiskUsage = latestData.disk?.disk_space?.usage_percent || 0;
  const currentDiskUsed = latestData.disk?.disk_space?.used || 0;
  const currentDiskTotal = latestData.disk?.disk_space?.total || 0;

  const currentCpuAvgUsage = latestData.cpu && latestData.cpu.length > 0
    ? latestData.cpu.reduce((cum, curr) => cum + curr.usage_percent, 0) / latestData.cpu.length
    : 0;

  const currentMemUsed = latestData.memory?.used ?? 0;
  const currentMemTotal = latestData.memory?.total ?? 0;
  const currentMemUsage = currentMemTotal > 0 ? (currentMemUsed / currentMemTotal) * 100 : 0;
  
  const currentTemp = latestData.temperature && latestData.temperature.length > 0
    ? latestData.temperature.reduce((cum, curr) => cum + curr.temp_c, 0) / latestData.temperature.length
    : undefined;

  return (
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
          total={formatByteStr(currentMemTotal)}
          used={formatByteStr(currentMemUsed)}
        />
      </div>

      {currentTemp !== undefined && (
        <div className="border border-(--border) rounded-md bg-(--bg-elev) p-4 h-64">
          <Gauge
            value={currentTemp}
            label="Temperature"
            symbol="C°"
            total={"100 C°"}
            used={`${currentTemp.toFixed(2)} C°`}
          />
        </div>
      )}
    </div>
  );
}
