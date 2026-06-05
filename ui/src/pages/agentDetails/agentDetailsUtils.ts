import type { SnapshotEntry } from "../../api/models";

export function getSeries<T>(
  history: SnapshotEntry[],
  selector: (snap: SnapshotEntry) => T,
): T[] {
  if (history === undefined) return [];
  return history.map(selector);
}

export function getMultiSeries(
  history: SnapshotEntry[],
  selector: (snap: SnapshotEntry) => number[],
): number[][] {
  const first = selector(history[0]);

  return first.map((_, index) => history.map((snap) => selector(snap)[index]));
}
