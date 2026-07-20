import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { resolveApiAssetUrl } from '@/api/client';
import type { Language, ReadingHistoryFilter, ReadingHistoryItemDto } from '@/api/types';
import { usePreferencesStore } from '@/features/preferences/store';
import { Button } from '@/shared/ui/Button/Button';
import { Checkbox } from '@/shared/ui/Controls/Controls';
import { TarotCard } from '@/shared/ui/TarotCard';
import { useReadingHistory } from './useReadingHistory';

function formatReadingDate(value: string, language: Language) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat(language === 'ru' ? 'ru-RU' : 'en-US', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date);
}

function formatRetentionDate(value: string, language: Language) {
  const date = new Date(`${value}T00:00:00`);
  if (Number.isNaN(date.getTime())) return value;
  const formatted = new Intl.DateTimeFormat(language === 'ru' ? 'ru-RU' : 'en-US', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }).format(date);
  return language === 'ru' ? formatted.replace(/\.$/, '') : formatted;
}

export function ProfileHistory() {
  const { t } = useTranslation('profile');
  const language = usePreferencesStore((state) => state.language);
  const [showDailyCards, setShowDailyCards] = useState(true);
  const [showOtherReadings, setShowOtherReadings] = useState(true);
  const [page, setPage] = useState(1);
  const [expandedIds, setExpandedIds] = useState<string[]>([]);

  const filter: ReadingHistoryFilter | null = showDailyCards && showOtherReadings
    ? 'all'
    : showDailyCards
      ? 'card_of_day'
      : showOtherReadings
        ? 'three_card'
        : null;
  const history = useReadingHistory(filter ?? 'all', page, filter !== null);
  const readings = filter ? history.data?.readings ?? [] : [];
  const pagination = history.data?.pagination;

  const updateDailyFilter = (checked: boolean) => {
    setShowDailyCards(checked);
    setPage(1);
  };

  const updateOtherFilter = (checked: boolean) => {
    setShowOtherReadings(checked);
    setPage(1);
  };

  const toggleReading = (id: string) => {
    setExpandedIds((current) => (
      current.includes(id) ? current.filter((item) => item !== id) : [...current, id]
    ));
  };

  return (
    <div className="profile-history">
      <section className="profile-history-toolbar" aria-labelledby="profile-history-title">
        <div>
          <h2 id="profile-history-title">{t('history.title')}</h2>
          <p>{t('history.description')}</p>
        </div>
        <div className="profile-history-filters" role="group" aria-label={t('history.filterLabel')}>
          <Checkbox checked={showDailyCards} onChange={updateDailyFilter}>
            {t('history.filters.dailyCards')}
          </Checkbox>
          <Checkbox checked={showOtherReadings} onChange={updateOtherFilter}>
            {t('history.filters.otherReadings')}
          </Checkbox>
        </div>
        {history.data?.retention_info && (
          <p className="profile-history-retention">
            {t('history.retention', {
              days: history.data.retention_info.retention_days,
              date: formatRetentionDate(history.data.retention_info.oldest_date, language),
            })}
          </p>
        )}
      </section>

      {!filter ? (
        <HistoryState title={t('history.filtersEmptyTitle')} copy={t('history.filtersEmptyCopy')} />
      ) : history.isPending ? (
        <div className="profile-history-state" role="status">{t('history.loading')}</div>
      ) : history.isError ? (
        <div className="profile-history-state profile-history-state-error" role="alert">
          <span>{t('history.error')}</span>
          <Button type="button" variant="outline" size="sm" onClick={() => history.refetch()}>
            {t('history.retry')}
          </Button>
        </div>
      ) : readings.length === 0 ? (
        <HistoryState title={t('history.emptyTitle')} copy={t('history.emptyCopy')} />
      ) : (
        <div className="profile-history-list" aria-busy={history.isFetching}>
          {readings.map((reading) => (
            <ReadingHistoryItem
              key={reading.id}
              reading={reading}
              language={language}
              expanded={expandedIds.includes(reading.id)}
              onToggle={() => toggleReading(reading.id)}
            />
          ))}
        </div>
      )}

      {filter && pagination && pagination.total_pages > 1 && (
        <nav className="profile-history-pagination" aria-label={t('history.paginationLabel')}>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            disabled={pagination.page <= 1 || history.isFetching}
            onClick={() => setPage((current) => Math.max(1, current - 1))}
          >
            {t('history.previousPage')}
          </Button>
          <span>{t('history.pageIndicator', { page: pagination.page, total: pagination.total_pages })}</span>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            disabled={pagination.page >= pagination.total_pages || history.isFetching}
            onClick={() => setPage((current) => current + 1)}
          >
            {t('history.nextPage')}
          </Button>
        </nav>
      )}
    </div>
  );
}

