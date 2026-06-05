import type { EChartsOption, SeriesOption } from "echarts";
import ReactECharts from "echarts-for-react";

interface ILine {
  timestamps: number[];
  values: number[];
  secondaryValues?: number[][];
  valueFormatter?: (value: number) => string;
  label: string;
  tooltipItemPrefix?: string;
}

const color = (index: number, total: number) =>
  `hsl(${(index * 360) / Math.max(total, 1)}, 70%, 60%)`;

export default function Line({
  timestamps,
  values,
  secondaryValues,
  valueFormatter = (value: number) => value.toString(),
  label,
  tooltipItemPrefix = "Series",
}: ILine) {
  const data =
    timestamps?.map((ts, i) => {
      const date = new Date(ts);

      return [date.getTime(), values[i]];
    }) ?? values.map((v, i) => [i, v]);

  const secondaryData =
    secondaryValues?.map(
      (values) =>
        timestamps?.map((ts, i) => {
          const date = new Date(ts);

          return [date.getTime(), values[i]];
        }) ?? values.map((v, i) => [i, v]),
    ) ?? [];

  const option: EChartsOption = {
    backgroundColor: "#0b1220",

    title: {
      text: label,
      left: 10,
      top: 10,
      textStyle: {
        color: "#e5e7eb",
        fontSize: 14,
      },
    },

    tooltip: {
      trigger: "axis",
      position: (
        point: number[],
        _params: any,
        _dom: HTMLElement,
        _rect: any,
        size: any,
      ) => {
        const [x, y] = point;

        const tooltipWidth = size.contentSize[0];
        const viewWidth = size.viewSize[0];

        let left = x + 20;

        if (left + tooltipWidth > viewWidth) {
          left = x - tooltipWidth - 20;
        }

        return [left, y];
      },

      backgroundColor: "#111827",
      borderColor: "#374151",
      borderWidth: 1,

      textStyle: {
        color: "#e5e7eb",
        fontSize: 12,
      },

      padding: 8,

      axisPointer: {
        type: "cross",
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
          color: "#374151",
        },
      },

      axisLabel: {
        color: "#94a3b8",
      },
    },

    yAxis: {
      type: "value",

      axisLine: {
        show: false,
      },

      splitLine: {
        lineStyle: {
          color: "#1f2937",
        },
      },

      axisLabel: {
        color: "#94a3b8",
      },
    },

    series: [
      {
        name: "Average",

        type: "line",

        data,

        smooth: true,

        showSymbol: false,

        lineStyle: {
          width: 4,
          color: "#2563eb",
        },

        z: 100,

        areaStyle: {
          color: "rgba(37,99,235,0.15)",
        },
      },

      ...secondaryData.map(
        (data, index): SeriesOption => ({
          name: `${tooltipItemPrefix} ${index + 1}`,

          type: "line",

          data,

          smooth: true,

          showSymbol: false,

          lineStyle: {
            width: 1,
            opacity: 0.25,
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
        backgroundColor: "#0f172a",
        fillerColor: "rgba(96,165,250,0.25)",

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
