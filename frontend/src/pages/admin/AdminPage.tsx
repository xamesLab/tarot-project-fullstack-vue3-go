import { useMemo, type KeyboardEvent } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { AdminAssets } from '@/features/admin/AdminAssets';
import { AdminDashboard } from '@/features/admin/AdminDashboard';
import { AdminPolicies } from '@/features/admin/AdminPolicies';
import { AdminSettings } from '@/features/admin/AdminSettings';
import { AdminUsers } from '@/features/admin/AdminUsers';
import { Ornament } from '@/shared/ui/Ornament';
import { Footer } from '@/widgets/Footer/Footer';
import { Header } from '@/widgets/Header/Header';
import { BackgroundDecor } from '@/widgets/StarField/StarField';
import './admin.css';

export type AdminTab = 'overview' | 'users' | 'policies' | 'settings' | 'assets';
const tabIds: AdminTab[] = ['overview', 'users', 'policies', 'settings', 'assets'];

export default function AdminPage() {
  const { t } = useTranslation('admin');
  const [searchParams, setSearchParams] = useSearchParams();
  const requestedTab = searchParams.get('tab');
  const activeTab: AdminTab = tabIds.includes(requestedTab as AdminTab) ? requestedTab as AdminTab : 'overview';
  const tabs = useMemo(() => tabIds.map((id) => ({ id, label: t(`tabs.${id}`) })), [t]);
  const setTab = (tab: AdminTab) => setSearchParams(tab === 'overview' ? {} : { tab });
  const selectAdjacentTab = (event: KeyboardEvent<HTMLButtonElement>, current: AdminTab) => {
    if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) return;
    event.preventDefault();
    const index = tabIds.indexOf(current);
    const nextIndex = event.key === 'Home' ? 0 : event.key === 'End' ? tabIds.length - 1 : (index + (event.key === 'ArrowRight' ? 1 : -1) + tabIds.length) % tabIds.length;
    const next = tabIds[nextIndex];
    setTab(next);
    document.getElementById(`admin-tab-${next}`)?.focus();
  };
  const content = activeTab === 'overview' ? <AdminDashboard onOpenTab={setTab} /> : activeTab === 'users' ? <AdminUsers /> : activeTab === 'policies' ? <AdminPolicies /> : activeTab === 'settings' ? <AdminSettings /> : <AdminAssets />;

  return (
    <div className="app-shell admin-shell">
      <BackgroundDecor />
      <Header />
      <main className="page-wrap admin-page">
        <header className="admin-page-heading"><span>{t('eyebrow')}</span><h1>{t('title')}</h1><p>{t('subtitle')}</p></header>
        <div className="admin-ornament"><Ornament /></div>
        <nav className="admin-tabs" role="tablist" aria-label={t('tabs.label')}>
          {tabs.map((tab) => <button id={`admin-tab-${tab.id}`} key={tab.id} type="button" role="tab" aria-selected={activeTab === tab.id} aria-controls={`admin-panel-${tab.id}`} tabIndex={activeTab === tab.id ? 0 : -1} className={activeTab === tab.id ? 'active' : ''} onClick={() => setTab(tab.id)} onKeyDown={(event) => selectAdjacentTab(event, tab.id)}>{tab.label}</button>)}
        </nav>
        <section id={`admin-panel-${activeTab}`} className="admin-tab-panel" role="tabpanel" aria-labelledby={`admin-tab-${activeTab}`} tabIndex={0}>
          <header className="admin-tab-heading"><span>{t(`sections.${activeTab}.eyebrow`)}</span><h2>{t(`sections.${activeTab}.title`)}</h2><p>{t(`sections.${activeTab}.copy`)}</p></header>
          {content}
        </section>
      </main>
      <Footer />
    </div>
  );
}