function ReadingHistoryItem({
  reading,
  language,
  expanded,
  onToggle,
}: {
  reading: ReadingHistoryItemDto;
  language: Language;
  expanded: boolean;
  onToggle: () => void;
}) {
  const { t } = useTranslation('profile');
  const detailsId = `profile-history-details-${reading.id}`;
  const cardName = (card: ReadingHistoryItemDto['cards'][number]) => (
    language === 'ru' ? card.name_ru || card.name : card.name || card.name_ru
  );
  const interpretation = (card: ReadingHistoryItemDto['cards'][number]) => (
    language === 'ru'
      ? card.interpretation_ru || card.interpretation
      : card.interpretation || card.interpretation_ru
  );

  return (
    <article className="profile-history-item">
      <div className="profile-history-summary">
        <div>
          <span className={`profile-history-type profile-history-type-${reading.type}`}>
            {reading.type === 'card_of_day' ? t('history.types.dailyCard') : t('history.types.otherReading')}
          </span>
          <h3>{formatReadingDate(reading.generated_at, language)}</h3>
          <p>{t('history.accountingDate', { date: reading.accounting_date })}</p>
        </div>
        <div className="profile-history-summary-actions">
          {reading.ai?.content && <span className="profile-history-ai-badge">{t('history.aiReady')}</span>}
          <Button
            type="button"
            variant="ghost"
            size="sm"
            aria-expanded={expanded}
            aria-controls={detailsId}
            onClick={onToggle}
          >
            {expanded ? t('history.hideDetails') : t('history.showDetails')}
          </Button>
        </div>
      </div>

      <div className="profile-history-chips" aria-label={t('history.cardsLabel')}>
        {reading.cards.map((card) => (
          <span key={`${reading.id}-${card.key}-${card.position}`}>
            {cardName(card)}{card.reversed ? ` · ${t('history.reversed')}` : ''}
          </span>
        ))}
      </div>

      {expanded && (
        <div id={detailsId} className="profile-history-details">
          <div className="profile-history-cards">
            {reading.cards.map((card) => (
              <article className="profile-history-card" key={`${reading.id}-details-${card.key}-${card.position}`}>
                <div className="profile-history-card-visual">
                  <TarotCard
                    roman={card.arcana === 'major' ? '✦' : ''}
                    name={cardName(card)}
                    faceImageUrl={resolveApiAssetUrl(card.image_url)}
                    faceImageAlt={cardName(card)}
                    backImageUrl={resolveApiAssetUrl(card.back_image_url)}
                  />
                </div>
                <div>
                  <h4>{cardName(card)}</h4>
                  <p>{interpretation(card)}</p>
                </div>
              </article>
            ))}
          </div>
          {reading.ai?.content && (
            <section className="profile-history-meaning" aria-label={t('history.meaningTitle')}>
              <span>{t('history.meaningTitle')}</span>
              <p>{reading.ai.content}</p>
            </section>
          )}
        </div>
      )}
    </article>
  );
}

function HistoryState({ title, copy }: { title: string; copy: string }) {
  return (
    <div className="profile-history-empty">
      <span aria-hidden="true">☾</span>
      <h3>{title}</h3>
      <p>{copy}</p>
    </div>
  );
}
