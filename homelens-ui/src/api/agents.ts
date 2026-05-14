import { client } from './client';

export type { Agent } from './models';
import type { Agent } from './models';

export const getAgents = (): Promise<Agent[]> => client.get<Agent[]>('/api/agents');
