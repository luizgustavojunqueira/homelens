import { client } from "./client";

import type {
  GetAlertConfigResponse,
  UpdateAlertConfigRequest,
} from "./models";

export const getAlertConfig = (): Promise<GetAlertConfigResponse> =>
  client.get<GetAlertConfigResponse>("/api/alerts");

export const saveAlertConfig = (
  request: UpdateAlertConfigRequest,
): Promise<boolean> => client.post<boolean>("/api/alerts", request);
