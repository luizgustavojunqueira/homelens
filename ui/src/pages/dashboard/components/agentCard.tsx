import { Icon } from '@iconify/react';
import { useNavigate } from 'react-router-dom';
import type { Agent } from '../../../api/models';
import Tooltip from '../../../components/tooltip';
import MetricBar from '../../../components/metricBar';
import { formatByteStr } from '../../../utils';
import NetworkUsage from '../../../components/networkUsage';

export function AgentCard({ id, name, online, latest_snapshot }: Agent) {
  const navigate = useNavigate();
  const snap = latest_snapshot?.data;
  const cpuPct = snap?.cpu_usage?.cpu_avg ?? 0;
  const memUsed = (snap?.memory?.used ?? 0) * 1024;
  const memTotal = (snap?.memory?.total ?? 0) * 1024;
  const memPct = memTotal > 0 ? (memUsed / memTotal) * 100 : 0;
  const diskPct = snap?.disk_space?.usage_percent ?? 0;
  const temp = snap?.temperature?.temp_avg ?? 0;

  const totalNetRx = snap?.net_usage?.reduce((sum, net) => sum + net.rx_bps, 0) ?? 0;
  const totalNetTx = snap?.net_usage?.reduce((sum, net) => sum + net.tx_bps, 0) ?? 0;

  return (
    <div
      className="w-full grid grid-cols-[minmax(150px,1fr)_repeat(5,1.1fr)_24px] items-center gap-6 px-5 py-0 text-left hover:bg-(--bg-hover) transition-colors cursor-pointer min-h-16 overflow-hidden overflow-x-auto"
      onClick={() => navigate(`/agents/${id}`)}
    >
      <div className="flex items-center gap-3 min-w-0">
        <span className={`dot ${online ? '' : 'off'}`}></span>
        <span className="text-base font-medium text-(--text) truncate">{name}</span>
      </div>

      {snap ? (
        <>
          <Tooltip
            content={
              <div className="text-left min-w-60">
                <div><strong>Detailed CPU Usage</strong></div>
                {snap.cpu_usage.cpu_info.map(({ name, usage_percent }) => (
                  <MetricBar key={name} name={name} value={usage_percent} labelWidth="w-14" />
                ))
                }
              </div>
            }>
            <MetricBar name={"CPU"} value={cpuPct} />
          </Tooltip>


          <Tooltip
            content={
              <div className="text-center">
                <div><strong>Detailed MEM Usage</strong></div>
                <span>{formatByteStr(snap.memory.used, "KB")} /  {formatByteStr(snap.memory.total, "KB")}</span>
              </div>
            }>
            <MetricBar name={"MEM"} value={memPct} />
          </Tooltip>


          <Tooltip
            content={
              <div className="text-center">
                <div><strong>Detailed DISK Usage</strong></div>
                <span>{formatByteStr(snap.disk_space.used,)} /  {formatByteStr(snap.disk_space.total,)}</span>
                <hr className="my-2 border-(--border)" />
                <div className="text-left">
                  {snap.disk_io_usage.map(({ name, read_mbps, write_mbps }) => (
                    <div key={name} className="flex flex-col text-sm">
                      <span>{name}</span>
                      <span className="text-(--text-dim)">{`R: ${read_mbps.toFixed(2)} MB/s  W: ${write_mbps.toFixed(2)} MB/s`}</span>
                    </div>
                  ))}
                </div>
              </div>
            }>
            <MetricBar name={"DSK"} value={diskPct} />
          </Tooltip>


          <Tooltip
            content={
              <div className="text-center min-w-60">
                <div><strong>Detailed TEMP Usage</strong></div>
                {snap.temperature.temp_info.map(({ temp_c, zone }) => (
                  <MetricBar key={zone} name={zone} value={temp_c} isTemp labelWidth="w-28" />
                ))}

              </div>
            }>
            <MetricBar name={"TMP"} value={temp} isTemp />
          </Tooltip>


          <Tooltip
            content={
              <div className="text-center min-w-60">
                <div><strong>Detailed NET Usage</strong></div>
                {snap.net_usage.map(({ name, rx_bps, tx_bps }) => (
                  <NetworkUsage key={name} name={name} rx={rx_bps} tx={tx_bps} labelWidth="w-20" />
                ))}
              </div>
            }>
            <NetworkUsage name='NET' rx={totalNetRx} tx={totalNetTx} />
          </Tooltip>
        </>
      ) : (
        <div className="col-span-4 text-sm text-(--text-faint)">no data</div>
      )}

      <Icon icon="lucide:chevron-right" width="16" className="text-(--text-faint)" />
    </div>
  );
}
