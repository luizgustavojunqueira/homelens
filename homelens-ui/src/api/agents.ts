import { client } from './client';

export interface Agent {
  id: string;
  name: string;
  last_seen: string;
  online: boolean;
}

export const getAgents = (): Promise<Agent[]> => client.get<Agent[]>('/api/agents');
