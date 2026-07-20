import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { AuthModal, type AuthModalMode } from '@/features/auth/AuthModal';
import { useAuthStore } from '@/features/auth/store';
import { useDailyCardBootstrap } from '@/features/daily-card/useDailyCardBootstrap';
import { useDailyCardInterpretation } from '@/features/daily-card/useDailyCardInterpretation';
import { usePreferencesStore } from '@/features/preferences/store';
import { Tooltip } from '@/shared/ui/Content';
import { Ornament, SectionLabel } from '@/shared/ui/Ornament';
import { OracleChat } from '@/features/oracle-chat/OracleChat';
import { CardOfDayDetails } from '@/widgets/CardOfDayDetails/CardOfDayDetails';
import { CardOfDayStage } from '@/widgets/CardOfDayStage/CardOfDayStage';
import { CelestialStrip } from '@/widgets/CelestialStrip/CelestialStrip';
import { Footer } from '@/widgets/Footer/Footer';
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

  const interpretationQuery = useDailyCardInterpretation({
    card: backendCard,
    dailyCard: dailyCardQuery.data,
    language,
    isGuest,
    isAuthenticated: authStatus === 'authenticated',
    isRevealed: detailsRevealed,
    userId: user?.id,
  });

  return (
    <div className="app-shell">
      <BackgroundDecor />
      <Header onAuthRequest={(mode) => setAuthModal({ open: true, mode })} />
      <main className="page-wrap home-page">
        <CelestialStrip timeZone={user?.timezone} />
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
      <Footer />
      <AuthModal open={authModal.open} initialMode={authModal.mode} onClose={() => setAuthModal((current) => ({ ...current, open: false }))} />
    </div>
  );
}
