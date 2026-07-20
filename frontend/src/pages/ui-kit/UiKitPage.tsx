import { useEffect, useRef, useState, type KeyboardEvent } from 'react';
import { usePreferencesStore } from '@/features/preferences/store';
import { CloseIcon, LogoSigil } from '@/shared/icons';
import { ControlSections } from './ControlSections';
import { FoundationSections } from './FoundationSections';
import { ShowcaseSections } from './ShowcaseSections';
import './ui-kit.css';

const links = [
  ['foundations', 'Основы'],
  ['typography', 'Типографика'],
  ['buttons', 'Кнопки'],
  ['chips', 'Чипы & Бейджи'],
  ['forms', 'Поля ввода'],
  ['controls', 'Переключатели'],
  ['avatars', 'Аватары'],
  ['nav', 'Навигация'],
  ['ornaments', 'Орнаменты'],
  ['card', 'Карта Таро'],
  ['aspects', 'Аспект-блоки'],
  ['panel', 'Панель'],
  ['tooltips', 'Подсказки'],
  ['modal', 'Модалка'],
  ['icons', 'Иконки & Декор'],
  ['logo', 'Логотип'],
] as const;

export default function UiKitPage() {
  const theme = usePreferencesStore((state) => state.theme);
  const toggleTheme = usePreferencesStore((state) => state.toggleTheme);
  const [menuOpen, setMenuOpen] = useState(false);
  const menuButtonRef = useRef<HTMLButtonElement>(null);
  const drawerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!menuOpen) return undefined;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    const focusTimer = window.setTimeout(() => drawerRef.current?.querySelector<HTMLAnchorElement>('a')?.focus(), 20);

    return () => {
      window.clearTimeout(focusTimer);
      document.body.style.overflow = previousOverflow;
    };
  }, [menuOpen]);

  const closeMenu = () => {
    setMenuOpen(false);
    window.setTimeout(() => menuButtonRef.current?.focus(), 0);
  };

  const handleDrawerKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
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
    <div className="ui-kit">
      <div className="kit-orb kit-orb-one" />
      <div className="kit-orb kit-orb-two" />
      <div className="kit-shell">
        <header className="kit-header">
          <button
            ref={menuButtonRef}
            type="button"
            className="kit-menu-button"
            aria-label="Открыть меню компонентов"
            aria-expanded={menuOpen}
            aria-controls="kit-mobile-drawer"
            onClick={() => setMenuOpen(true)}
          >
            <span /><span /><span />
          </button>
          <div className="kit-brand">
            <LogoSigil width="36" height="36" />
            <div>
              <div className="kit-brand-name">ARCANA</div>
              <div className="kit-brand-sub">UI Library</div>
            </div>
          </div>
          <div className="kit-header-actions">
            <button
              type="button"
              className="kit-theme"
              aria-label={theme === 'dark' ? 'Включить светлую тему' : 'Включить тёмную тему'}
              onClick={toggleTheme}
            >
              <span className="kit-theme-indicator" />
              <span className="kit-theme-label">{theme === 'dark' ? 'Тёмная тема' : 'Светлая тема'}</span>
            </button>
          </div>
        </header>
        <aside className="kit-rail" aria-label="Разделы библиотеки компонентов">
          <div className="rail-kicker">Компоненты</div>
          <nav aria-label="Навигация по UI Kit">
            {links.map(([id, label]) => (
              <a href={`#${id}`} key={id}>
                {label}
              </a>
            ))}
          </nav>
        </aside>
        <main className="kit-content">
          <header className="kit-lede">
            <span>Дизайн-система · {theme === 'dark' ? 'Тёмная тема' : 'Светлая тема'}</span>
            <h1>Библиотека компонентов ARCANA</h1>
            <p>
              Переиспользуемые элементы интерфейса, извлечённые из карты дня: золото на чернильном
              фоне, засечные акценты Cinzel и Cormorant, тонкие астральные орнаменты.
            </p>
          </header>
          <FoundationSections />
          <ControlSections />
          <ShowcaseSections />
        </main>
      </div>
      {menuOpen && (
        <div
          className="kit-drawer-overlay"
          onMouseDown={(event) => {
            if (event.target === event.currentTarget) closeMenu();
          }}
        >
          <div
            ref={drawerRef}
            id="kit-mobile-drawer"
            className="kit-drawer"
            role="dialog"
            aria-modal="true"
            aria-labelledby="kit-drawer-title"
            onKeyDown={handleDrawerKeyDown}
          >
            <header className="kit-drawer-head">
              <div>
                <span>UI Library</span>
                <h2 id="kit-drawer-title">Компоненты</h2>
              </div>
              <button type="button" className="kit-drawer-close" aria-label="Закрыть меню компонентов" onClick={closeMenu}>
                <CloseIcon />
              </button>
            </header>
            <nav className="kit-drawer-nav" aria-label="Мобильная навигация по UI Kit">
              {links.map(([id, label], index) => (
                <a href={`#${id}`} key={id} onClick={() => setMenuOpen(false)}>
                  <span>{String(index + 1).padStart(2, '0')}</span>
                  {label}
                </a>
              ))}
            </nav>
          </div>
        </div>
      )}
    </div>
  );
}
