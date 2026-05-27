
import { client } from './client';

import type { GetSnapshotsResponse } from './models';

export const getSnapshots = (agent_id: string): Promise<GetSnapshotsResponse> => client.get<GetSnapshotsResponse>(`/api/agents/${agent_id}`)
