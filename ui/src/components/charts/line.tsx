import type {
  EChartsOption,
  SeriesOption,
  TooltipComponentFormatterCallbackParams,
} from "echarts";
import ReactECharts from "echarts-for-react";
import type { CallbackDataParams } from "echarts/types/dist/shared";

const MAX_POINTS = 500;

export interface ISeries {
  name: string;
  values: number[];
  subtle?: boolean;
}

interface ILine {
  timestamps: number[];
  series: ISeries[];
  valueFormatter?: (value: number) => string;
  label: string;
  isTotalAverage?: boolean;
}

const color = (index: number, total: number) =>
  `hsl(${(index * 360) / Math.max(total, 1)}, 65%, 55%)`;

function getPointValue(param: CallbackDataParams): number | null {
  if (!Array.isArray(param.value)) {
    return null;
  }

  const value = param.value[1];

  return typeof value === "number" ? value : null;
}

export default function Line({
  timestamps,
  series,
  valueFormatter = (value: number) => value.toString(),
  label,
  isTotalAverage = false,
}: ILine) {
  const start = Math.max(0, timestamps.length - MAX_POINTS);

  const recentTimestamps = timestamps.slice(start);

  const chartSeries = series.map((serie) => ({
    ...serie,
    data: recentTimestamps.map((ts, i) => [
      ts,
      serie.values[start + i] ?? null,
    ]),
  }));

  const seriesColors = new Map<string, string>();

  chartSeries.forEach((serie, index) => {
    const serieColor =
      index === 0
        ? "#3b82f6"
        : color(index - 1, Math.max(chartSeries.length - 1, 1));

    seriesColors.set(serie.name, serieColor);
  });

  const option: EChartsOption = {
    backgroundColor: "transparent",

    tooltip: {
      trigger: "axis",

      position: (point, _params, _dom, _rect, size) => {
        const [x, y] = point;

        const tooltipWidth = size.contentSize[0];
        const viewWidth = size.viewSize[0];

        let left = x + 20;

        if (left + tooltipWidth > viewWidth) {
          left = x - tooltipWidth - 20;
        }

        return [left, y];
      },

      backgroundColor: "rgba(24,24,27,0.96)",
      borderColor: "#3f3f46",
      borderWidth: 1,

      textStyle: {
        color: "#fafafa",
        fontSize: 12,
      },

      padding: 8,

      axisPointer: {
        type: "cross",
        animation: false,

        lineStyle: {
          color: "#374151",
        },
      },

      formatter: (params: TooltipComponentFormatterCallbackParams) => {
        if (!Array.isArray(params) || params.length === 0) {
          return "";
        }

        const first = params[0];

        const axisValue =
          Array.isArray(first.value) && first.value.length > 0
            ? Number(first.value[0])
            : Date.now();

        const date = new Date(axisValue);

        const visibleSeries = params
          .map((p) => ({
            param: p,
            value: getPointValue(p),
          }))
          .filter(
            (
              item,
            ): item is {
              param: CallbackDataParams;
              value: number;
            } => item.value !== null,
          );

        const ordered = [...visibleSeries].sort((a, b) => b.value - a.value);

        const topSeries = ordered.slice(0, 8);
        const remaining = ordered.length - topSeries.length;

        const summary =
          ordered.length > 1
            ? isTotalAverage
              ? ordered.reduce((sum, item) => sum + item.value, 0) /
                ordered.length
              : ordered.reduce((sum, item) => sum + item.value, 0)
            : null;

        const summaryLabel = isTotalAverage ? "Average" : "Total";

        return `
<div style="font-weight:600;margin-bottom:6px">
  ${label}
</div>

<div style="opacity:0.75;margin-bottom:8px">
  ${date.toLocaleString()}
</div>

${
  summary !== null
    ? `
<div
  style="
    display:flex;
    justify-content:space-between;
    gap:16px;
    margin-bottom:${topSeries.length > 0 ? "8px" : "0"};
    font-weight:600;
    color:#93c5fd;
  "
>
  <span>${summaryLabel}</span>
  <span>${valueFormatter(summary)}</span>
</div>
`
    : ""
}

${
  topSeries.length > 0
    ? `
<div
  style="
    border-top:1px solid #374151;
    margin-top:6px;
    padding-top:6px;
  "
>
  ${topSeries
    .map(({ param, value }) => {
      const serieColor = seriesColors.get(param.seriesName ?? "") ?? "#ffffff";

      return `
<div
  style="
    display:flex;
    justify-content:space-between;
    gap:16px;
    margin-bottom:2px;
  "
>
<span style="color:${serieColor}">
  <span
    style="
      display:inline-block;
      width:10px;
      height:10px;
      border-radius:50%;
      background:${serieColor};
      margin-right:6px;
      vertical-align:middle;
    "
  ></span>
  ${param.seriesName}
</span>

  <span>
    ${valueFormatter(value)}
  </span>
</div>
`;
    })
    .join("")}

  ${
    remaining > 0
      ? `
<div
  style="
    margin-top:4px;
    opacity:0.65;
    font-style:italic;
  "
>
  +${remaining} more series
</div>
`
      : ""
  }
</div>
`
    : ""
}
`;
      },
    },

    grid: {
      left: 50,
      right: 20,
      top: 20,
      bottom: 50,
    },

    xAxis: {
      type: "time",

      axisLine: {
        lineStyle: {
          color: "#3f3f46",
        },
      },

      axisLabel: {
        color: "#a1a1aa",
      },
    },

    yAxis: {
      type: "value",

      axisLine: {
        show: false,
      },

      splitLine: {
        lineStyle: {
          color: "rgba(255,255,255,0.06)",
        },
      },

      axisLabel: {
        color: "#a1a1aa",
      },
    },

    series: chartSeries.map(
      (serie): SeriesOption => ({
        name: serie.name,

        type: "line",

        data: serie.data,

        connectNulls: false,

        sampling: "lttb",
        progressive: 5000,
        progressiveThreshold: 10000,

        animation: false,
        showSymbol: false,

        smooth: true,

        silent: !!serie.subtle,

        lineStyle: {
          width: serie.subtle ? 1 : 3,
          opacity: serie.subtle ? 0.15 : 1,
          color: seriesColors.get(serie.name),
        },

        areaStyle: serie.subtle
          ? undefined
          : {
              opacity: 0.05,
              color: seriesColors.get(serie.name),
            },

        z: serie.subtle ? 1 : 100,
      }),
    ),

    dataZoom: [
      {
        type: "inside",
      },
      {
        type: "slider",

        height: 10,
        bottom: 10,

        borderColor: "transparent",
        backgroundColor: "rgba(255,255,255,0.04)",
        fillerColor: "rgba(59,130,246,0.20)",

        handleStyle: {
          color: "#60a5fa",
          borderColor: "#60a5fa",
          borderWidth: 1,
        },

        moveHandleStyle: {
          color: "#1f2937",
        },

        textStyle: {
          color: "#94a3b8",
        },

        showDetail: false,
      },
    ],
  };

  return (
    <ReactECharts
      option={option}
      lazyUpdate
      style={{ width: "100%", height: "100%" }}
    />
  );
}
