export type ApiResponse<T> = {
  ok: boolean;
  message: string;
  data: T;
  timestamp: string;
};

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers || {}),
    },
    ...init,
  });
  const body = (await res.json()) as ApiResponse<T>;
  if (!res.ok || !body.ok) {
    throw new Error(body?.message || `HTTP ${res.status}`);
  }
  return body.data;
}

export const api = {
  me: () => request<any>("/v1/auth/me"),
  repos: () => request<any[]>("/v1/git/repos"),
  projects: () => request<any[]>("/v1/projects"),
  createProject: (payload: unknown) =>
    request<any>("/v1/projects", { method: "POST", body: JSON.stringify(payload) }),
  connectRepo: (projectId: string, payload: unknown) =>
    request<any>(`/v1/projects/${projectId}/repo/connect`, {
      method: "POST",
      body: JSON.stringify(payload),
    }),
  deployments: (projectId: string) =>
    request<any[]>(`/v1/projects/${projectId}/deployments`),
  createDeployment: (projectId: string) =>
    request<any>(`/v1/projects/${projectId}/deployments`, { method: "POST" }),
  deployment: (deploymentId: string) => request<any>(`/v1/deployments/${deploymentId}`),
  deploymentLogs: (deploymentId: string) =>
    request<{ logs: string; deployment_id: string }>(`/v1/deployments/${deploymentId}/logs`),
  runtimeStatus: (projectId: string) => request<any>(`/v1/projects/${projectId}/runtime`),
  runtimeAction: (projectId: string, action: "start" | "stop" | "restart") =>
    request<any>(`/v1/projects/${projectId}/runtime/${action}`, { method: "POST" }),
  envList: (projectId: string) => request<any[]>(`/v1/projects/${projectId}/env`),
  envUpsert: (projectId: string, key: string, value: string) =>
    request<any>(`/v1/projects/${projectId}/env`, {
      method: "POST",
      body: JSON.stringify({ key, value }),
    }),
  envDelete: (projectId: string, key: string) =>
    request<any>(`/v1/projects/${projectId}/env/${encodeURIComponent(key)}`, {
      method: "DELETE",
    }),
  telegramStatus: (projectId: string) =>
    request<any>(`/v1/projects/${projectId}/telegram/status`),
  telegramConfig: (projectId: string, payload: unknown) =>
    request<any>(`/v1/projects/${projectId}/telegram/config`, {
      method: "POST",
      body: JSON.stringify(payload),
    }),
};
