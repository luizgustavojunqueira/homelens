import type { EChartsOption, SeriesOption } from "echarts";
import ReactECharts from "echarts-for-react";

const MAX_POINTS = 500;

interface ILine {
  timestamps: number[];
  values: number[];
  secondaryValues?: number[][];
  valueFormatter?: (value: number) => string;
  label: string;
  tooltipItemPrefix?: string;
  seriesNames?: string[];
}

const color = (index: number, total: number) =>
  `hsl(${(index * 360) / Math.max(total, 1)}, 65%, 55%)`;

export default function Line({
  timestamps,
  values,
  secondaryValues,
  valueFormatter = (value: number) => value.toString(),
  label,
  tooltipItemPrefix = "Series",
  seriesNames,
}: ILine) {
  const start = Math.max(0, timestamps.length - MAX_POINTS);

  const recentTimestamps = timestamps.slice(start);
  const recentValues = values.slice(start);

  const data = recentTimestamps.map((ts, i) => [ts, recentValues[i]]);

  const secondaryData =
    secondaryValues?.map((series) =>
      recentTimestamps.map((ts, i) => [ts, series[start + i]]),
    ) ?? [];

  const option: EChartsOption = {
    backgroundColor: "transparent",

    title: {
      left: 10,
      top: 10,
      textStyle: {
        color: "#e5e7eb",
        fontSize: 14,
      },
    },

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

      formatter: (params: any) => {
        if (!params?.length) {
          return "";
        }

        const date = new Date(params[0].axisValue);

        const average = params.find((p: any) => p.seriesName === "Average");

        const secondary = params
          .filter((p: any) => p.seriesName !== "Average")
          .sort((a: any, b: any) => b.value[1] - a.value[1]);

        const topSecondary = secondary.slice(0, 8);

        const remaining = secondary.length - topSecondary.length;

        return `
          <div style="font-weight:600;margin-bottom:6px">
            ${label}
          </div>

          <div style="opacity:0.75;margin-bottom:8px">
            ${date.toLocaleString()}
          </div>

          ${
            average
              ? `
            <div
              style="
                display:flex;
                justify-content:space-between;
                gap:16px;
                margin-bottom:8px;
                font-weight:600;
                color:#93c5fd;
              "
            >
              <span>Average</span>
              <span>${valueFormatter(average.value[1])}</span>
            </div>
          `
              : ""
          }

          ${
            topSecondary.length > 0
              ? `
            <div
              style="
                border-top:1px solid #374151;
                margin-top:6px;
                padding-top:6px;
              "
            >
              ${topSecondary
                .map(
                  (item: any) => `
                  <div
                    style="
                      display:flex;
                      justify-content:space-between;
                      gap:16px;
                      margin-bottom:2px;
                    "
                  >
                    <span style="color:${item.color}">
                      ${item.marker} ${item.seriesName}
                    </span>
                    <span>${valueFormatter(item.value[1])}</span>
                  </div>
                `,
                )
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
    grid: {
      left: 50,
      right: 20,
      top: 20,
      bottom: 50,
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

    series: [
      {
        name: "Average",

        type: "line",

        data,

        lineStyle: {
          width: 4,
          color: "#3b82f6",
        },

        z: 100,

        areaStyle: {
          color: "rgba(59,130,246,0.12)",
        },

        sampling: "lttb",
        progressive: 5000,
        progressiveThreshold: 10000,
        animation: false,
        showSymbol: false,
        smooth: true,
      },

      ...secondaryData.map(
        (data, index): SeriesOption => ({
          name:
            seriesNames === undefined
              ? `${tooltipItemPrefix} ${index + 1}`
              : `${seriesNames[index]}`,

          type: "line",

          data,
          silent: true,
          sampling: "lttb",
          progressive: 5000,
          progressiveThreshold: 10000,
          animation: false,
          showSymbol: false,
          smooth: true,

          lineStyle: {
            width: 1,
            opacity: 0.15,
            color: color(index, secondaryData.length),
          },

          z: 1,

          areaStyle: undefined,
        }),
      ),
    ],

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
