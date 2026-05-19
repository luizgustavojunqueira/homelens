interface IMetricBar {
  name: string;
  value: number;
  warn?: number;
  crit?: number;
  isTemp?: boolean;
  labelWidth?: string;
}

function level(pct: number, warn: number, crit: number): '' | 'warn' | 'crit' {
  if (pct >= crit) return 'crit';
  if (pct >= warn) return 'warn';
  return '';
}

export default function MetricBar({ name, value, warn = 70, crit = 85, isTemp, labelWidth = 'w-fit' }: IMetricBar) {
  return (
    <div className="flex items-center gap-2 min-w-0 h-full">
      <span className={`text-sm text-(--text-dim) truncate ${labelWidth}`}>{name}</span>
      <div className={`bar ${level(value, warn, crit)} flex-1`}>
        <span style={{ width: `${Math.min(100, value)}%` }}></span>
      </div>
      <span className="tnum text-base text-(--text) w-12 text-right">
        {value.toFixed(0)}{isTemp ? '°C' : '%'}
      </span>
    </div>
  )
}
