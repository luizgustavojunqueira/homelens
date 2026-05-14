
import { client } from './client';

export interface Snapshot {
  id: number;
  agentId: string;
  timestamp: string;
  data: any
}

export interface Snapshots {
  snapshots: Snapshot[];
}

export const getSnapshots = (): Promise<Snapshots[]> => client.get<Snapshots[]>('/api/snapshots');



