type Metric = "B" | "KB" | "MB" | "GB" | "TB" | "PB" | "EB" | "ZB" | "YB";

const metrics: Metric[] = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
export function formatByteStr(
  value: number,
  currentMetric: Metric = "B",
): string {
  if (value === 0) {
    return `${value.toFixed(2)} ${currentMetric}`;
  }

  const metricIndex = metrics.indexOf(currentMetric);

  const factor = Math.pow(1024, metricIndex);
  const valueInBytes = value * factor;

  const newMetricIndex = Math.floor(Math.log(valueInBytes) / Math.log(1024));
  const newValue = valueInBytes / Math.pow(1024, newMetricIndex);

  return `${newValue.toFixed(2)} ${metrics[newMetricIndex]}`;
}

export function convertByteToMetric(
  value: number,
  targetMetric: Metric,
  currentMetric: Metric = "B",
): number {
  if (value === 0) {
    return value;
  }

  const currentMetricIndex = metrics.indexOf(currentMetric);

  const factorCurrent = Math.pow(1024, currentMetricIndex);
  const valueInBytes = value * factorCurrent;

  const targetIndex = metrics.indexOf(targetMetric);
  if (targetIndex === -1) {
    throw new Error(`Invalid target metric: ${targetMetric}`);
  }

  const factor = Math.pow(1024, targetIndex);
  return valueInBytes / factor;
}
