import { ChevronRight } from "lucide-react";
import { useNavigate } from "react-router-dom";
import Tooltip from "../../../components/tooltip";
import MetricBar from "../../../components/metricBar";
import { formatByteStr } from "../../../utils";
import NetworkUsage from "../../../components/networkUsage";
import { useAgents } from "../../../store/agentsStore";

export function AgentCard({ guid }: { guid: string }) {
  const navigate = useNavigate();
  const agent = useAgents((state) => state.agents[guid]);

  if (!agent) return null;

  const { name, online, latest_snapshot } = agent;
  const snap = latest_snapshot?.data;

  let cpuPct = 0;
  let memPct = 0;
  let diskPct = 0;
  let temp = 0;
  let totalNetRx = 0;
  let totalNetTx = 0;

  if (snap) {
    cpuPct = snap.cpu.length > 0
      ? snap.cpu.reduce((cum, curr) => cum + curr.usage_percent, 0) / snap.cpu.length
      : 0;
    const memUsed = snap.memory?.used ?? 0;
    const memTotal = snap.memory?.total ?? 0;
    memPct = memTotal > 0 ? (memUsed / memTotal) * 100 : 0;
    diskPct = snap.disk?.disk_space?.usage_percent || 0;
    temp = snap.temperature && snap.temperature.length > 0
      ? snap.temperature.reduce((cum, curr) => cum + curr.temp_c, 0) / snap.temperature.length
      : 0;
    totalNetRx = snap.network?.reduce((sum, net) => sum + net.rx_bps, 0) || 0;
    totalNetTx = snap.network?.reduce((sum, net) => sum + net.tx_bps, 0) || 0;
  }

  return (
    <div
      className="w-full grid grid-cols-[minmax(150px,1fr)_repeat(4,1fr)_minmax(250px,1.4fr)_24px] items-center gap-6 px-5 py-0 text-left hover:bg-(--bg-hover) transition-colors cursor-pointer min-h-16 overflow-hidden overflow-x-auto"
      onClick={() => navigate(`/agents/${guid}`)}
    >
      <div className="flex items-center gap-3 min-w-0">
        <span className={`dot ${online ? "" : "off"}`}></span>
        <div className="flex flex-col min-w-0">
          <span className="text-base font-medium text-(--text) truncate">
            {name || snap?.host?.hostname || "Agent"}
          </span>
          {snap?.host?.os && (
            <span className="text-xs text-(--text-dim) truncate">
              {snap.host.os} {snap.host.platform ? `(${snap.host.platform})` : ""}
            </span>
          )}
        </div>
      </div>

      {snap ? (
        <>
          <Tooltip
            content={
              <div className="text-left min-w-60">
                <div>
                  <strong>Detailed CPU Usage</strong>
                </div>
                {snap.cpu.map(({ name, usage_percent }) => (
                  <MetricBar
                    key={name}
                    name={name}
                    value={usage_percent}
                    labelWidth="w-14"
                  />
                ))}
              </div>
            }
          >
            <MetricBar name={"CPU"} value={cpuPct} />
          </Tooltip>

          <Tooltip
            content={
              <div className="text-center">
                <div>
                  <strong>Detailed MEM Usage</strong>
                </div>
                <span>
                  {formatByteStr(snap.memory?.used ?? 0)} /{" "}
                  {formatByteStr(snap.memory?.total ?? 0)}
                </span>
              </div>
            }
          >
            <MetricBar name={"MEM"} value={memPct} />
          </Tooltip>

          <Tooltip
            content={
              <div className="text-center">
                <div>
                  <strong>Detailed DISK Usage</strong>
                </div>
                <span>
                  {formatByteStr(snap.disk?.disk_space?.used ?? 0)} /{" "}
                  {formatByteStr(snap.disk?.disk_space?.total ?? 0)}
                </span>
                <hr className="my-2 border-(--border)" />
                <div className="text-left">
                  {snap.disk?.disk_io_usage?.map(
                    ({ name, read_mbps, write_mbps }, index) => (
                      <div
                        key={`${name}-${index}`}
                        className="flex flex-col text-sm"
                      >
                        <span>{name}</span>
                        <span className="text-(--text-dim)">{`R: ${read_mbps.toFixed(2)} MB/s  W: ${write_mbps.toFixed(2)} MB/s`}</span>
                      </div>
                    ),
                  )}
                </div>
              </div>
            }
          >
            <MetricBar name={"DSK"} value={diskPct} />
          </Tooltip>

          <Tooltip
            content={
              <div className="text-center min-w-60">
                <div>
                  <strong>Detailed TEMP Usage</strong>
                </div>
                {snap.temperature &&
                  snap.temperature.map(({ temp_c, zone }, index) => (
                    <MetricBar
                      key={`${zone}-${index}`}
                      name={zone}
                      value={temp_c}
                      isTemp
                      labelWidth="w-28"
                    />
                  ))}
              </div>
            }
          >
            <MetricBar name={"TMP"} value={temp} isTemp />
          </Tooltip>

          <Tooltip
            content={
              <div className="text-center min-w-60">
                <div>
                  <strong>Detailed NET Usage</strong>
                </div>
                {snap.network?.map(({ name, rx_bps, tx_bps }, index) => (
                  <NetworkUsage
                    key={`${name}-${index}`}
                    name={name}
                    rx={rx_bps}
                    tx={tx_bps}
                    labelWidth="w-20"
                  />
                ))}
              </div>
            }
          >
            <NetworkUsage name="NET" rx={totalNetRx} tx={totalNetTx} />
          </Tooltip>
        </>
      ) : (
        <div className="col-span-4 text-sm text-(--text-faint)">no data</div>
      )}

      <ChevronRight
        width="16"
        className="text-(--text-faint)"
      />
    </div>
  );
}
