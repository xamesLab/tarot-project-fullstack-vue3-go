import { useEffect, useMemo, useState, type CSSProperties } from 'react';
import { useTranslation } from 'react-i18next';
import { ApiError, resolveApiAssetUrl } from '@/api/client';
import type { BackendReadingCardDto } from '@/api/types';
import { useAuthStore } from '@/features/auth/store';
import { useThreeCardSpread } from '@/features/spreads/useThreeCardSpread';
import { usePreferencesStore } from '@/features/preferences/store';
import { ChevronIcon, ZodiacWheel } from '@/shared/icons';
import { Button } from '@/shared/ui/Button';
import { Ornament } from '@/shared/ui/Ornament';
import { TarotCard } from '@/shared/ui/TarotCard';
import { CelestialStrip } from '@/widgets/CelestialStrip/CelestialStrip';
import { Footer } from '@/widgets/Footer/Footer';
import { Header } from '@/widgets/Header/Header';
import { BackgroundDecor } from '@/widgets/StarField/StarField';
import './spreads.css';

const romanPositions = ['I', 'II', 'III'];

function localizedCard(card: BackendReadingCardDto, language: 'ru' | 'en') {
  return {
    name: language === 'ru' ? card.name_ru : card.name,
    interpretation: language === 'ru' ? card.interpretation_ru : card.interpretation,
    tags: (language === 'ru' ? card.tags_ru : card.tags)
      .split(',')
      .map((tag) => tag.trim())
      .filter(Boolean),
  };
}

function EyeIcon() {
  return <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true"><path d="M2 12q10-8 20 0-10 8-20 0Z" strokeWidth="1.4"/><circle cx="12" cy="12" r="3.4" strokeWidth="1.4"/></svg>;
}

