import { client } from "./client";

import type { Agent, UpdateNameRequest } from "./models";

export const getAgents = (): Promise<Agent[]> =>
  client.get<Agent[]>("/api/agents");

export const updateName = (request: UpdateNameRequest): Promise<boolean> =>
  client.post<boolean>("/api/agents/update-name", request);
