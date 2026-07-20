import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { api } from '@/api';
import type { UserRole } from '@/api/types';
import { Button } from '@/shared/ui/Button';
import { adminKeys, getAdminError } from './shared';

export function AdminDashboard({ onOpenTab }: { onOpenTab: (tab: 'users' | 'assets') => void }) {
  const { t } = useTranslation(['admin', 'common']);
  const query = useQuery({ queryKey: adminKeys.dashboard, queryFn: api.admin.getDashboard });

  if (query.isPending) return <AdminState>{t('admin:states.loading')}</AdminState>;
  if (query.isError) return <AdminState error>{getAdminError(query.error, t('admin:states.error'))}</AdminState>;

  const metrics = query.data.metrics;
  const metricCards = [
    { key: 'users', label: t('admin:dashboard.totalUsers'), value: metrics.users.total, meta: t('admin:dashboard.usersMeta', metrics.users.by_status) },
    { key: 'today', label: t('admin:dashboard.readingsToday'), value: metrics.readings.today.total, meta: t('admin:dashboard.readingsMeta', metrics.readings.today) },
    { key: 'week', label: t('admin:dashboard.readingsWeek'), value: metrics.readings.this_week.total, meta: t('admin:dashboard.readingsMeta', metrics.readings.this_week) },
    { key: 'decks', label: t('admin:dashboard.totalDecks'), value: metrics.decks.total, meta: t('admin:dashboard.decksMeta', metrics.decks.by_status) },
  ];
  const roles: UserRole[] = ['seeker', 'adept', 'master', 'keeper'];

  return (
    <div className="admin-stack">
      <section className="admin-metrics" aria-label={t('admin:dashboard.metricsLabel')}>
        {metricCards.map((card) => <article className="admin-metric-card" key={card.key}><span>{card.label}</span><strong>{card.value}</strong><p>{card.meta}</p></article>)}
      </section>
      <section className="admin-surface admin-dashboard-section">
        <div className="admin-section-head"><div><span>{t('admin:dashboard.usersEyebrow')}</span><h2>{t('admin:dashboard.usersTitle')}</h2><p>{t('admin:dashboard.usersCopy')}</p></div><Button variant="outline" size="sm" onClick={() => onOpenTab('users')}>{t('admin:dashboard.openUsers')}</Button></div>
        <div className="admin-summary-grid">
          {roles.map((role) => <article key={role}><strong>{metrics.users.by_role[role]}</strong><span>{t(`common:roles.${role}`)}</span></article>)}
        </div>
      </section>
      <div className="admin-dashboard-grid">
        <section className="admin-surface admin-dashboard-section">
          <div className="admin-section-head"><div><span>{t('admin:dashboard.readingsEyebrow')}</span><h2>{t('admin:dashboard.readingsTitle')}</h2><p>{t('admin:dashboard.readingsCopy')}</p></div></div>
          <div className="admin-summary-grid admin-summary-grid-2">
            <article><strong>{metrics.readings.today.card_of_day}</strong><span>{t('admin:dashboard.dailyToday')}</span></article>
            <article><strong>{metrics.readings.today.three_card}</strong><span>{t('admin:dashboard.spreadsToday')}</span></article>
            <article><strong>{metrics.readings.this_week.card_of_day}</strong><span>{t('admin:dashboard.dailyWeek')}</span></article>
            <article><strong>{metrics.readings.this_week.three_card}</strong><span>{t('admin:dashboard.spreadsWeek')}</span></article>
          </div>
        </section>
        <section className="admin-surface admin-dashboard-section">
          <div className="admin-section-head"><div><span>{t('admin:dashboard.decksEyebrow')}</span><h2>{t('admin:dashboard.decksTitle')}</h2><p>{t('admin:dashboard.decksCopy')}</p></div><Button variant="outline" size="sm" onClick={() => onOpenTab('assets')}>{t('admin:dashboard.openAssets')}</Button></div>
          <div className="admin-summary-grid admin-summary-grid-2">
            {(['done', 'processing', 'disabled', 'archived'] as const).map((status) => <article key={status}><strong>{metrics.decks.by_status[status]}</strong><span>{t(`admin:statuses.${status}`)}</span></article>)}
          </div>
        </section>
      </div>
    </div>
  );
}

export function AdminState({ children, error = false }: { children: React.ReactNode; error?: boolean }) {
  return <div className={`admin-state ${error ? 'admin-state-error' : ''}`} role={error ? 'alert' : 'status'}>{children}</div>;
}
