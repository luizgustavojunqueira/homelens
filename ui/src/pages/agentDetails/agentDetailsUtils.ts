import type { SnapshotEntry } from "../../api/models";

export function getSeries<T>(
  history: SnapshotEntry[],
  selector: (snap: SnapshotEntry) => T,
): T[] {
  if (!history || history.length === 0) return [];
  return history.map(selector);
}

export function getMultiSeries(
  history: SnapshotEntry[],
  selector: (snap: SnapshotEntry) => number[],
): number[][] {
  if (!history || history.length === 0) return [];
  
  const first = selector(history[0]);
  if (!first) return [];

  return first.map((_, index) => history.map((snap) => selector(snap)[index]));
}
