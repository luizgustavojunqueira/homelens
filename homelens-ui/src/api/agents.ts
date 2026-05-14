import { client } from './client';

import type { Agent } from './models';

export const getAgents = (): Promise<Agent[]> => client.get<Agent[]>('/api/agents');
