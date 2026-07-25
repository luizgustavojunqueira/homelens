import { formatByteStr } from "../utils";

interface INetworkUsage {
  name: string;
  rx: number;
  tx: number;
  labelWidth?: string;
}

export default function NetworkUsage({ name, rx, tx, labelWidth = 'w-fit' }: INetworkUsage) {
  return (
    <div className="flex flex-row gap-2 justify-end items-center min-w-0 h-full whitespace-nowrap">
      <span className={`text-sm text-(--text-dim) truncate ${labelWidth}`}>{name}</span>
      <div className="flex items-center gap-1 min-w-0 whitespace-nowrap">
        <span className="text-green-500">↓</span>
        <span className="tnum text-base text-(--text) whitespace-nowrap text-right min-w-24">
          {formatByteStr(rx, "B")}/s
        </span>
      </div>
      <div className="flex items-center gap-1 min-w-0 whitespace-nowrap">
        <span className="text-red-500">↑</span>
        <span className="tnum text-base text-(--text) whitespace-nowrap text-right min-w-24">
          {formatByteStr(tx, "B")}/s
        </span>
      </div>
    </div>
  )
}
