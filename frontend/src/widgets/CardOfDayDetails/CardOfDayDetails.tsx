import { useTranslation } from 'react-i18next';
import type { BackendReadingCardDto, DailyReadingDto, Language } from '@/api/types';
import { getDailyCardPresentation } from '@/features/daily-card/presentation';
import { StarIcon } from '@/shared/icons';
import { Interpretation } from '@/widgets/Interpretation/Interpretation';

interface CardOfDayDetailsProps {
  backendCard?: BackendReadingCardDto;
  reading?: DailyReadingDto;
  language: Language;
  revealed: boolean;
}

export function CardOfDayDetails({ backendCard, reading, language, revealed }: CardOfDayDetailsProps) {
  const { t } = useTranslation('home');
  const card = backendCard ? getDailyCardPresentation(backendCard, language) : null;
  const subtitle = card ? `${t(`arcana.${card.arcana}`)} · ${t(`orientation.${card.reversed ? 'reversed' : 'upright'}`)}` : '···';

  return <section className={`card-details-shell ${revealed ? 'revealed' : 'concealed'}`} data-revealed={revealed}>
    <div className="card-details-body" aria-hidden={!revealed}>
      <div className="card-label"><span className={`card-number ${card?.arcana === 'minor' ? 'card-number-minor' : ''}`}>{card?.marker ?? '···'}</span><span className="card-title">{card?.name ?? '···'}</span><span className="card-subtitle">{subtitle}</span></div>
      <Interpretation reading={reading} backendCard={backendCard} language={language} revealed={revealed}/>
    </div>
    {!revealed ? <div className="card-details-lock" role="status"><span className="card-details-lock-icon" aria-hidden="true"><StarIcon/></span><span>{t('revealDescription')}</span></div> : null}
  </section>;
}
