import { client } from './client';
import type { GetSnapshotsResponse } from './models';

export type { GetSnapshotsResponse, AgentSnapshots, SnapshotEntry } from './models';

export const getSnapshots = (): Promise<GetSnapshotsResponse> =>
  client.get<GetSnapshotsResponse>('/api/snapshots');
