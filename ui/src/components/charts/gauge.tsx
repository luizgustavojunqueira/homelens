import type { EChartsOption } from 'echarts'
import ReactECharts from 'echarts-for-react'

interface IGauge {
  value: number
  label: string
  used?: string;
  total?: string;
}

export default function Gauge({ value, label, used, total }: IGauge) {
  const option: EChartsOption = {
    backgroundColor: '#0b1220',

    animationDuration: 800,

    tooltip: {
      backgroundColor: '#111827',
      borderColor: '#374151',
      borderWidth: 1,

      textStyle: {
        color: '#e5e7eb',
        fontSize: 12,
      },
      padding: 10,
      formatter: () => {
        if (used !== undefined && total !== undefined) {
          return `
            <div style="font-weight:600;margin-bottom:4px">
              ${label}
            </div>
            <div style="opacity:0.8">
              ${used} / ${total}
            </div>
          `
        }
        return `
          <div style="font-weight:600">
            ${label}
          </div>
          <div>${value.toFixed(1)}%</div>
        `
      },
      position: (point) => {
        return [point[0] - 60, point[1] + 20]
      },
    },

    series: [
      {
        type: 'gauge',

        startAngle: 225,
        endAngle: -45,

        min: 0,
        max: 100,

        pointer: {
          show: false
        },

        progress: {
          show: true,
          roundCap: true,
          width: 18,
        },

        axisLine: {
          lineStyle: {
            width: 18,
            color: [[1, '#1e293b']],
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
          fontWeight: 'bold',
          color: '#ffffff',
          offsetCenter: [0, '-5%'],
          formatter: (v: number) => `${v.toFixed(1)}%`,
        },

        title: {
          offsetCenter: [0, '20%'],
          color: '#94a3b8',
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
  }

  return (
    <ReactECharts
      option={option}
      notMerge
      lazyUpdate
      style={{ width: 260, height: 260 }}
    />
  )
}
