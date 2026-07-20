import { ApiError } from '@/api/client';

export const adminKeys = {
  all: ['admin'] as const,
  dashboard: ['admin', 'dashboard'] as const,
  users: (filters: object) => ['admin', 'users', filters] as const,
  policies: ['admin', 'policies'] as const,
  settings: ['admin', 'settings'] as const,
  decks: (scope: string, page: number, limit: number) => ['admin', 'decks', scope, page, limit] as const,
  cardBacks: (scope: string, page: number, limit: number) => ['admin', 'card-backs', scope, page, limit] as const,
};

export function getAdminError(error: unknown, fallback: string) {
  return error instanceof ApiError || error instanceof Error ? error.message : fallback;
}