export default function SpreadsPage() {
  const { t } = useTranslation(['common', 'spreads']);
  const language = usePreferencesStore((state) => state.language);
  const userTimeZone = useAuthStore((state) => state.user?.timezone);
  const spread = useThreeCardSpread();
  const [expanded, setExpanded] = useState(false);

  useEffect(() => setExpanded(false), [spread.reading?.id]);

  const cards = useMemo(
    () => [...(spread.reading?.cards || [])].sort((left, right) => left.position - right.position),
    [spread.reading?.cards],
  );
  const allRevealed = cards.length === 3 && cards.every((card) => card.viewed);
  const anyRevealing = spread.revealingPositions.length > 0;
  const localizedCards = cards.map((card) => localizedCard(card, language));
  const keywords = [...new Set(localizedCards.flatMap((card) => card.tags))].slice(0, 3);

  const generateErrorCode = spread.generateError instanceof ApiError ? spread.generateError.code : '';
  const generateErrorText = generateErrorCode === 'daily_limit_reached'
    ? t('spreads:errors.limit')
    : spread.generateError
      ? t('spreads:errors.generate')
      : '';

  const content = (() => {
    if (spread.isLoading) {
      return <div className="spread-state" role="status"><span>✦</span>{t('spreads:states.loading')}</div>;
    }

    if (spread.loadError) {
      return (
        <div className="spread-state spread-state-error" role="alert">
          <span>☾</span>
          <p>{t('spreads:errors.load')}</p>
          <Button variant="outline" size="sm" onClick={() => void spread.reload()}>{t('spreads:actions.retry')}</Button>
        </div>
      );
    }

    if (!spread.reading) {
      return (
        <section className="spread-empty" aria-labelledby="spread-empty-title">
          <div className="spread-empty-sigil" aria-hidden="true">✦</div>
          <h2 id="spread-empty-title">{t('spreads:empty.title')}</h2>
          <p>{t('spreads:empty.copy')}</p>
          {spread.limitInfo && <div className="spread-limit">{t('spreads:remaining', { count: spread.limitInfo.count, limit: spread.limitInfo.limit })}</div>}
          <Button
            onClick={() => spread.generate()}
            disabled={spread.isGenerating || spread.limitInfo?.count === 0}
          >
            {spread.isGenerating ? t('spreads:actions.generating') : t('spreads:actions.generate')}
          </Button>
          {generateErrorText && <div className="spread-action-error" role="alert">{generateErrorText}</div>}
        </section>
      );
    }

    return (
      <>
        <div className="spread-toolbar">
          {spread.limitInfo
            ? <span className="spread-limit">{t('spreads:remaining', { count: spread.limitInfo.count, limit: spread.limitInfo.limit })}</span>
            : <span className="spread-limit">{t('spreads:remainingUnknown')}</span>}
          <Button
            variant="ghost"
            size="sm"
            onClick={() => spread.generate()}
            disabled={spread.isGenerating || spread.limitInfo?.count === 0 || anyRevealing}
          >
            {spread.isGenerating ? t('spreads:actions.generating') : t('spreads:actions.new')}
          </Button>
        </div>

        <section className="spread-stage" aria-label={t('spreads:stageLabel')}>
          <ZodiacWheel className="spread-zodiac" />
          <div className="spread-mobile-captions">
            {cards.map((card, index) => (
              <div className="spread-mobile-caption" key={`mobile-caption-${card.position}`}>
                <div className="spread-position">{t(`spreads:positions.${index}.label`)}</div>
                <div className="spread-position-sub">{t(`spreads:positions.${index}.subtitle`)}</div>
              </div>
            ))}
          </div>
          <div className="spread-row">
            {cards.map((card, index) => {
              const localized = localizedCards[index];
              const revealing = spread.revealingPositions.includes(card.position);
              return (
                <article className={`spread-card-shell ${card.viewed ? 'revealed' : ''}`} key={`${spread.reading?.id}-${card.position}`} style={{ '--deal-index': index } as CSSProperties}>
                  <div className="spread-card-caption">
                    <div className="spread-position">{t(`spreads:positions.${index}.label`)}</div>
                    <div className="spread-position-sub">{t(`spreads:positions.${index}.subtitle`)}</div>
                  </div>
                  <div className="spread-card-float">
                    <div className="spread-card-glow" />
                    <div className="spread-card-visual">
                      <TarotCard
                        roman={romanPositions[index] || String(index + 1)}
                        name={localized.name}
                        faceImageUrl={resolveApiAssetUrl(card.image_url)}
                        faceImageAlt={localized.name}
                        backImageUrl={resolveApiAssetUrl(card.back_image_url)}
                        flippable
                        revealed={card.viewed}
                        revealLoading={revealing}
                        onReveal={() => void spread.revealCard(card.position).catch(() => undefined)}
                        flipLabel={card.viewed
                          ? t('spreads:cardRevealed', { name: localized.name })
                          : t('spreads:revealCard', { position: t(`spreads:positions.${index}.label`) })}
                      />
                    </div>
                  </div>
                  <div className="spread-tap-hint">
                    {revealing
                      ? t('spreads:states.revealing')
                      : card.viewed
                        ? localized.name
                        : t('spreads:actions.tap')}
                  </div>
                </article>
              );
            })}
          </div>
          <div className="spread-mobile-hints">
            {cards.map((card, index) => {
              const localized = localizedCards[index];
              const revealing = spread.revealingPositions.includes(card.position);
              return (
                <div className={`spread-mobile-hint ${card.viewed ? 'revealed' : ''}`} key={`mobile-hint-${card.position}`}>
                  {revealing
                    ? t('spreads:states.revealing')
                    : card.viewed
                      ? localized.name
                      : t('spreads:actions.tap')}
                </div>
              );
            })}
          </div>

          {!allRevealed && (
            <button className="spread-reveal-all" type="button" onClick={() => void spread.revealAll()} disabled={anyRevealing}>
              <EyeIcon />
              <span>{anyRevealing ? t('spreads:states.revealing') : t('spreads:actions.revealAll')}</span>
            </button>
          )}
          {spread.revealError && <div className="spread-action-error" role="alert">{t('spreads:errors.reveal')}</div>}
        </section>

        {allRevealed && (
          <section className="spread-reading" aria-label={t('spreads:interpretation.label')}>
            {keywords.length > 0 && <div className="spread-keywords">{keywords.map((keyword) => <span className="chip" key={keyword}>{keyword}</span>)}</div>}
            <p className="spread-brief">{t('spreads:interpretation.brief', { cards: localizedCards.map((card) => card.name).join(' · ') })}</p>
            <div className="spread-interpretation-ornament"><Ornament /></div>
            <button
              className={`spread-expand ${expanded ? 'open' : ''}`}
              type="button"
              aria-expanded={expanded}
              aria-controls="spread-interpretation-details"
              onClick={() => setExpanded((value) => !value)}
            >
              <span>{expanded ? t('spreads:actions.collapse') : t('spreads:actions.expand')}</span>
              <ChevronIcon />
            </button>
            {expanded && (
              <div id="spread-interpretation-details" className="spread-interpretation-details">
                <div className="spread-card-meanings">
                  {cards.map((card, index) => {
                    const localized = localizedCards[index];
                    return (
                      <article className="spread-card-meaning" key={`meaning-${card.position}`}>
                        <div className="spread-meaning-symbol" aria-hidden="true">{index === 0 ? '☾' : index === 1 ? '✦' : '☉'}</div>
                        <div className="spread-meaning-position">{t(`spreads:positions.${index}.label`)}</div>
                        <h3>{localized.name}</h3>
                        <p>{localized.interpretation}</p>
                      </article>
                    );
                  })}
                </div>
                <div className="spread-synthesis">
                  <div className="spread-synthesis-title">{t('spreads:interpretation.connection')}</div>
                  {spread.isInterpretationLoading && <p role="status">{t('spreads:states.interpreting')}</p>}
                  {!spread.isInterpretationLoading && spread.reading.ai?.content && <p>{spread.reading.ai.content}</p>}
                  {!spread.isInterpretationLoading && spread.interpretationError && (
                    <div className="spread-ai-error" role="alert">
                      <p>{t('spreads:errors.interpretation')}</p>
                      <Button variant="outline" size="sm" onClick={spread.retryInterpretation}>{t('spreads:actions.retry')}</Button>
                    </div>
                  )}
                </div>
              </div>
            )}
          </section>
        )}
        {generateErrorText && <div className="spread-action-error spread-new-error" role="alert">{generateErrorText}</div>}
      </>
    );
  })();

  return (
    <div className="app-shell spreads-shell">
      <BackgroundDecor />
      <Header />
      <main className="page-wrap spreads-page">
        <CelestialStrip timeZone={userTimeZone} />
        <header className="spreads-heading">
          <h1>{t('spreads:title')}</h1>
          <p>{t('spreads:subtitle')}</p>
        </header>
        {content}
      </main>
      <Footer />
    </div>
  );
}
