import { useQuery } from '@tanstack/react-query';
import { api } from '@/api';
import type { BackendReadingCardDto, Language } from '@/api/types';
import type { DailyCardBootstrapResult } from './bootstrap';

interface DailyCardInterpretationOptions {
  card?: BackendReadingCardDto;
  dailyCard?: DailyCardBootstrapResult;
  language: Language;
  isGuest: boolean;
  isAuthenticated: boolean;
  isRevealed: boolean;
  userId?: string;
}

export function useDailyCardInterpretation({
  card,
  dailyCard,
  language,
  isGuest,
  isAuthenticated,
  isRevealed,
  userId,
}: DailyCardInterpretationOptions) {
  const readingId = dailyCard?.reading.id;
  const claimToken = dailyCard?.claimToken;

  return useQuery({
    queryKey: [
      'daily-card-interpretation',
      isGuest ? 'guest' : userId ?? 'user',
      readingId ?? '',
      card?.key ?? '',
      claimToken ?? '',
      language,
    ],
    queryFn: () => {
      if (!card) throw new Error('Daily card is not ready');
      if (isGuest) {
        if (!claimToken) throw new Error('Guest claim token is missing');
        return api.dailyCard.getGuestInterpretation({ claim_token: claimToken, language });
      }
      if (!readingId) throw new Error('Reading id is missing');
      return api.dailyCard.getUserInterpretation({ reading_id: readingId, language });
    },
    enabled: isRevealed && Boolean(card) && (isGuest ? Boolean(claimToken) : isAuthenticated && Boolean(readingId)),
    staleTime: Infinity,
    retry: 1,
  });
}
