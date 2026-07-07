import { usePreferencesStore } from '@/features/preferences/store';
import { LogoSigil } from '@/shared/icons';
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
  return (
    <div className="ui-kit">
      <div className="kit-orb kit-orb-one" />
      <div className="kit-orb kit-orb-two" />
      <div className="kit-shell">
        <header className="kit-header">
          <div className="kit-brand">
            <LogoSigil width="36" height="36" />
            <div>
              <div className="kit-brand-name">ARCANA</div>
              <div className="kit-brand-sub">UI Library</div>
            </div>
          </div>
          <button type="button" className="kit-theme" onClick={toggleTheme}>
            <span />
            {theme === 'dark' ? 'Dark Theme' : 'Light Theme'}
          </button>
        </header>
        <aside className="kit-rail">
          <div className="rail-kicker">Компоненты</div>
          <nav>
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
    </div>
  );
}
