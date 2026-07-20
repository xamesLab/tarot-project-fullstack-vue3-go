import { useEffect, useRef, useState, type KeyboardEvent, type MouseEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { useLocation, useNavigate } from 'react-router-dom';
import type { AuthModalMode } from '@/features/auth/AuthModal';
import { useAuthStore } from '@/features/auth/store';
import { usePreferencesStore } from '@/features/preferences/store';
import { CloseIcon, LogoutIcon } from '@/shared/icons';
import { UserPill } from '@/shared/ui/Avatar';
import { SegmentedControl } from '@/shared/ui/Controls';
import { Logo } from '@/shared/ui/Logo';
import { NavBar, type NavItem } from '@/shared/ui/NavBar';

function getInitials(username: string) {
  return username.replace(/[^a-zA-Z0-9а-яА-ЯёЁ]/g, '').slice(0, 2).toUpperCase() || '✦';
}

export function Header({ onAuthRequest }: { onAuthRequest?: (mode: AuthModalMode) => void }) {
  const { t } = useTranslation(['common', 'admin']);
  const location = useLocation();
  const navigate = useNavigate();
  const theme = usePreferencesStore((state) => state.theme);
  const language = usePreferencesStore((state) => state.language);
  const setTheme = usePreferencesStore((state) => state.setTheme);
  const setLanguage = usePreferencesStore((state) => state.setLanguage);
  const toggleTheme = usePreferencesStore((state) => state.toggleTheme);
  const user = useAuthStore((state) => state.user);
  const token = useAuthStore((state) => state.token);
  const status = useAuthStore((state) => state.status);
  const logout = useAuthStore((state) => state.logout);
  const [menuOpen, setMenuOpen] = useState(false);
  const menuButtonRef = useRef<HTMLButtonElement>(null);
  const drawerRef = useRef<HTMLElement>(null);
  const isAuthenticated = Boolean(user && token);

  const items: NavItem[] = [
    { id: 'daily', label: t('nav.daily') },
    { id: 'spreads', label: t('nav.spreads'), href: isAuthenticated ? '/spreads' : undefined },
    ...(isAuthenticated ? [{ id: 'profile', label: t('nav.profile'), href: '/profile' }] : []),
    ...(user?.role === 'keeper' ? [{ id: 'admin', label: t('admin:nav'), href: '/admin' }] : []),
  ];
  const displayName = status === 'loading' && token && !user ? t('user.loading') : user?.username || t('user.guest');
  const secondary = user ? t(`roles.${user.role}`) : t('user.signIn');
  const initials = user ? getInitials(user.username) : language === 'ru' ? 'Г' : 'G';
  const activeItem = location.pathname.startsWith('/admin')
    ? 'admin'
    : location.pathname === '/profile'
      ? 'profile'
      : location.pathname === '/spreads'
        ? 'spreads'
        : 'daily';

  useEffect(() => {
    if (!menuOpen) return undefined;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    drawerRef.current?.focus();
    return () => { document.body.style.overflow = previousOverflow; };
  }, [menuOpen]);

  const closeMenu = (restoreFocus = true) => {
    setMenuOpen(false);
    if (restoreFocus) window.setTimeout(() => menuButtonRef.current?.focus(), 0);
  };

  const selectItem = (item: NavItem) => {
    if (item.id === 'daily') {
      if (location.pathname !== '/') navigate('/');
      return;
    }
    if (item.id === 'spreads') {
      if (isAuthenticated) {
        if (location.pathname !== '/spreads') navigate('/spreads');
      } else {
        onAuthRequest?.('register');
      }
      return;
    }
    if (item.id === 'profile') {
      if (location.pathname !== '/profile') navigate('/profile');
      return;
    }
    if (item.id === 'admin' && !location.pathname.startsWith('/admin')) navigate('/admin');
  };

  const selectDrawerItem = (item: NavItem, event?: MouseEvent<HTMLAnchorElement>) => {
    if (event && (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey)) return;
    event?.preventDefault();
    closeMenu(false);
    selectItem(item);
  };

  const openAccount = () => {
    closeMenu(false);
    if (isAuthenticated) {
      if (location.pathname !== '/profile') navigate('/profile');
    } else if (!(status === 'loading' && token)) {
      onAuthRequest?.('login');
    }
  };

  const handleDrawerKeyDown = (event: KeyboardEvent<HTMLElement>) => {
    if (event.key === 'Escape') {
      closeMenu();
      return;
    }
    if (event.key !== 'Tab' || !drawerRef.current) return;
    const focusable = Array.from(drawerRef.current.querySelectorAll<HTMLElement>('a, button:not([disabled])'));
    const first = focusable[0];
    const last = focusable[focusable.length - 1];
    if (event.shiftKey && document.activeElement === first) {
      event.preventDefault();
      last?.focus();
    } else if (!event.shiftKey && document.activeElement === last) {
      event.preventDefault();
      first?.focus();
    }
  };

  return (
    <header className="site-header">
      <div className="header-inner">
        <button
          ref={menuButtonRef}
          type="button"
          className="site-menu-button"
          aria-label={t('menu.open')}
          aria-expanded={menuOpen}
          aria-controls="site-mobile-drawer"
          onClick={() => setMenuOpen(true)}
        >
          <span /><span /><span />
        </button>
        <Logo tagline={t('tagline')} />
        <NavBar items={items} active={activeItem} label={t('navigationLabel')} onSelect={selectItem} />
        <div className="header-actions">
          <UserPill initials={initials} name={displayName} sign={secondary} onClick={openAccount} />
          {isAuthenticated && (
            <button className="logout-button" type="button" onClick={logout} aria-label={t('auth.logout')}>
              <LogoutIcon width="16" height="16" />
            </button>
          )}
          <SegmentedControl
            value={language}
            options={[{ value: 'ru', label: 'RU' }, { value: 'en', label: 'EN' }]}
            onChange={setLanguage}
            label={t('languageLabel')}
          />
          <button className="theme-button" type="button" onClick={toggleTheme} aria-label={t(`theme.${theme}`)}>
            {theme === 'light' ? '☀' : '☾'}
          </button>
        </div>
      </div>

      {menuOpen && (
        <div
          className="site-drawer-overlay"
          onMouseDown={(event) => { if (event.target === event.currentTarget) closeMenu(); }}
        >
          <aside
            ref={drawerRef}
            id="site-mobile-drawer"
            className="site-mobile-drawer"
            role="dialog"
            aria-modal="true"
            aria-labelledby="site-drawer-title"
            tabIndex={-1}
            onKeyDown={handleDrawerKeyDown}
          >
            <div className="site-drawer-head">
              <Logo size="sm" tagline={t('tagline')} />
              <button type="button" className="site-drawer-close" aria-label={t('menu.close')} onClick={() => closeMenu()}>
                <CloseIcon />
              </button>
            </div>
            <div className="site-drawer-body">
              <section className="site-drawer-section">
                <h2 id="site-drawer-title">{t('menu.title')}</h2>
                <nav className="site-drawer-nav" aria-label={t('navigationLabel')}>
                  {items.map((item, index) => {
                    const content = <><span aria-hidden="true">{String(index + 1).padStart(2, '0')}</span>{item.label}</>;
                    const className = item.id === activeItem ? 'active' : '';
                    return item.href ? (
                      <a key={item.id} className={className} href={item.href} aria-current={item.id === activeItem ? 'page' : undefined} onClick={(event) => selectDrawerItem(item, event)}>{content}</a>
                    ) : (
                      <button key={item.id} type="button" className={className} aria-current={item.id === activeItem ? 'page' : undefined} onClick={() => selectDrawerItem(item)}>{content}</button>
                    );
                  })}
                </nav>
              </section>
              <section className="site-drawer-section site-drawer-account">
                <h2>{t('menu.account')}</h2>
                <UserPill initials={initials} name={displayName} sign={secondary} onClick={openAccount} />
                <div className="site-drawer-preferences">
                  <span>{t('languageLabel')}</span>
                  <SegmentedControl
                    value={language}
                    options={[{ value: 'ru', label: 'RU' }, { value: 'en', label: 'EN' }]}
                    onChange={setLanguage}
                    label={t('languageLabel')}
                  />
                </div>
                <div className="site-drawer-preferences">
                  <span>{t('themeLabel')}</span>
                  <SegmentedControl
                    value={theme}
                    options={[
                      { value: 'dark', label: '☾', ariaLabel: t('themeOption.dark') },
                      { value: 'light', label: '☀', ariaLabel: t('themeOption.light') },
                    ]}
                    onChange={setTheme}
                    label={t('themeLabel')}
                  />
                </div>
                {isAuthenticated && (
                  <button type="button" className="site-drawer-logout" onClick={() => { closeMenu(false); logout(); }}>
                    <LogoutIcon width="15" height="15" />
                    {t('auth.logout')}
                  </button>
                )}
              </section>
            </div>
          </aside>
        </div>
      )}
    </header>
  );
}
