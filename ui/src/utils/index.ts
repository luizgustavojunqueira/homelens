type Metric = "B" | "KB" | "MB" | "GB" | "TB" | "PB" | "EB" | "ZB" | "YB";

export function formatByteStr(bytes: number, currentMetric: Metric = "B"): string {
  const metrics: Metric[] = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  let value = bytes;
  let metricIndex = metrics.indexOf(currentMetric);

  while (value >= 1024 && metricIndex < metrics.length - 1) {
    value /= 1024;
    metricIndex++;
  }

  return `${value.toFixed(2)} ${metrics[metricIndex]}`;
}
