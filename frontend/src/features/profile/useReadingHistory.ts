import { useQuery } from '@tanstack/react-query';
import { api } from '@/api';
import type { ReadingHistoryFilter } from '@/api/types';

const HISTORY_PAGE_SIZE = 10;

export function useReadingHistory(type: ReadingHistoryFilter, page: number, enabled: boolean) {
  return useQuery({
    queryKey: ['reading-history', type, page, HISTORY_PAGE_SIZE],
    queryFn: () => api.history.getHistory({ type, page, limit: HISTORY_PAGE_SIZE }),
    enabled,
    staleTime: 60 * 1000,
  });
}
