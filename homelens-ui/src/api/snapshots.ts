import { client } from './client';
import type { GetSnapshotsResponse } from './models';

export type { GetSnapshotsResponse, SnapshotEntry } from './models';

export const getSnapshots = (agentId: string): Promise<GetSnapshotsResponse> =>
  client.get<GetSnapshotsResponse>(`/api/agents/${agentId}`);
