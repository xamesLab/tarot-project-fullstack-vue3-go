import { useTranslation } from 'react-i18next';
import { usePreferencesStore } from '@/features/preferences/store';
import { Logo } from '@/shared/ui/Logo';
import { NavBar, type NavItem } from '@/shared/ui/NavBar';
import { SegmentedControl } from '@/shared/ui/Controls';
import { UserPill } from '@/shared/ui/Avatar';

export function Header({ onGatedAction }: { onGatedAction: (trigger: HTMLButtonElement) => void }) {
  const { t } = useTranslation();
  const theme = usePreferencesStore((state) => state.theme);
  const language = usePreferencesStore((state) => state.language);
  const setLanguage = usePreferencesStore((state) => state.setLanguage);
  const toggleTheme = usePreferencesStore((state) => state.toggleTheme);
  const items: NavItem[] = [
    { id: 'daily', label: t('nav.daily') }, { id: 'spreads', label: t('nav.spreads') },
    { id: 'horoscope', label: t('nav.horoscope') }, { id: 'learn', label: t('nav.learn') },
  ];
  return <header className="site-header"><div className="header-inner"><Logo tagline={t('tagline')}/><NavBar items={items} active="daily" label={t('navigationLabel')} onSelect={(_item, trigger) => onGatedAction(trigger)}/><div className="header-actions"><UserPill initials={language === 'ru' ? 'МА' : 'MA'} name={t('user.name')} sign={t('user.sign')}/><SegmentedControl value={language} options={[{ value: 'ru', label: 'RU' }, { value: 'en', label: 'EN' }]} onChange={setLanguage} label={t('languageLabel')}/><button className="theme-button" type="button" onClick={toggleTheme} aria-label={t(`theme.${theme}`)}>{theme === 'light' ? '☀' : '☾'}</button></div></div></header>;
}
