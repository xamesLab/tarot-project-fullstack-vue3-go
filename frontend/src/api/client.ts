import type { ApiErrorDetails } from '@/api/types';

export const TOKEN_STORAGE_KEY = 'tarot_daily_token';
export const AUTH_UNAUTHORIZED_EVENT = 'nebula:auth-unauthorized';

export function resolveApiAssetUrl(path?: string): string {
  if (!path) return '';
  if (/^[a-z][a-z\d+.-]*:/i.test(path)) return path;

  const apiUrl = new URL(import.meta.env.VITE_API_URL || '/api', window.location.origin);
  apiUrl.pathname = apiUrl.pathname.replace(/\/api\/?$/, '/');
  if (!apiUrl.pathname.endsWith('/')) apiUrl.pathname += '/';
  apiUrl.search = '';
  apiUrl.hash = '';

  return new URL(path.replace(/^\/+/, ''), apiUrl).toString();
}

interface ErrorPayload {
  error?: string;
  message?: string;
  details?: ApiErrorDetails;
  current_limit?: ApiErrorDetails['current_limit'];
  next_reset_at?: string;
}

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    public readonly code: string,
    message: string,
    public readonly details?: ApiErrorDetails,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export async function request<T>(path: string, init?: RequestInit, options: { auth?: boolean } = {}): Promise<T> {
  const baseUrl = (import.meta.env.VITE_API_URL || '/api').replace(/\/$/, '');
  const token = localStorage.getItem(TOKEN_STORAGE_KEY);
  const useAuth = options.auth !== false;
  const headers = new Headers(init?.headers);
  const isFormData = typeof FormData !== 'undefined' && init?.body instanceof FormData;
  if (init?.body && !isFormData && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json');
  if (useAuth && token && !headers.has('Authorization')) headers.set('Authorization', `Bearer ${token}`);

  const response = await fetch(`${baseUrl}${path}`, {
    ...init,
    headers,
  });
  const payload = await response.json().catch(() => ({})) as T & ErrorPayload;
  if (!response.ok) {
    if (response.status === 401 && useAuth && token) {
      window.dispatchEvent(new CustomEvent(AUTH_UNAUTHORIZED_EVENT));
    }
    throw new ApiError(
      response.status,
      payload.error || 'request_failed',
      payload.message || `API request failed: ${response.status}`,
      payload.details || payload.current_limit || payload.next_reset_at
        ? { ...payload.details, current_limit: payload.current_limit, next_reset_at: payload.next_reset_at }
        : undefined,
    );
  }
  return payload;
}

export function post<TResponse, TPayload>(
  path: string,
  payload: TPayload,
  options?: { auth?: boolean },
): Promise<TResponse> {
  return request(path, { method: 'POST', body: JSON.stringify(payload) }, options);
}
