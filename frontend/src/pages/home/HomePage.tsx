import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import { AuthModal, type AuthModalMode } from '@/features/auth/AuthModal';
import { useAuthStore } from '@/features/auth/store';
import { useDailyCardBootstrap } from '@/features/daily-card/useDailyCardBootstrap';
import { usePreferencesStore } from '@/features/preferences/store';
import { Tooltip } from '@/shared/ui/Content';
import { Ornament, SectionLabel } from '@/shared/ui/Ornament';
import { OracleChat } from '@/features/oracle-chat/OracleChat';
import { CardOfDayDetails } from '@/widgets/CardOfDayDetails/CardOfDayDetails';
import { CardOfDayStage } from '@/widgets/CardOfDayStage/CardOfDayStage';
import { DateStrip } from '@/widgets/DateStrip/DateStrip';
import { Header } from '@/widgets/Header/Header';
import { BackgroundDecor } from '@/widgets/StarField/StarField';
import './home.css';

export function HomePage() {
  const { t } = useTranslation(['common', 'home']);
  const language = usePreferencesStore((state) => state.language);
  const authStatus = useAuthStore((state) => state.status);
  const hydrated = useAuthStore((state) => state.hydrated);
  const user = useAuthStore((state) => state.user);
  const [authModal, setAuthModal] = useState<{ open: boolean; mode: AuthModalMode }>({ open: false, mode: 'register' });
  const [revealedCardKey, setRevealedCardKey] = useState<string | null>(null);
  const dailyCardQuery = useDailyCardBootstrap();
  const backendCard = dailyCardQuery.data?.reading.cards[0];
  const isGuest = hydrated && authStatus === 'guest';
  const detailsRevealed = Boolean(backendCard && revealedCardKey === backendCard.key);

  const interpretationQuery = useQuery({
    queryKey: [
      'daily-card-interpretation',
      isGuest ? 'guest' : user?.id ?? 'user',
      dailyCardQuery.data?.reading.id ?? '',
      backendCard?.key ?? '',
      dailyCardQuery.data?.claimToken ?? '',
      language,
    ],
    queryFn: () => {
      if (!backendCard) throw new Error('Daily card is not ready');
      if (isGuest) {
        if (!dailyCardQuery.data?.claimToken) throw new Error('Guest claim token is missing');
        return api.dailyCard.getGuestInterpretation({ claim_token: dailyCardQuery.data.claimToken, language });
      }

      const readingId = dailyCardQuery.data?.reading.id;
      if (!readingId) throw new Error('Reading id is missing');
      return api.dailyCard.getUserInterpretation({ reading_id: readingId, language });
    },
    enabled: detailsRevealed && Boolean(backendCard) && (
      isGuest ? Boolean(dailyCardQuery.data?.claimToken) : authStatus === 'authenticated' && Boolean(dailyCardQuery.data?.reading.id)
    ),
    staleTime: Infinity,
    retry: 1,
  });

  return (
    <div className="app-shell">
      <BackgroundDecor />
      <Header onAuthRequest={(mode) => setAuthModal({ open: true, mode })} />
      <main className="page-wrap">
        <DateStrip date={new Date('2026-05-28T12:00:00')} label={t('home:date')} />
        <div className="home-ornament"><Ornament /></div>
        <SectionLabel>
          <span className="daily-card-heading">
            <span>{t('home:cardOfDay')}{isGuest ? ` ${t('home:guestCardSuffix')}` : ''}</span>
            {isGuest ? (
              <Tooltip below label={t('home:guestCardTooltip')}>
                <button type="button" className="info-dot" aria-label={t('home:guestCardTooltip')}>i</button>
              </Tooltip>
            ) : null}
          </span>
        </SectionLabel>
        <div className="stage">
          <CardOfDayStage
            backendCard={backendCard}
            language={language}
            hint={t('home:flipHint')}
            onRevealChange={(revealed) => setRevealedCardKey(revealed ? backendCard?.key ?? null : null)}
          />
          {dailyCardQuery.isError ? (
            <p className="reading-error">{t('home:error')}</p>
          ) : (
            <CardOfDayDetails
              backendCard={backendCard}
              interpretation={interpretationQuery.data?.interpretation}
              interpretationLoading={interpretationQuery.isFetching}
              interpretationError={interpretationQuery.isError ? t('home:error') : ''}
              language={language}
              revealed={detailsRevealed}
            />
          )}
          <OracleChat />
        </div>
      </main>
      <footer className="site-footer">
        <div><span>{t('common:copyright')}</span><span> · · · </span><span>{t('common:footer')}</span></div>
      </footer>
      <AuthModal open={authModal.open} initialMode={authModal.mode} onClose={() => setAuthModal((current) => ({ ...current, open: false }))} />
    </div>
  );
}
