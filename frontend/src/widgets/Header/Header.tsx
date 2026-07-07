import { useTranslation } from 'react-i18next';
import { usePreferencesStore } from '@/features/preferences/store';
import { Logo } from '@/shared/ui/Logo';
import { NavBar, type NavItem } from '@/shared/ui/NavBar';
import { SegmentedControl } from '@/shared/ui/Controls';
import { UserPill } from '@/shared/ui/Avatar';
import { useAuthStore } from '@/features/auth/store';
import type { AuthModalMode } from '@/features/auth/AuthModal';
import { LogoutIcon } from '@/shared/icons';

function getInitials(username: string) {
  return username.replace(/[^a-zA-Z0-9а-яА-ЯёЁ]/g, '').slice(0, 2).toUpperCase() || '✦';
}

export function Header({ onAuthRequest }: { onAuthRequest: (mode: AuthModalMode) => void }) {
  const { t } = useTranslation();
  const theme = usePreferencesStore((state) => state.theme);
  const language = usePreferencesStore((state) => state.language);
  const setLanguage = usePreferencesStore((state) => state.setLanguage);
  const toggleTheme = usePreferencesStore((state) => state.toggleTheme);
  const user = useAuthStore((state) => state.user);
  const token = useAuthStore((state) => state.token);
  const status = useAuthStore((state) => state.status);
  const logout = useAuthStore((state) => state.logout);
  const isAuthenticated = Boolean(user && token);
  const items: NavItem[] = [
    { id: 'daily', label: t('nav.daily') }, { id: 'spreads', label: t('nav.spreads') },
    { id: 'horoscope', label: t('nav.horoscope') }, { id: 'learn', label: t('nav.learn') },
  ];
  const displayName = status === 'loading' && token && !user ? t('user.loading') : user?.username || t('user.guest');
  const secondary = user ? t(`roles.${user.role}`) : t('user.signIn');
  const initials = user ? getInitials(user.username) : language === 'ru' ? 'Г' : 'G';

  return <header className="site-header"><div className="header-inner"><Logo tagline={t('tagline')}/><NavBar items={items} active="daily" label={t('navigationLabel')} onSelect={(item) => { if (item.id === 'daily' || isAuthenticated) return; onAuthRequest('register'); }}/><div className="header-actions"><UserPill initials={initials} name={displayName} sign={secondary} onClick={isAuthenticated || (status === 'loading' && Boolean(token)) ? undefined : () => onAuthRequest('login')}/>{isAuthenticated && <button className="logout-button" type="button" onClick={logout} aria-label={t('auth.logout')}><LogoutIcon width="16" height="16"/></button>}<SegmentedControl value={language} options={[{ value: 'ru', label: 'RU' }, { value: 'en', label: 'EN' }]} onChange={setLanguage} label={t('languageLabel')}/><button className="theme-button" type="button" onClick={toggleTheme} aria-label={t(`theme.${theme}`)}>{theme === 'light' ? '☀' : '☾'}</button></div></div></header>;
}
