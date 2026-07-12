import type { EChartsOption } from "echarts";
import ReactECharts from "echarts-for-react";
import { useMemo } from "react";

interface IGauge {
  value: number;
  symbol?: string;
  label: string;
  used?: string;
  total?: string;
}

const escapeHtml = (unsafe: string) => {
  return String(unsafe)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
};

export default function Gauge({
  value,
  symbol = "%",
  label,
  used,
  total,
}: IGauge) {
  const option: EChartsOption = useMemo(() => ({
    backgroundColor: "transparent",

    animationDuration: 800,

    tooltip: {
      backgroundColor: "#111827",
      borderColor: "#374151",
      borderWidth: 1,

      textStyle: {
        color: "#e5e7eb",
        fontSize: 12,
      },
      padding: 10,
      formatter: () => {
        if (used !== undefined && total !== undefined) {
          return `
            <div style="font-weight:600;margin-bottom:4px">
              ${escapeHtml(label)}
            </div>
            <div style="opacity:0.8">
              ${escapeHtml(String(used))} / ${escapeHtml(String(total))}
            </div>
          `;
        }
        return `
          <div style="font-weight:600">
            ${escapeHtml(label)}
          </div>
          <div>${value.toFixed(1)}${escapeHtml(symbol)}</div>
        `;
      },
      position: (point) => {
        return [point[0] - 60, point[1] + 20];
      },
    },

    series: [
      {
        type: "gauge",
        radius: "95%",

        startAngle: 225,
        endAngle: -45,

        min: 0,
        max: 100,

        pointer: {
          show: false,
        },

        progress: {
          show: true,
          roundCap: true,
          width: 18,
          itemStyle: {
            color: value > 90 ? "#ef4444" : value > 75 ? "#f59e0b" : "#3b82f6",
          },
        },

        axisLine: {
          lineStyle: {
            width: 18,
            color: [[1, "rgba(255,255,255,0.08)"]],
          },
        },

        splitLine: {
          show: false,
        },

        axisTick: {
          show: false,
        },

        axisLabel: {
          show: false,
        },

        detail: {
          valueAnimation: true,
          fontSize: 24,
          fontWeight: "bold",
          color: "#fafafa",
          offsetCenter: [0, "-5%"],
          formatter: (v: number) => `${v.toFixed(1)} ${symbol}`,
        },

        title: {
          offsetCenter: [0, "20%"],
          color: "#a1a1aa",
          fontSize: 14,
        },

        data: [
          {
            value,
            name: label,
          },
        ],
      },
    ],
  }), [value, symbol, label, used, total]);

  return (
    <ReactECharts
      option={option}
      lazyUpdate
      style={{ width: "100%", height: "100%" }}
    />
  );
}
