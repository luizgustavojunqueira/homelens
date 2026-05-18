import { Icon } from '@iconify/react';
import { useNavigate } from 'react-router-dom';
import type { Agent } from '../../../api/models';

function level(pct: number, warn = 70, crit = 90): '' | 'warn' | 'crit' {
  if (pct >= crit) return 'crit';
  if (pct >= warn) return 'warn';
  return '';
}

function tempLevel(t: number): '' | 'warn' | 'crit' {
  if (t >= 80) return 'crit';
  if (t >= 65) return 'warn';
  return '';
}

export function AgentCard({ id, name, online, latest_snapshot }: Agent) {
  const navigate = useNavigate();
  const snap = latest_snapshot?.data;
  const cpuPct = snap?.cpu_usage?.cpu_avg ?? 0;
  const memUsed = (snap?.memory?.used ?? 0) * 1024;
  const memTotal = (snap?.memory?.total ?? 0) * 1024;
  const memPct = memTotal > 0 ? (memUsed / memTotal) * 100 : 0;
  const diskPct = snap?.disk_space?.usage_percent ?? 0;
  const temp = snap?.temperature?.temp_avg ?? 0;

  return (
    <button
      className="w-full grid grid-cols-[minmax(180px,1.4fr)_repeat(4,minmax(0,1fr))_24px] items-center gap-6 px-5 py-4 text-left hover:bg-(--bg-hover) transition-colors cursor-pointer"
      onClick={() => navigate(`/agents/${id}`)}
    >
      <div className="flex items-center gap-3 min-w-0">
        <span className={`dot ${online ? '' : 'off'}`}></span>
        <span className="text-base font-medium text-(--text) truncate">{name}</span>
      </div>

      {snap ? (
        <>
          <div className="flex items-center gap-2 min-w-0">
            <span className="text-sm text-(--text-dim) w-10">CPU</span>
            <div className={`bar ${level(cpuPct)} flex-1`}>
              <span style={{ width: `${Math.min(100, cpuPct)}%` }}></span>
            </div>
            <span className="tnum text-base text-(--text) w-12 text-right">
              {cpuPct.toFixed(0)}%
            </span>
          </div>

          <div className="flex items-center gap-2 min-w-0">
            <span className="text-sm text-(--text-dim) w-10">MEM</span>
            <div className={`bar ${level(memPct)} flex-1`}>
              <span style={{ width: `${Math.min(100, memPct)}%` }}></span>
            </div>
            <span className="tnum text-base text-(--text) w-12 text-right">
              {memPct.toFixed(0)}%
            </span>
          </div>

          <div className="flex items-center gap-2 min-w-0">
            <span className="text-sm text-(--text-dim) w-10">DSK</span>
            <div className={`bar ${level(diskPct, 80, 95)} flex-1`}>
              <span style={{ width: `${Math.min(100, diskPct)}%` }}></span>
            </div>
            <span className="tnum text-base text-(--text) w-12 text-right">
              {diskPct.toFixed(0)}%
            </span>
          </div>

          <div className="flex items-center gap-2 min-w-0">
            <span className="text-sm text-(--text-dim) w-10">TMP</span>
            <div className={`bar ${tempLevel(temp)} flex-1`}>
              <span style={{ width: `${Math.min(100, temp)}%` }}></span>
            </div>
            <span className="tnum text-base text-(--text) w-12 text-right">
              {temp.toFixed(0)}°
            </span>
          </div>
        </>
      ) : (
        <div className="col-span-4 text-sm text-(--text-faint)">no data</div>
      )}

      <Icon icon="lucide:chevron-right" width="16" className="text-(--text-faint)" />
    </button>
  );
}
