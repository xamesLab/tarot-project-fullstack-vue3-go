import { usePreferencesStore } from '@/features/preferences/store';
import { Logo } from '@/shared/ui/Logo';
import { ControlSections } from './ControlSections';
import { FoundationSections } from './FoundationSections';
import { ShowcaseSections } from './ShowcaseSections';
import './ui-kit.css';

const links = ['palette', 'typography', 'buttons', 'chips', 'fields', 'switchers', 'avatars', 'navigation', 'ornaments', 'tarot', 'aspects', 'panels', 'tooltips', 'modal', 'icons', 'logo'];

export default function UiKitPage() {
  const theme = usePreferencesStore((state) => state.theme); const toggleTheme = usePreferencesStore((state) => state.toggleTheme);
  return <div className="ui-kit"><aside className="kit-rail"><Logo size="sm"/><div className="rail-kicker">UI LIBRARY · V1</div><nav>{links.map((link, index) => <a href={`#${link}`} key={link}><span>{String(index + 1).padStart(2, '0')}</span>{link}</a>)}</nav><button className="kit-theme" onClick={toggleTheme}>{theme === 'dark' ? '☾ Dark Theme' : '☀ Light Theme'}</button><a className="home-link" href="/">← Home</a></aside><main className="kit-content"><header className="kit-hero"><span>NEBULA-TAROT · DESIGN SYSTEM</span><h1>Библиотека<br/><em>интерфейса</em></h1><p>Живая витрина компонентов, токенов и интерактивных состояний.</p></header><FoundationSections/><ControlSections/><ShowcaseSections/></main></div>;
}
