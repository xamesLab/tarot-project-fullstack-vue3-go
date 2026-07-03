export class ApiError extends Error {
  constructor(public readonly status: number, message: string) { super(message); this.name = 'ApiError'; }
}

export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${import.meta.env.VITE_API_URL ?? ''}${path}`, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...init?.headers },
  });
  if (!response.ok) throw new ApiError(response.status, `API request failed: ${response.status}`);
  return response.json() as Promise<T>;
}
