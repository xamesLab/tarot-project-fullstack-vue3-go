import { useQuery } from '@tanstack/react-query';
import { useEffect } from 'react';
import { useAuthStore } from '@/features/auth/store';
import { bootstrapDailyCard, getAccountingDate } from './bootstrap';

export function useDailyCardBootstrap() {
  const hydrated = useAuthStore((state) => state.hydrated);
  const user = useAuthStore((state) => state.user);
  const timezone = user?.timezone || Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';
  const accountingDate = getAccountingDate(new Date(), timezone);
  const query = useQuery({
    queryKey: ['backend-daily-card', user?.id ?? 'guest', accountingDate],
    queryFn: () => bootstrapDailyCard(user),
    enabled: hydrated,
    staleTime: Infinity,
  });

  useEffect(() => {
    if (query.data) console.info('[daily-card]', query.data);
  }, [query.data]);

  useEffect(() => {
    if (query.error) console.error('[daily-card]', query.error);
  }, [query.error]);

  return query;
}
