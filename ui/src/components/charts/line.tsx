import type { EChartsOption } from "echarts"
import ReactECharts from "echarts-for-react"

interface ILine {
  timestamps: number[]
  values: number[]
  valueFormatter?: (value: number) => string
  label: string
}

export default function Line({
  timestamps,
  values,
  valueFormatter = (value: number) => value.toString(),
  label,
}: ILine) {

  const data =
    timestamps?.map((ts, i) => {
      const date = new Date(ts)

      return [date.getTime(), values[i]]
    }) ?? values.map((v, i) => [i, v])

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

      backgroundColor: "#111827",
      borderColor: "#374151",
      borderWidth: 1,
      textStyle: {
        color: "#e5e7eb",
        fontSize: 12,
      },

      padding: 10,

      axisPointer: {
        type: "cross",
        lineStyle: {
          color: "#374151",
        },
      },

      formatter: (params: any) => {
        const p = params?.[0]

        const date = new Date(p.axisValue)

        return `
          <div style="font-weight:600;margin-bottom:4px">
            ${label}
          </div>

          <div style="opacity:0.8;margin-bottom:4px">
            ${date.toLocaleString()}
          </div>

          <div style="font-weight:600">
            ${valueFormatter(p.value[1])}
          </div>
        `
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
        type: "line",

        data,

        smooth: true,

        showSymbol: false,

        lineStyle: {
          width: 2,
          color: "#60a5fa",
        },

        areaStyle: {
          color: "rgba(96,165,250,0.08)",
        },
      },
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
  }

  return (
    <ReactECharts
      option={option}
      notMerge
      lazyUpdate
      style={{ width: "100%", height: "100%" }}
    />
  )
}
