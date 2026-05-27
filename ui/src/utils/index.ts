
type Metric = "B" | "KB" | "MB" | "GB" | "TB" | "PB" | "EB" | "ZB" | "YB";

const metrics: Metric[] = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
export function formatByteStr(bytes: number, currentMetric: Metric = "B"): string {
  let value = bytes;
  let metricIndex = metrics.indexOf(currentMetric);

  const factor = Math.pow(1024, metricIndex);
  value = bytes / factor;

  return `${value.toFixed(2)} ${metrics[metricIndex]}`;
}

export function convertByteToMetric(bytes: number, targetMetric: Metric): number {
  const targetIndex = metrics.indexOf(targetMetric);
  if (targetIndex === -1) {
    throw new Error(`Invalid target metric: ${targetMetric}`);
  }
  const factor = Math.pow(1024, targetIndex);
  return bytes / factor;
}
