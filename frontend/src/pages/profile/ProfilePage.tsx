import { useState, type KeyboardEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { ProfileHistory } from '@/features/profile/ProfileHistory';
import { ProfileSettings } from '@/features/profile/ProfileSettings';
import { useProfile } from '@/features/profile/useProfile';
import { usePreferencesStore } from '@/features/preferences/store';
import { ProfileIcon } from '@/shared/icons';
import { Ornament } from '@/shared/ui/Ornament';
import { Footer } from '@/widgets/Footer/Footer';
import { Header } from '@/widgets/Header/Header';
import { BackgroundDecor } from '@/widgets/StarField/StarField';
import './profile.css';

type ProfileTab = 'settings' | 'history';

function getInitials(username: string) {
  return username.replace(/[^a-zA-Z0-9а-яА-ЯёЁ]/g, '').slice(0, 2).toUpperCase() || '✦';
}

function formatMemberSince(value: string | undefined, language: 'ru' | 'en', fallback: string) {
  if (!value) return fallback;
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return fallback;
  if (language === 'ru') {
    return new Intl.DateTimeFormat('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' })
      .format(date)
      .replace(/^\d+\s+/, '')
      .replace(/\s+г\.$/, '');
  }
  return new Intl.DateTimeFormat('en-US', { month: 'long', year: 'numeric' }).format(date);
}

export default function ProfilePage() {
  const { t } = useTranslation(['common', 'profile']);
  const language = usePreferencesStore((state) => state.language);
  const profile = useProfile();
  const [activeTab, setActiveTab] = useState<ProfileTab>('settings');
  const tabs: ProfileTab[] = ['settings', 'history'];

  const selectAdjacentTab = (event: KeyboardEvent<HTMLButtonElement>, current: ProfileTab) => {
    if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) return;
    event.preventDefault();
    const currentIndex = tabs.indexOf(current);
    const nextIndex = event.key === 'Home'
      ? 0
      : event.key === 'End'
        ? tabs.length - 1
        : (currentIndex + (event.key === 'ArrowRight' ? 1 : -1) + tabs.length) % tabs.length;
    const nextTab = tabs[nextIndex];
    setActiveTab(nextTab);
    document.getElementById(`profile-tab-${nextTab}`)?.focus();
  };

  const content = (() => {
    if (profile.isPending) {
      return <div className="profile-state" role="status">{t('profile:loading')}</div>;
    }
    if (profile.isError || !profile.data) {
      return <div className="profile-state profile-state-error" role="alert">{t('profile:error')}</div>;
    }

    const user = profile.data;
    const memberSince = formatMemberSince(user.created_at, language, t('profile:mock.memberSince'));

    return (
      <>
        <div className="profile-top-grid">
          <aside className="profile-identity-card" aria-label={t('profile:identityLabel')}>
            <div className="profile-card-halo" aria-hidden="true" />
            <div className="profile-avatar" aria-hidden="true">
              <span className="profile-avatar-ring" />
              {getInitials(user.username)}
            </div>
            <div className="profile-name">{user.username}</div>
            <div className="profile-sign">☽ {t('profile:mock.sign')} · {t('profile:memberSince', { date: memberSince })}</div>
            <div className="profile-chips">
              <span>{t('profile:mock.chipMajor')}</span>
              <span>{t('profile:mock.chipLunar')}</span>
            </div>
            <div className="profile-stats">
              <div><strong>{t('profile:mock.readingsValue')}</strong><span>{t('profile:stats.readings')}</span></div>
              <div><strong>{t('profile:mock.streakValue')}</strong><span>{t('profile:stats.streak')}</span></div>
              <div><strong>{t('profile:mock.soulCardValue')}</strong><span>{t('profile:stats.soulCard')}</span></div>
            </div>
            <blockquote className="profile-quote">{t('profile:mock.quote')}</blockquote>
          </aside>

          <section className="profile-panel" aria-labelledby="profile-information-title">
            <div className="profile-panel-heading">
              <ProfileIcon />
              <h2 id="profile-information-title">{t('profile:sectionProfile')}</h2>
            </div>
            <dl className="profile-fields">
              <div><dt>{t('profile:fields.name')}</dt><dd>{user.username}</dd></div>
              <div><dt>{t('profile:fields.birthDate')}</dt><dd>{t('profile:mock.birthDate')}</dd></div>
              <div><dt>{t('profile:fields.birthPlace')}</dt><dd>{t('profile:mock.birthPlace')}</dd></div>
              <div><dt>{t('profile:fields.email')}</dt><dd>{user.email}</dd></div>
            </dl>
            <div className="profile-note">
              <span>{t(`common:roles.${user.role}`)} · {user.timezone}</span>
              <p>{t('profile:profileNote')}</p>
            </div>
          </section>
        </div>

        <div className="profile-tabs" role="tablist" aria-label={t('profile:tabs.label')}>
          {tabs.map((tab) => (
            <button
              id={`profile-tab-${tab}`}
              key={tab}
              type="button"
              role="tab"
              className={activeTab === tab ? 'active' : ''}
              aria-selected={activeTab === tab}
              aria-controls={`profile-panel-${tab}`}
              tabIndex={activeTab === tab ? 0 : -1}
              onClick={() => setActiveTab(tab)}
              onKeyDown={(event) => selectAdjacentTab(event, tab)}
            >
              {t(`profile:tabs.${tab}`)}
            </button>
          ))}
        </div>
        {activeTab === 'settings' ? (
          <div
            id="profile-panel-settings"
            className="profile-settings-tab-panel"
            role="tabpanel"
            aria-labelledby="profile-tab-settings"
            tabIndex={0}
          >
            <ProfileSettings user={user} />
          </div>
        ) : (
          <div
            id="profile-panel-history"
            className="profile-history-tab-panel"
            role="tabpanel"
            aria-labelledby="profile-tab-history"
            tabIndex={0}
          >
            <ProfileHistory />
          </div>
        )}
      </>
    );
  })();

  return (
    <div className="app-shell profile-shell">
      <BackgroundDecor />
      <Header />
      <main className="page-wrap profile-page">
        <header className="profile-page-heading">
          <h1>{t('profile:title')}</h1>
          <p>{t('profile:subtitle')}</p>
        </header>
        <div className="profile-ornament"><Ornament /></div>
        <div className="profile-cabinet">{content}</div>
      </main>
      <Footer />
    </div>
  );
}
