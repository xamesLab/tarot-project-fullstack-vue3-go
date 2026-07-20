import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';
import { api } from '@/api';
import { ApiError } from '@/api/client';
import type { BackendReadingDto, RemainingTodayDto } from '@/api/types';
import { useAuthStore } from '@/features/auth/store';
import { getAccountingDate } from '@/features/daily-card/bootstrap';

const spreadKeys = {
  today: (timezone: string) => ['three-card-spread', 'today', timezone] as const,
};

function mergeReading(previous: BackendReadingDto | null | undefined, next: BackendReadingDto): BackendReadingDto {
  if (!previous) return next;
  return {
    ...next,
    remaining_today: next.remaining_today || previous.remaining_today,
  };
}

function asMessage(error: unknown): string {
  return error instanceof Error ? error.message : '';
}

export function useThreeCardSpread() {
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const timezone = user?.timezone || 'UTC';
  const queryKey = spreadKeys.today(timezone);
  const [revealingPositions, setRevealingPositions] = useState<number[]>([]);
  const [revealError, setRevealError] = useState('');
  const [limitInfo, setLimitInfo] = useState<RemainingTodayDto | null>(null);
  const interpretationAttempted = useRef('');

  const spreadQuery = useQuery({
    queryKey,
    enabled: Boolean(user),
    queryFn: async (): Promise<BackendReadingDto | null> => {
      const response = await api.history.getHistory({ type: 'three_card', page: 1, limit: 20 });
      let accountingDate: string;
      try {
        accountingDate = getAccountingDate(new Date(), timezone);
      } catch {
        accountingDate = getAccountingDate(new Date(), 'UTC');
      }
      return response.readings.find((reading) => reading.accounting_date === accountingDate) || null;
    },
  });

  const reading = spreadQuery.data || null;

  useEffect(() => {
    if (reading?.remaining_today) setLimitInfo(reading.remaining_today);
  }, [reading?.remaining_today]);

  const interpretationMutation = useMutation({
    mutationFn: (readingId: string) => api.spreads.requestInterpretation(readingId),
    onMutate: (readingId) => {
      queryClient.setQueryData<BackendReadingDto | null>(queryKey, (current) => current
        ? { ...current, ai: { ...current.ai, status: 'processing', error_code: undefined, error_message: undefined } }
        : current);
      interpretationAttempted.current = readingId;
    },
    onSuccess: ({ request }) => {
      queryClient.setQueryData<BackendReadingDto | null>(queryKey, (current) => current
        ? {
            ...current,
            ai: {
              request_id: request.id,
              provider: request.provider,
              model: request.model,
              prompt_version: request.prompt_version,
              status: request.status,
              content: request.response_excerpt,
              error_code: request.error_code,
              error_message: request.error_message,
              completed_at: request.completed_at,
            },
          }
        : current);
    },
    onError: (error) => {
      queryClient.setQueryData<BackendReadingDto | null>(queryKey, (current) => current
        ? { ...current, ai: { ...current.ai, status: 'failed', error_message: asMessage(error) } }
        : current);
    },
  });

  useEffect(() => {
    if (!reading?.id || !reading.cards.length || !reading.cards.every((card) => card.viewed)) return;
    if (reading.ai?.content || reading.ai?.status === 'failed' || reading.ai?.status === 'processing') return;
    if (interpretationAttempted.current === reading.id) return;
    interpretationMutation.mutate(reading.id);
  }, [reading?.id, reading?.cards, reading?.ai?.content, reading?.ai?.status, interpretationMutation]);

  const generateMutation = useMutation({
    mutationFn: () => api.spreads.generateThreeCard({ timezone }),
    onSuccess: ({ reading: nextReading }) => {
      setRevealError('');
      setLimitInfo(nextReading.remaining_today || null);
      interpretationAttempted.current = '';
      queryClient.setQueryData(queryKey, nextReading);
    },
    onError: (error) => {
      if (error instanceof ApiError && error.code === 'daily_limit_reached' && error.details?.current_limit) {
        setLimitInfo({ count: 0, limit: error.details.current_limit.limit });
      }
    },
  });

  const revealCard = async (position: number) => {
    const current = queryClient.getQueryData<BackendReadingDto | null>(queryKey);
    if (!current?.id || current.cards.find((card) => card.position === position)?.viewed) return;
    setRevealError('');
    setRevealingPositions((positions) => positions.includes(position) ? positions : [...positions, position]);
    try {
      const response = await api.spreads.revealCard(current.id, position);
      queryClient.setQueryData<BackendReadingDto | null>(queryKey, (latest) => mergeReading(latest, response.reading));
    } catch (error) {
      setRevealError(asMessage(error));
      throw error;
    } finally {
      setRevealingPositions((positions) => positions.filter((item) => item !== position));
    }
  };

  const revealAll = async () => {
    const hiddenPositions = (queryClient.getQueryData<BackendReadingDto | null>(queryKey)?.cards || [])
      .filter((card) => !card.viewed)
      .map((card) => card.position);
    for (const position of hiddenPositions) {
      try {
        await revealCard(position);
      } catch {
        break;
      }
    }
  };

  const retryInterpretation = () => {
    if (!reading?.id) return;
    interpretationAttempted.current = '';
    interpretationMutation.mutate(reading.id);
  };

  return {
    reading,
    limitInfo,
    isLoading: spreadQuery.isPending,
    loadError: spreadQuery.error,
    reload: spreadQuery.refetch,
    generate: generateMutation.mutate,
    isGenerating: generateMutation.isPending,
    generateError: generateMutation.error,
    revealCard,
    revealAll,
    revealingPositions,
    revealError,
    isInterpretationLoading: interpretationMutation.isPending || reading?.ai?.status === 'processing',
    interpretationError: reading?.ai?.status === 'failed' ? reading.ai.error_message : '',
    retryInterpretation,
  };
}
