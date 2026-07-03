import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import { RegistrationModal } from '@/features/auth/RegistrationModal';
import { usePreferencesStore } from '@/features/preferences/store';
import { OracleChat } from '@/features/oracle-chat/OracleChat';
import { Ornament, SectionLabel } from '@/shared/ui/Ornament';
import { CardOfDayStage } from '@/widgets/CardOfDayStage/CardOfDayStage';
import { DateStrip } from '@/widgets/DateStrip/DateStrip';
import { Header } from '@/widgets/Header/Header';
import { Interpretation } from '@/widgets/Interpretation/Interpretation';
import { BackgroundDecor } from '@/widgets/StarField/StarField';
import './home.css';

export function HomePage() {
  const { t } = useTranslation(['common', 'home']);
  const language = usePreferencesStore((state) => state.language);
  const [registrationOpen, setRegistrationOpen] = useState(false);
  const { data: reading, isError } = useQuery({ queryKey: ['daily-reading', language], queryFn: () => api.readings.getDaily(language) });
  return <div className="app-shell"><BackgroundDecor/><Header onGatedAction={() => setRegistrationOpen(true)}/><main className="page-wrap"><DateStrip date={new Date('2026-05-28T12:00:00')} label={t('home:date')}/><div className="home-ornament"><Ornament/></div><SectionLabel>{t('home:cardOfDay')}</SectionLabel><div className="stage"><CardOfDayStage card={reading?.card} hint={t('home:flipHint')}/>{isError ? <p className="reading-error">{t('home:error')}</p> : <Interpretation reading={reading}/>}<OracleChat/></div></main><footer className="site-footer"><div><span>{t('common:copyright')}</span><span>· · ·</span><span>{t('common:footer')}</span></div></footer><RegistrationModal open={registrationOpen} onClose={() => setRegistrationOpen(false)}/></div>;
}
