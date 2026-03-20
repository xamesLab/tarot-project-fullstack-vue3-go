import type { GatewayHealthResponse } from "../types/reading";
import { request } from "./apiClient";

export function getGatewayHealth(token: string): Promise<GatewayHealthResponse> {
  return request<GatewayHealthResponse>("/api/gateway/health", { token });
}
